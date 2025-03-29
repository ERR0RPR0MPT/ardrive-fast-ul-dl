package main

import (
	"github.com/ERR0RPR0MPT/ardrive-fast-ul-dl/ardrive_stream"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	go ardrive_stream.CacheCleaner()

	r := gin.Default()
	r.GET("/ardrive/file/:folderId", ardrive_stream.HandleFileRequest)
	log.Fatal(r.Run(":8080"))
}
