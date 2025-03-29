package ardrive_stream

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ERR0RPR0MPT/ardrive-fast-ul-dl/ardrive_fast_dl"
	"github.com/ERR0RPR0MPT/ardrive-fast-ul-dl/rs_splitter"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	ArweaveGateway = "arweave.net"
	maxRetries     = 9999999
	maxConcurrency = 64
	rangePrefix    = "bytes="
	cacheTTL       = 60 * time.Minute
	//shardIDPrefixLen = 5 // 00001.bin
)

type FileMeta struct {
	DataShards   int
	ParityShards int
	ChunkSize    int
	OriginalSize int64
	Filename     string
	MIMEType     string
	ShardMap     map[int]ardrive_fast_dl.ArDriveEntity
	LastAccessed time.Time
}

type CachedMeta struct {
	meta       *FileMeta
	expiration time.Time
}

var (
	fileCache sync.Map
	cacheLock sync.Mutex
	client    = &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: maxConcurrency,
		},
		Timeout: 30 * time.Second,
	}
)

func HandleFileRequest(c *gin.Context) {
	folderID := c.Param("folderId")

	// Handle range request
	rangeHeader := c.GetHeader("Range")
	if rangeHeader == "" {
		log.Println("send full file")
	} else {
		log.Println("send partial content:", rangeHeader)
	}

	meta, err := getFileMeta(folderID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Println("DataShards:", meta.DataShards)
	log.Println("ParityShards:", meta.ParityShards)
	log.Println("ChunkSize:", meta.ChunkSize)
	log.Println("OriginalSize:", meta.OriginalSize)
	log.Println("MIMEType:", meta.MIMEType)
	log.Println("LastAccessed:", meta.LastAccessed)

	// Set response headers
	c.Header("Content-Type", meta.MIMEType)
	c.Header("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", meta.Filename))
	c.Header("Accept-Ranges", "bytes")

	// Handle range request
	rangeHeader = c.GetHeader("Range")
	if rangeHeader == "" {
		sendFullFile(c, meta)
		return
	}

	start, end, err := parseRange(rangeHeader, meta.OriginalSize)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	sendPartialContent(c, meta, start, end)
}

func sendFullFile(c *gin.Context, meta *FileMeta) {
	c.Header("Content-Length", strconv.FormatInt(meta.OriginalSize, 10))

	// Create pipe for streaming
	pr, pw := io.Pipe()
	defer pr.Close()

	go func() {
		defer pw.Close()
		if err := streamData(pw, meta, 0, meta.OriginalSize-1); err != nil {
			log.Printf("Stream error: %v", err)
		}
	}()

	c.Stream(func(w io.Writer) bool {
		io.Copy(w, pr)
		return false
	})
}

