package main

import (
	"fmt"
	"github.com/ERR0RPR0MPT/ardrive-fast-ul-dl/rs_splitter"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	RSSplitterSuffix = "-rs-splitter"
)

func IsDir(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return info.IsDir(), nil
}

func OutputCmdHelp() {
	fmt.Println("Usage: rs-splitter encode <input_file> [dataShards] [parityShards] [chunkSize] | " +
		"encode <input_dir> [output_dir] | " +
		"decode <input_file> | " +
		"decode <input_dir> [output_dir]")
	os.Exit(1)
}

func main() {
	if len(os.Args) < 3 {
		OutputCmdHelp()
	}

	cmd := os.Args[1]
	input := os.Args[2]

	isDir, err := IsDir(input)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	dataShards, parityShards, chunkSize, err := rs_splitter.GetDefaultShardsNum(input)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	switch cmd {
	case "encode":
		if !isDir {
			if len(os.Args) >= 4 {
				if ds, err := strconv.Atoi(os.Args[3]); err == nil {
					dataShards = ds
				}
			}
			if len(os.Args) >= 5 {
				if ps, err := strconv.Atoi(os.Args[4]); err == nil {
					parityShards = ps
				}
			}
			if len(os.Args) >= 6 {
				if cs, err := strconv.Atoi(os.Args[5]); err == nil {
					chunkSize = cs
				}
			}
			fmt.Println("dataShards:", dataShards)
			fmt.Println("parityShards:", parityShards)
			fmt.Println("chunkSize:", chunkSize)
			outputDir := strings.ReplaceAll(filepath.Base(input), ".", "-")
			if err := rs_splitter.RSSplitterEncode(input, outputDir, dataShards, parityShards, chunkSize); err != nil {
				fmt.Printf("Encode error: %v\n", err)
				os.Exit(1)
			}
		} else {
			//if len(os.Args) < 4 {
			//	OutputCmdHelp()
			//}
			inputDir := strings.TrimSuffix(strings.TrimSuffix(input, "\\"), "/")
			outputDir := inputDir + RSSplitterSuffix
			if len(os.Args) >= 4 {
				outputDir = strings.TrimSuffix(strings.TrimSuffix(os.Args[3], "\\"), "/")
			}
			fmt.Println("inputDir:", inputDir)
			fmt.Println("outputDir:", outputDir)
			if err := rs_splitter.RSSplitterEncodeDir(inputDir, outputDir); err != nil {
				fmt.Printf("EncodeDir error: %v\n", err)
				os.Exit(1)
			}
		}
	case "decode":
		if !isDir {
			outputFile := strings.ReplaceAll(filepath.Base(input), "-", ".")
			if err := rs_splitter.RSSplitterDecode(input, outputFile); err != nil {
				fmt.Printf("Decode error: %v\n", err)
				os.Exit(1)
			}
		} else {
			//if len(os.Args) < 4 {
			//	OutputCmdHelp()
			//}
			inputDir := strings.TrimSuffix(strings.TrimSuffix(input, "\\"), "/")
			outputDir := strings.ReplaceAll(inputDir, RSSplitterSuffix, "")
			if len(os.Args) >= 4 {
				outputDir = strings.TrimSuffix(strings.TrimSuffix(os.Args[3], "\\"), "/")
			}
			fmt.Println("inputDir:", inputDir)
			fmt.Println("outputDir:", outputDir)
			if err := rs_splitter.RSSplitterDecodeDir(inputDir, outputDir); err != nil {
				fmt.Printf("DecodeDir error: %v\n", err)
				os.Exit(1)
			}
		}
	default:
		fmt.Println("Unknown command. Use encode or decode.")
		os.Exit(1)
	}
}
