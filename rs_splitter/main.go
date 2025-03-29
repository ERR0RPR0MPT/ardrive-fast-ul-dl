package rs_splitter

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/klauspost/reedsolomon"
	"io"
	"math"
	"mime"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	defaultChunkSize        = 97 * 1024
	defaultMaxFolderNum     = 74
	DefaultFileInfoJsonName = "fileInfo.json"
)

type FileInfoManifest struct {
	DataShards   int    `json:"data_shards"`
	ParityShards int    `json:"parity_shards"`
	ChunkSize    int    `json:"chunk_size"`
	IndexLength  int    `json:"index_length"`
	OriginalSize int64  `json:"original_size"`
	Filename     string `json:"filename"`
	MIMETypes    string `json:"mime_types"`
}

func CheckShardNumber(path string) (bool, int) {
	base := filepath.Base(path)
	if !strings.HasSuffix(base, ".bin") {
		return false, 0
	}
	numStr := strings.TrimSuffix(base, ".bin")
	num, err := strconv.Atoi(numStr)
	if err != nil {
		return false, 0
	}
	return true, num
}

func GetShardNumber(path string) int {
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

// getMimeType 根据文件名返回 MIME 类型
func getMimeType(fileName string) string {
	// 获取文件扩展名
	ext := filepath.Ext(fileName)
	// 根据扩展名查找 MIME 类型
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		return "application/octet-stream" // 默认类型
	}
	return mimeType
}

func GetDefaultShardsNum(filePath string) (int, int, int, error) {
	// 获取文件信息
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return 0, 0, 0, err
	}

	fileSize := fileInfo.Size() // 获取文件大小，单位为字节

	// 根据要求计算 a
	a := int(math.Ceil(float64(fileSize) / defaultMaxFolderNum / defaultChunkSize))

	// 计算 b
	b := int(math.Ceil(float64(a) * 0.25))

	return a, b, defaultChunkSize, nil
}

// CalculateShards 计算并处理 allShardsNum
func CalculateShards(originalSize int64, chunkSize, dataShards, parityShards int) int {
	if originalSize <= int64(chunkSize) {
		return 3
	}
	allShardsNum := float64(originalSize) / float64(dataShards) * float64(dataShards+parityShards) / float64(chunkSize)
	// 计算总的 shards 数量
	totalShards := dataShards + parityShards

	// 如果 allShardsNum 不是 totalShards 的倍数，向上进位
	if int(allShardsNum)%totalShards != 0 {
		allShardsNum = math.Ceil(allShardsNum/float64(totalShards)) * float64(totalShards)
	}

	return int(allShardsNum) + 1
}

// RSSplitterEncodeDir encodes all files in the input directory
func RSSplitterEncodeDir(inputDir, outputDir string) error {
	return filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			// Get the relative path for output
			relPath, _ := filepath.Rel(inputDir, path)
			outputFile := filepath.Join(outputDir, relPath)

			// Create the output directory if it doesn't exist
			if err := os.MkdirAll(filepath.Dir(outputFile), 0755); err != nil {
				return fmt.Errorf("failed to create directory: %v", err)
			}

			dataShards, parityShards, chunkSize, err := GetDefaultShardsNum(path)
			if err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}

			// Encode the file
			fmt.Println("dataShards:", dataShards)
			fmt.Println("parityShards:", parityShards)
			fmt.Println("chunkSize:", chunkSize)
			outputFile = filepath.Join(filepath.Dir(outputFile), strings.ReplaceAll(filepath.Base(outputFile), ".", "-"))
			fmt.Println("Encode:", path, "->", outputFile)
			if err := RSSplitterEncode(path, outputFile, dataShards, parityShards, chunkSize); err != nil {
				return fmt.Errorf("error encoding file %s: %v", path, err)
			}
		}
		return nil
	})
}

