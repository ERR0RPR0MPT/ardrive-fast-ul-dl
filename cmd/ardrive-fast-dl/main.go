package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

type ArDriveEntity struct {
	Name         string `json:"name"`
	DataTxId     string `json:"dataTxId"`
	EntityType   string `json:"entityType"`
	Path         string `json:"path"`
	EntityId     string `json:"entityId"`
	EntityIdPath string `json:"entityIdPath"`
}

func worker(tasks <-chan ArDriveEntity, wg *sync.WaitGroup, folderID string, wgDoneFlag bool) {
	if wgDoneFlag {
		defer wg.Done()
	}
	for file := range tasks {
		err := downloadFile(file, folderID)
		if err != nil {
			fmt.Printf("Error downloading %s: %v\n", file.Name, err)
			worker(tasks, wg, folderID, false)
		}
	}
}

func processFolderRecursive(folderID string) ([]ArDriveEntity, error) {
	log.Println("获取文件夹数据:", folderID)
	var files []ArDriveEntity
	entities, err := listFolder(folderID)
	if err != nil {
		return nil, err
	}

	for _, entity := range entities {
		log.Println("获取到文件路径:", entity.Path)
		if entity.EntityType == "folder" {
			subFiles, err := processFolderRecursive(entity.EntityId)
			if err != nil {
				return nil, err
			}
			files = append(files, subFiles...)
		} else {
			files = append(files, entity)
		}
	}
	return files, nil
}

func listFolder(parentFolderID string) ([]ArDriveEntity, error) {
	sArr := []string{"list-folder", "--parent-folder-id", parentFolderID}
	cmd := exec.Command("ardrive", sArr...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("command failed: %v", err, string(output), sArr)
	}

	var entities []ArDriveEntity
	if err := json.Unmarshal(output, &entities); err != nil {
		return nil, fmt.Errorf("JSON parse error: %v", err)
	}

	return entities, nil
}

func findIndex(name string, parts []string) (bool, int) {
	isFound := false
	foundIndex := -1
	for k, v := range parts {
		if v == name {
			isFound = true
			foundIndex = k
		}
	}
	return isFound, foundIndex
}

func downloadFile(file ArDriveEntity, folderID string) error {
	// 生成本地路径
	cleanPath := strings.TrimPrefix(file.Path, "/")
	pathParts := strings.Split(cleanPath, "/")

	//生成相对路径
	cleanEntityIdPath := strings.TrimPrefix(file.EntityIdPath, "/")
	entityIdPathParts := strings.Split(cleanEntityIdPath, "/")

	isFound, idx := findIndex(folderID, entityIdPathParts)
	if !isFound {
		return fmt.Errorf("findIndex(folderID, entityIdPathParts) error")
	}

	pathParts = pathParts[idx:]
	localPath := filepath.Join(".", filepath.Join(pathParts...))

	// 创建目录
	dir := filepath.Dir(localPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create directory failed: %v", err)
	}

	// 下载文件
	log.Println("下载文件:", file.Name)
	resp, err := http.Get(fmt.Sprintf("https://arweave.net/%s", file.DataTxId))
	if err != nil {
		return fmt.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %d", resp.StatusCode)
	}

	// 创建文件
	outFile, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("file create failed: %v", err)
	}
	defer outFile.Close()

	// 保存内容
	if _, err := io.Copy(outFile, resp.Body); err != nil {
		return fmt.Errorf("write file failed: %v", err)
	}
	//log.Println("文件下载完成:", file.Name)

	log.Println("Successfully downloaded:", localPath)
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ardrive-fast-dl <folder_id> [threads]")
		os.Exit(1)
	}
	folderID := os.Args[1]
	threads := 64
	if len(os.Args) > 2 && len(os.Args) < 3 {
		var err error
		threads, err = strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Usage: ardrive-fast-dl <folder_id> [threads]")
			os.Exit(1)
		}
	}

	startTime := time.Now()

	// 获取所有文件
	files, err := processFolderRecursive(folderID)
	if err != nil {
		fmt.Printf("Error processing folder: %v\n", err)
		os.Exit(1)
	}

	endTime := time.Now()
	elapsed := endTime.Sub(startTime)
	log.Println("\n获取文件用时: ", elapsed, "\n")

	// 创建任务通道
	tasks := make(chan ArDriveEntity, len(files))
	var wg sync.WaitGroup

	// 启动64个worker
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go worker(tasks, &wg, folderID, true)
	}

	// 添加任务
	for _, file := range files {
		tasks <- file
	}
	close(tasks)

	// 等待所有任务完成
	wg.Wait()
	endTime = time.Now()
	elapsed = endTime.Sub(startTime)
	fmt.Printf("执行用时: %s\n", elapsed)
	fmt.Println("All downloads completed!")
}
