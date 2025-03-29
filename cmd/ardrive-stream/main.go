package main

import (
	"fmt"
	"github.com/ERR0RPR0MPT/ardrive-fast-ul-dl/ardrive_stream"
)

func main() {
	if err := ardrive_stream.ConfigInit(); err != nil {
		fmt.Println("Error initializing config:", err)
		return
	}
	ardrive_stream.CacheInit("")
	ardrive_stream.RunServer()
}