func sendPartialContent(c *gin.Context, meta *FileMeta, start, end int64) {
	log.Println("sendPartialContent:", start, end)
	contentLength := end - start + 1
	c.Header("Content-Length", strconv.FormatInt(contentLength, 10))
	c.Header("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, meta.OriginalSize))
	c.Status(http.StatusPartialContent)

	// Create pipe for streaming
	pr, pw := io.Pipe()
	defer pr.Close()

	go func() {
		defer pw.Close()
		if err := streamData(pw, meta, start, end); err != nil {
			log.Printf("Stream error: %v", err)
		}
	}()

	c.Stream(func(w io.Writer) bool {
		io.Copy(w, pr)
		return false
	})
}

type fetchResult struct {
	blockIndex int64
	data       []byte
	err        error
}

func streamData(w io.Writer, meta *FileMeta, start, end int64) error {
	return streamDataMulti(w, meta, start, end)
}

func streamDataSync(w io.Writer, meta *FileMeta, start, end int64) error {
	blockSize := int64(meta.DataShards * meta.ChunkSize)

	current := start
	for current <= end {
		blockIndex := current / blockSize
		blockStart := blockIndex * blockSize
		blockEnd := blockStart + blockSize - 1
		if blockEnd > meta.OriginalSize-1 {
			blockEnd = meta.OriginalSize - 1
		}

		readStart := max(current, blockStart)
		readEnd := min(end, blockEnd)

		// Synchronously fetch the block
		data, err := fetchBlock(context.Background(), meta, blockIndex)
		if err != nil {
			return err
		}

		blockOffset := readStart - blockStart
		length := readEnd - readStart + 1
		if _, err := w.Write(data[blockOffset : blockOffset+length]); err != nil {
			return err
		}

		current = readEnd + 1
	}

	return nil
}

func streamDataMulti(w io.Writer, meta *FileMeta, start, end int64) error {
	blockSize := meta.DataShards * meta.ChunkSize
	startBlockIndex := start / int64(blockSize)
	endBlockIndex := end / int64(blockSize)
	sem := make(chan struct{}, maxConcurrency)

	resultChan := make(chan fetchResult)
	writeErrChan := make(chan error, 1)
	var wg sync.WaitGroup

	// 启动写入器协程
	go func() {
		cache := make(map[int64][]byte)
		currentBlockIndex := startBlockIndex

		for res := range resultChan {
			if res.err != nil {
				writeErrChan <- res.err
				return
			}

			blockIndex := res.blockIndex
			data := res.data

			if blockIndex == currentBlockIndex {
				if err := writeBlock(w, meta, start, end, blockIndex, data); err != nil {
					writeErrChan <- err
					return
				}
				currentBlockIndex++
				// 检查缓存中的后续块
				for {
					cachedData, ok := cache[currentBlockIndex]
					if !ok {
						break
					}
					if err := writeBlock(w, meta, start, end, currentBlockIndex, cachedData); err != nil {
						writeErrChan <- err
						return
					}
					delete(cache, currentBlockIndex)
					currentBlockIndex++
				}
			} else if blockIndex > currentBlockIndex {
				cache[blockIndex] = data
			}
		}

		// 处理剩余缓存
		for currentBlockIndex <= endBlockIndex {
			cachedData, ok := cache[currentBlockIndex]
			if !ok {
				writeErrChan <- fmt.Errorf("missing block %d", currentBlockIndex)
				return
			}
			if err := writeBlock(w, meta, start, end, currentBlockIndex, cachedData); err != nil {
				writeErrChan <- err
				return
			}
			delete(cache, currentBlockIndex)
			currentBlockIndex++
		}

		writeErrChan <- nil
	}()

	// 启动多个协程异步拉取数据块
	for i := startBlockIndex; i <= endBlockIndex; i++ {
		wg.Add(1)
		go func(bIndex int64) {
			defer wg.Done()
			sem <- struct{}{}        // 获取信号量
			defer func() { <-sem }() // 释放信号量

			data, err := fetchBlock(context.Background(), meta, bIndex)
			if err != nil {
				resultChan <- fetchResult{blockIndex: bIndex, err: err}
				return
			}
			resultChan <- fetchResult{blockIndex: bIndex, data: data}
		}(i)
	}

	// 等待所有拉取完成并关闭通道
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	return <-writeErrChan
}

func writeBlock(w io.Writer, meta *FileMeta, start, end, blockIndex int64, data []byte) error {
	blockSize := meta.DataShards * meta.ChunkSize
	blockStart := blockIndex * int64(blockSize)
	blockEnd := blockStart + int64(blockSize) - 1
	if blockEnd > meta.OriginalSize-1 {
		blockEnd = meta.OriginalSize - 1
	}

	readStart := max(start, blockStart)
	readEnd := min(end, blockEnd)
	if readStart > readEnd {
		return nil // 无需写入此块
	}

	offset := readStart - blockStart
	length := readEnd - readStart + 1
	if int64(len(data)) < offset+length {
		return fmt.Errorf("block %d data too short (expected %d, got %d)", blockIndex, offset+length, len(data))
	}

	if _, err := w.Write(data[offset : offset+length]); err != nil {
		return err
	}
	return nil
}

func fetchBlock(ctx context.Context, meta *FileMeta, blockIndex int64) ([]byte, error) {
	log.Println("fetchBlock:", "blockIndex:", blockIndex)

	shardStart := 1 + blockIndex*int64(meta.DataShards+meta.ParityShards)
	blockResult := make([]byte, meta.DataShards*meta.ChunkSize)

	type result struct {
		index int
		data  []byte
		err   error
	}

	resultChan := make(chan result, meta.DataShards)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for i := 0; i < meta.DataShards; i++ {
		go func(shardNum int) {
			shardID := int(shardStart) + shardNum
			dataTxID, exists := meta.ShardMap[shardID]
			if !exists {
				resultChan <- result{err: fmt.Errorf("missing shard %d", shardID)}
				return
			}

			data, err := downloadWithRetry(ctx, dataTxID.DataTxId)
			if err != nil {
				resultChan <- result{err: err}
				return
			}

			resultChan <- result{
				index: shardNum,
				data:  data,
			}
		}(i)
	}

	received := 0
	var firstError error
	for received < meta.DataShards {
		select {
		case res := <-resultChan:
			if res.err != nil {
				if firstError == nil {
					firstError = res.err
					cancel()
				}
			} else {
				copy(blockResult[res.index*meta.ChunkSize:], res.data)
			}
			received++
		case <-ctx.Done():
			return nil, firstError
		}
	}

	if firstError != nil {
		return nil, firstError
	}
	return blockResult, nil
}

func downloadWithRetry(ctx context.Context, dataTxID string) ([]byte, error) {
	for attempt := 1; ; attempt++ {
		req, _ := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("https://%s/%s", ArweaveGateway, dataTxID), nil)
		resp, err := client.Do(req)
		if err == nil && resp.StatusCode == http.StatusOK {
			data, err := io.ReadAll(resp.Body)
			resp.Body.Close()
			return data, err
		}
		if resp != nil {
			resp.Body.Close()
		}

		if attempt >= maxRetries {
			return nil, fmt.Errorf("max retries reached for %s", dataTxID)
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(time.Duration(attempt) * time.Second):
		}
	}
}

