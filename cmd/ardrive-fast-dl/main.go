package main

import (
	"fmt"
	"github.com/ERR0RPR0MPT/ardrive-fast-ul-dl/ardrive_fast_dl"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ardrive-fast-dl <folder_id> [threads] [arweave_gateway]")
		os.Exit(1)
	}
	folderID := os.Args[1]
	threads := 64
	arGate := ardrive_fast_dl.ArweaveGateway

	if len(os.Args) >= 3 {
		var err error
		threads, err = strconv.Atoi(os.Args[2])
		if err != nil {
			log.Println("Invalid threads value, using default 64")
		}
	}

	if len(os.Args) >= 4 {
		arGate = os.Args[3]
	}

	startTime := time.Now()

	err := ardrive_fast_dl.FastDL(threads, folderID, arGate)
	if err != nil {
		log.Println("Invalid threads value, using default 64")
	}

	endTime := time.Now()
	log.Printf("Total execution time: %s\n", endTime.Sub(startTime))
	log.Println("All downloads completed!")
}
