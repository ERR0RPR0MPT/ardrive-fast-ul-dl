package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/klauspost/reedsolomon"
)

const (
	dataShards   = 200
	parityShards = 50
	totalShards  = dataShards + parityShards
	chunkSize    = 97 * 1024 // 97KB
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: rs-splitter encode|decode <input>")
		os.Exit(1)
	}

	cmd := os.Args[1]
	input := os.Args[2]

	switch cmd {
	case "encode":
		outputDir := strings.ReplaceAll(filepath.Base(input), ".", "-") + "-split"
		if err := encode(input, outputDir); err != nil {
			fmt.Printf("Encode error: %v\n", err)
			os.Exit(1)
		}
	case "decode":
		outputFile := strings.ReplaceAll(filepath.Base(input), "-split", "")
		if err := decode(input, outputFile); err != nil {
			fmt.Printf("Decode error: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Println("Unknown command. Use encode or decode.")
		os.Exit(1)
	}
}

func encode(inputFile, outputDir string) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}

	file, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}
	originalSize := fileInfo.Size()

	encoder, err := reedsolomon.New(dataShards, parityShards)
	if err != nil {
		return err
	}

	group := 0
	buffer := make([]byte, dataShards*chunkSize)

	for {
		n, err := io.ReadFull(file, buffer)
		if err != nil && err != io.EOF && !errors.Is(err, io.ErrUnexpectedEOF) {
			return err
		}

		if n == 0 {
			break
		}

		// Pad last block with zeros if necessary
		if n < len(buffer) {
			clear(buffer[n:])
		}

		// Split into data shards
		shards := make([][]byte, totalShards)
		for i := 0; i < dataShards; i++ {
			shards[i] = buffer[i*chunkSize : (i+1)*chunkSize]
		}

		// Initialize parity shards
		for i := dataShards; i < totalShards; i++ {
			shards[i] = make([]byte, chunkSize)
		}

		// Generate parity shards
		if err := encoder.Encode(shards); err != nil {
			return err
		}

		// Calculate directory name
		startNum := group*totalShards + 1
		endNum := (group + 1) * totalShards
		subDir := fmt.Sprintf("%d-%d", startNum, endNum)
		dirPath := filepath.Join(outputDir, subDir)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return err
		}

		// Write all shards
		for i := 0; i < totalShards; i++ {
			shardNum := startNum + i
			shardPath := filepath.Join(dirPath, fmt.Sprintf("%d.bin", shardNum))
			if err := os.WriteFile(shardPath, shards[i], 0644); err != nil {
				return err
			}
		}

		group++

		if err == io.EOF {
			break
		}
	}

	// Create manifest
	manifest := map[string]interface{}{
		"data_shards":   dataShards,
		"parity_shards": parityShards,
		"chunk_size":    chunkSize,
		"original_size": originalSize,
	}

	manifestData, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(outputDir, "manifest.json"), manifestData, 0644)
}

func decode(inputDir, outputFile string) error {
	// Read manifest
	manifestData, err := os.ReadFile(filepath.Join(inputDir, "manifest.json"))
	if err != nil {
		return err
	}

	var manifest struct {
		DataShards   int   `json:"data_shards"`
		ParityShards int   `json:"parity_shards"`
		ChunkSize    int   `json:"chunk_size"`
		OriginalSize int64 `json:"original_size"`
	}
	if err := json.Unmarshal(manifestData, &manifest); err != nil {
		return err
	}

	encoder, err := reedsolomon.New(manifest.DataShards, manifest.ParityShards)
	if err != nil {
		return err
	}

	// Collect all shard files
	var shardFiles []string
	filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".bin") {
			shardFiles = append(shardFiles, path)
		}
		return nil
	})

	// Sort shards by number
	sort.Slice(shardFiles, func(i, j int) bool {
		return getShardNumber(shardFiles[i]) < getShardNumber(shardFiles[j])
	})

	outFile, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Process groups
	for i := 0; i < len(shardFiles); {
		groupShards := shardFiles[i:getMin(i+totalShards, len(shardFiles))]
		i += len(groupShards)

		// Initialize shard buffers
		shards := make([][]byte, totalShards)
		present := make([]bool, totalShards)

		// Load shards
		for _, path := range groupShards {
			num := getShardNumber(path)
			groupIndex := (num - 1) % totalShards
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			shards[groupIndex] = data
			present[groupIndex] = true
		}

		// Reconstruct missing data shards
		needReconstruct := false
		for i := 0; i < dataShards; i++ {
			if !present[i] {
				needReconstruct = true
				break
			}
		}

		if needReconstruct {
			// Verify we have enough shards
			available := 0
			for _, p := range present {
				if p {
					available++
				}
			}
			if available < dataShards {
				return fmt.Errorf("not enough shards for reconstruction")
			}

			if err := encoder.Reconstruct(shards); err != nil {
				return err
			}
		}

		// Write data shards to output
		for i := 0; i < dataShards; i++ {
			if _, err := outFile.Write(shards[i]); err != nil {
				return err
			}
		}
	}

	// Truncate to original size
	if err := outFile.Truncate(manifest.OriginalSize); err != nil {
		return err
	}

	return nil
}

func getShardNumber(path string) int {
	base := filepath.Base(path)
	numStr := strings.TrimSuffix(base, ".bin")
	num, _ := strconv.Atoi(numStr)
	return num
}

func getMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}