func parseRange(rangeHeader string, fileSize int64) (int64, int64, error) {
	if !strings.HasPrefix(rangeHeader, rangePrefix) {
		return 0, 0, fmt.Errorf("invalid range prefix")
	}

	rangeStr := strings.TrimPrefix(rangeHeader, rangePrefix)
	parts := strings.Split(rangeStr, "-")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid range format")
	}

	start, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, 0, err
	}

	var end int64
	if parts[1] == "" {
		end = fileSize - 1
	} else {
		end, err = strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return 0, 0, err
		}
	}

	if start < 0 || end >= fileSize || start > end {
		return 0, 0, fmt.Errorf("invalid range bounds")
	}

	return start, end, nil
}

func getFileMeta(folderID string) (*FileMeta, error) {
	cacheLock.Lock()
	defer cacheLock.Unlock()

	if cached, ok := fileCache.Load(folderID); ok {
		cm := cached.(CachedMeta)
		if time.Now().Before(cm.expiration) {
			return cm.meta, nil
		}
	}

	meta, err := fetchFileMeta(folderID)
	if err != nil {
		return nil, err
	}

	fileCache.Store(folderID, CachedMeta{
		meta:       meta,
		expiration: time.Now().Add(cacheTTL),
	})

	return meta, nil
}

func fetchFileMeta(folderID string) (*FileMeta, error) {
	// 1. Collect data shards
	shardMap, shardSlice, err := collectAllShards(folderID)
	if err != nil {
		return nil, err
	}

	// 2. Find fileInfo.json
	fileInfo, err := findFileInfo(shardSlice)
	if err != nil {
		return nil, err
	}

	// 2. Parse manifest
	manifest, err := downloadManifest(fileInfo.DataTxId)
	if err != nil {
		return nil, err
	}

	return &FileMeta{
		DataShards:   manifest.DataShards,
		ParityShards: manifest.ParityShards,
		ChunkSize:    manifest.ChunkSize,
		OriginalSize: manifest.OriginalSize,
		Filename:     manifest.Filename,
		MIMEType:     manifest.MIMETypes,
		ShardMap:     shardMap,
	}, nil
}

func findFileInfo(entities []ardrive_fast_dl.ArDriveEntity) (*ardrive_fast_dl.ArDriveEntity, error) {
	isFound := false
	var ade *ardrive_fast_dl.ArDriveEntity
	for _, v := range entities {
		if v.Name == rs_splitter.DefaultFileInfoJsonName {
			ade = &v
			isFound = true
		}
	}
	if !isFound {
		return nil, fmt.Errorf("cannot find fileInfo file")
	}
	return ade, nil
}

func downloadManifest(dataTxID string) (*rs_splitter.FileInfoManifest, error) {
	resp, err := client.Get(fmt.Sprintf("https://%s/%s", ArweaveGateway, dataTxID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var manifest rs_splitter.FileInfoManifest
	if err := json.NewDecoder(resp.Body).Decode(&manifest); err != nil {
		return nil, err
	}

	return &manifest, nil
}

func collectAllShards(folderID string) (map[int]ardrive_fast_dl.ArDriveEntity, []ardrive_fast_dl.ArDriveEntity, error) {
	tasks := make(chan ardrive_fast_dl.ArDriveEntity, 9999)
	processErrCh := make(chan error, 1)
	var wg sync.WaitGroup
	resultMap := make(map[int]ardrive_fast_dl.ArDriveEntity)
	resultSlice := make([]ardrive_fast_dl.ArDriveEntity, 0)

	// Start add goroutines
	wg.Add(1)
	go func(tasks <-chan ardrive_fast_dl.ArDriveEntity, wg *sync.WaitGroup, folderID string) {
		defer wg.Done()
		for file := range tasks {
			resultSlice = append(resultSlice, file)
			isShardFile, indexShard := rs_splitter.CheckShardNumber(file.Name)
			if !isShardFile {
				log.Println(file.Name, "is not shard file")
				continue
			}
			resultMap[indexShard] = file
		}
	}(tasks, &wg, folderID)

	// Process folders and handle errors
	go func() {
		defer close(tasks)
		processErrCh <- ardrive_fast_dl.ProcessFolderRecursive(folderID, tasks)
	}()

	// Wait for folder processing to complete and get error
	processErr := <-processErrCh

	// Wait for all workers to finish processing remaining tasks
	wg.Wait()

	if processErr != nil {
		return nil, nil, processErr
	}

	return resultMap, resultSlice, nil
}

func CacheCleaner() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		<-ticker.C
		fileCache.Range(func(key, value interface{}) bool {
			cm := value.(CachedMeta)
			if time.Now().After(cm.expiration) {
				fileCache.Delete(key)
			}
			return true
		})
	}
}
