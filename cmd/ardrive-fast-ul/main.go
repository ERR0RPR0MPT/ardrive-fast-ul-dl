package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

type FileTask struct {
	LocalPath string
	ParentID  string
}

type CreateFolderResult struct {
	Created []struct {
		Type       string `json:"type"`
		EntityId   string `json:"entityId"`
		EntityName string `json:"entityName"`
	} `json:"created"`
}

func createFolder(name, parentID, wallet string, turboEnable bool) (string, error) {
	t := ""
	if turboEnable {
		t = "--turbo"
	}
	cmd := exec.Command("ardrive", "create-folder",
		"--parent-folder-id", parentID,
		"--folder-name", name,
		"-w", wallet, t)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("command failed: %v\nOutput: %s", err, string(output))
	}

	var result CreateFolderResult
	if err := json.Unmarshal(output, &result); err != nil {
		return "", fmt.Errorf("parse failed: %v\nOutput: %s", err, string(output))
	}

	if len(result.Created) == 0 {
		return "", fmt.Errorf("no folder created")
	}

	return result.Created[0].EntityId, nil
}

func uploadFile(path, parentID, wallet string, turboEnable bool) error {
	t := ""
	if turboEnable {
		t = "--turbo"
	}
	cmd := exec.Command("ardrive", "upload-file",
		"--local-path", path,
		"--parent-folder-id", parentID,
		"-w", wallet, t)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("upload failed: %v\nOutput: %s", err, string(output))
	}
	return nil
}

func main() {
	startTime := time.Now()

	if len(os.Args) < 4 {
		log.Fatal("Usage: ardrive-fast-ul <parent-folder-id> <local-path> <wallet-path> [threads] [turbo_enable]")
	}

	parentFolderID := os.Args[1]
	localPath := os.Args[2]
	walletPath := os.Args[3]
	threads := 32
	turboEnable := true

	if len(os.Args) > 4 {
		t, err := strconv.Atoi(os.Args[4])
		if err == nil {
			threads = t
		}
		log.Println("threads:", threads)
	}

	if len(os.Args) > 5 {
		t, err := strconv.ParseBool(os.Args[5])
		if err == nil {
			turboEnable = t
		}
		log.Println("turboEnable:", turboEnable)
	}

	// 创建目录映射和任务列表
	dirMap := make(map[string]string)
	var tasks []FileTask

	// 处理根目录
	absPath, _ := filepath.Abs(localPath)
	rootName := filepath.Base(absPath)
	rootID, err := createFolder(rootName, parentFolderID, walletPath, turboEnable)
	if err != nil {
		log.Fatal("Create root folder failed:", err)
	}
	dirMap[absPath] = rootID
	log.Printf("Created root folder: %s (%s)\n", rootName, rootID)

	// 遍历目录结构
	err = filepath.Walk(absPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if path == absPath {
			return nil
		}

		if info.IsDir() {
			parentDir := filepath.Dir(path)
			parentID, ok := dirMap[parentDir]
			if !ok {
				return fmt.Errorf("parent directory not found: %s", parentDir)
			}

			dirName := info.Name()
			newID, err := createFolder(dirName, parentID, walletPath, turboEnable)
			if err != nil {
				return err
			}
			dirMap[path] = newID
			log.Printf("Created folder: %s (%s)\n", dirName, newID)
		} else {
			dir := filepath.Dir(path)
			parentID, ok := dirMap[dir]
			if !ok {
				return fmt.Errorf("parent directory not found: %s", dir)
			}
			tasks = append(tasks, FileTask{path, parentID})
		}

		return nil
	})

	if err != nil {
		log.Fatal("Directory walk failed:", err)
	}

	// 准备并发上传
	taskChan := make(chan FileTask, len(tasks))
	var wg sync.WaitGroup

	// 启动worker池
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range taskChan {
				log.Printf("Uploading: %s\n", task.LocalPath)
				var err error
				const maxRetries = 99999
				for retries := 0; retries < maxRetries; retries++ {
					err = uploadFile(task.LocalPath, task.ParentID, walletPath, turboEnable)
					if err == nil {
						log.Printf("Upload success: %s\n", task.LocalPath)
						break
					}
					log.Printf("Upload failed: %s (%v), retrying... (%d/%d)\n", task.LocalPath, err, retries+1, maxRetries)
				}
				if err != nil {
					log.Printf("Upload permanently failed: %s (%v)\n", task.LocalPath, err)
				}
			}
		}()
	}

	// 分发任务
	for _, task := range tasks {
		taskChan <- task
	}
	close(taskChan)

	// 等待完成
	wg.Wait()

	elapsed := time.Since(startTime)
	log.Printf("Total execution time: %s\n", elapsed)
}
