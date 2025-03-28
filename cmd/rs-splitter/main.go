package main

import (
	"fmt"
	"github.com/ERR0RPR0MPT/ardrive-fast-ul-dl/rs_splitter"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: rs-splitter encode <input> [dataShards] [parityShards] [chunkSize] | decode <input>")
		os.Exit(1)
	}

	cmd := os.Args[1]
	input := os.Args[2]

	dataShards, parityShards, chunkSize, err := rs_splitter.GetDefaultShardsNum(input)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

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

	switch cmd {
	case "encode":
		fmt.Println("dataShards:", dataShards)
		fmt.Println("parityShards:", parityShards)
		fmt.Println("chunkSize:", chunkSize)
		outputDir := strings.ReplaceAll(filepath.Base(input), ".", "-")
		if err := rs_splitter.RSSplitterEncode(input, outputDir, dataShards, parityShards, chunkSize); err != nil {
			fmt.Printf("Encode error: %v\n", err)
			os.Exit(1)
		}
	case "decode":
		outputFile := strings.ReplaceAll(filepath.Base(input), "-", ".")
		if err := rs_splitter.RSSplitterDecode(input, outputFile); err != nil {
			fmt.Printf("Decode error: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Println("Unknown command. Use encode or decode.")
		os.Exit(1)
	}
}