func RSSplitterDecodeDir(inputDir, outputDir string) error {
	return filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			manifestPath := filepath.Join(path, DefaultFileInfoJsonName)
			if _, err := os.Stat(manifestPath); err == nil {
				// 读取fileInfo.json
				manifestData, err := os.ReadFile(manifestPath)
				if err != nil {
					return fmt.Errorf("error reading manifest in %s: %v", path, err)
				}
				var manifest FileInfoManifest
				if err := json.Unmarshal(manifestData, &manifest); err != nil {
					return fmt.Errorf("error parsing manifest in %s: %v", path, err)
				}

				// 构造输出文件路径
				relPath, err := filepath.Rel(inputDir, path)
				if err != nil {
					return fmt.Errorf("error getting relative path for %s: %v", path, err)
				}
				parentRelPath := filepath.Dir(relPath)
				outputParentDir := filepath.Join(outputDir, parentRelPath)
				outputFilePath := filepath.Join(outputParentDir, manifest.Filename)

				// 创建输出目录
				if err := os.MkdirAll(outputParentDir, 0755); err != nil {
					return fmt.Errorf("error creating output directory %s: %v", outputParentDir, err)
				}

				// 解码该目录
				fmt.Printf("Decoding: %s -> %s\n", path, outputFilePath)
				if err := RSSplitterDecode(path, outputFilePath); err != nil {
					return fmt.Errorf("error decoding %s: %v", path, err)
				}
			}
		}
		return nil
	})
}

func RSSplitterEncode(inputFile, outputDir string, dataShards, parityShards, chunkSize int) error {
	totalShards := dataShards + parityShards
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

	indexLength := 0
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
		indexLength = endNum
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

	manifest := FileInfoManifest{
		DataShards:   dataShards,
		ParityShards: parityShards,
		ChunkSize:    chunkSize,
		IndexLength:  indexLength,
		OriginalSize: originalSize,
		Filename:     filepath.Base(inputFile),
		MIMETypes:    getMimeType(inputFile),
	}

	manifestData, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(outputDir, DefaultFileInfoJsonName), manifestData, 0644)
}

func RSSplitterDecode(inputDir, outputFile string) error {
	manifestData, err := os.ReadFile(filepath.Join(inputDir, DefaultFileInfoJsonName))
	if err != nil {
		return err
	}

	var manifest FileInfoManifest
	if err := json.Unmarshal(manifestData, &manifest); err != nil {
		return err
	}
	dataShards := manifest.DataShards
	parityShards := manifest.ParityShards
	totalShards := dataShards + parityShards
	allShardsNum := CalculateShards(manifest.OriginalSize, manifest.ChunkSize, dataShards, parityShards)
	filename := manifest.Filename
	MIMETypes := manifest.MIMETypes

	fmt.Println("dataShards:", dataShards)
	fmt.Println("parityShards:", parityShards)
	fmt.Println("totalShards:", totalShards)
	fmt.Println("allShardsNum:", allShardsNum)
	fmt.Println("filename:", filename)
	fmt.Println("MIMETypes:", MIMETypes)

	encoder, err := reedsolomon.New(manifest.DataShards, manifest.ParityShards)
	if err != nil {
		return err
	}

	// Collect all shard files
	shardFiles := make([]string, allShardsNum)
	filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".bin") {
			shardFiles[GetShardNumber(path)] = path
		}
		return nil
	})

	//// Sort shards by number
	//sort.Slice(shardFiles, func(i, j int) bool {
	//	return GetShardNumber(shardFiles[i]) < GetShardNumber(shardFiles[j])
	//})

	outFile, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Process groups
	for i := 1; i < len(shardFiles); {
		//fmt.Println("Check:", i, "-", getMin(i+totalShards, len(shardFiles)))
		groupShards := shardFiles[i:getMin(i+totalShards, len(shardFiles))]
		i += len(groupShards)

		// Initialize shard buffers
		shards := make([][]byte, totalShards)
		present := make([]bool, totalShards)

		// Load shards
		for _, path := range groupShards {
			if path == "" {
				continue
			}
			num := GetShardNumber(path)
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
		for j := 0; j < dataShards; j++ {
			if !present[j] {
				fmt.Println("Reconstruct missing data shards index:", i+j-len(groupShards))
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
			//if available < dataShards {
			//	return fmt.Errorf("not enough shards for reconstruction")
			//}

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
