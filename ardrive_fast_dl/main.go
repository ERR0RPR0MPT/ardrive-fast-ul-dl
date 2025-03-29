package ardrive_fast_dl

import (
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/sync/errgroup"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	ArweaveGateway = "arweave.net"
	//ArweaveGateway = "permagate.io"
	maxRetries = 99999
	baseDelay  = 0 * time.Second
)

type ArDriveEntity struct {
	Name         string `json:"name"`
	DataTxId     string `json:"dataTxId"`
	EntityType   string `json:"entityType"`
	Path         string `json:"path"`
	EntityId     string `json:"entityId"`
	EntityIdPath string `json:"entityIdPath"`
}

func worker(tasks <-chan ArDriveEntity, wg *sync.WaitGroup, folderID, arGate string) {
	defer wg.Done()
	for file := range tasks {
		for {
			err := DownloadFile(file, folderID, arGate)
			if err != nil {
				log.Printf("Error downloading %s: %v\n", file.Name, err)
				continue
			}
			break
		}
	}
}

func ProcessFolderRecursive(folderID string, tasks chan<- ArDriveEntity) error {
	//retryDelay := baseDelay
	log.Println("Processing folder:", folderID)
	entities, err := listFolder(folderID)
	if err != nil {
		return err
	}

	g, ctx := errgroup.WithContext(context.Background())

	for _, entity := range entities {
		entity := entity // Capture loop variable
		if entity.EntityType == "folder" {
			g.Go(func() error {
				retryDelay := baseDelay
				for i := 0; i < maxRetries; i++ {
					err := ProcessFolderRecursive(entity.EntityId, tasks)
					if err == nil {
						return nil
					}
					log.Printf("Subfolder processing attempt %d/%d failed: %v", i+1, maxRetries, err)
					if i < maxRetries-1 {
						time.Sleep(retryDelay)
						retryDelay *= 2
					}
				}
				return fmt.Errorf("failed to process subfolder %s after %d attempts", entity.EntityId, maxRetries)
			})
		} else {
			select {
			case tasks <- entity:
				log.Println("Found entity path:", entity.Path)
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}

func extractBetweenBraces(s string) string {
	start := strings.Index(s, "[")
	if start == -1 {
		return ""
	}
	end := strings.LastIndex(s, "]")
	if end == -1 || end < start {
		return ""
	}
	return s[start : end+1]
}

func listFolder(parentFolderID string) ([]ArDriveEntity, error) {
	sArr := []string{"list-folder", "--parent-folder-id", parentFolderID}
	cmd := exec.Command("ardrive", sArr...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("command failed: %v, output: %s", err, string(output))
	}

	output = []byte(extractBetweenBraces(string(output)))

	var entities []ArDriveEntity
	if err := json.Unmarshal(output, &entities); err != nil {
		return nil, fmt.Errorf("JSON parse error: %v, output: %s", err, string(output))
	}

	return entities, nil
}

func findIndex(name string, parts []string) (bool, int) {
	for k, v := range parts {
		if v == name {
			return true, k
		}
	}
	return false, -1
}

func DownloadFile(file ArDriveEntity, folderID, arGate string) error {
	cleanPath := strings.TrimPrefix(file.Path, "/")
	pathParts := strings.Split(cleanPath, "/")

	cleanEntityIdPath := strings.TrimPrefix(file.EntityIdPath, "/")
	entityIdPathParts := strings.Split(cleanEntityIdPath, "/")

	isFound, idx := findIndex(folderID, entityIdPathParts)
	if !isFound {
		return fmt.Errorf("folderID not found in entityIdPath")
	}

	pathParts = pathParts[idx:]
	localPath := filepath.Join(".", filepath.Join(pathParts...))

	dir := filepath.Dir(localPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	log.Println("Downloading file:", file.Name)
	resp, err := http.Get(fmt.Sprintf("https://%s/%s", arGate, file.DataTxId))
	if err != nil {
		return fmt.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %d", resp.StatusCode)
	}

	outFile, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer outFile.Close()

	if _, err := io.Copy(outFile, resp.Body); err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	log.Println("Successfully downloaded:", localPath)
	return nil
}

func FastDL(threads int, folderID, arGate string) error {
	tasks := make(chan ArDriveEntity, 9999)
	processErrCh := make(chan error, 1)
	var wg sync.WaitGroup

	// Start worker goroutines
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go worker(tasks, &wg, folderID, arGate)
	}

	// Process folders and handle errors
	go func() {
		defer close(tasks)
		processErrCh <- ProcessFolderRecursive(folderID, tasks)
	}()

	// Wait for folder processing to complete and get error
	processErr := <-processErrCh

	// Wait for all workers to finish processing remaining tasks
	wg.Wait()

	if processErr != nil {
		return processErr
	}
	return nil
}
