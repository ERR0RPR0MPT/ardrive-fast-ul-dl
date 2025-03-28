package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

const fileNum = 100
const chunkSize = 1024 * 97

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: ./splitter encode|decode <input>")
		os.Exit(1)
	}

	cmd := os.Args[1]
	input := os.Args[2]

	switch cmd {
	case "encode":
		outputDir := strings.ReplaceAll(filepath.Base(input), ".", "-")
		if err := encode(input, outputDir); err != nil {
			fmt.Printf("Encode error: %v\n", err)
			os.Exit(1)
		}
	case "decode":
		outputFile := strings.ReplaceAll(filepath.Base(input), "-", ".")
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

	buf := make([]byte, chunkSize)
	index := 1

	for {
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		start := ((index-1)/fileNum)*fileNum + 1
		end := start + fileNum - 1
		subDir := fmt.Sprintf("%d-%d", start, end)
		dirPath := filepath.Join(outputDir, subDir)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return err
		}

		chunkFile := filepath.Join(dirPath, fmt.Sprintf("%d.bin", index))
		if err := os.WriteFile(chunkFile, buf[:n], 0644); err != nil {
			return err
		}

		index++
	}
	return nil
}

func decode(inputDir, outputFile string) error {
	var files []string

	err := filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".bin" {
			filename := filepath.Base(path)
			nameWithoutExt := strings.TrimSuffix(filename, ".bin")
			if _, err := strconv.Atoi(nameWithoutExt); err == nil {
				files = append(files, path)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	sort.Slice(files, func(i, j int) bool {
		return getIndex(files[i]) < getIndex(files[j])
	})

	out, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer out.Close()

	for _, fpath := range files {
		f, err := os.Open(fpath)
		if err != nil {
			return err
		}

		if _, err := io.Copy(out, f); err != nil {
			f.Close()
			return err
		}
		f.Close()
	}
	return nil
}

func getIndex(path string) int {
	filename := filepath.Base(path)
	name := strings.TrimSuffix(filename, ".bin")
	index, _ := strconv.Atoi(name)
	return index
}
