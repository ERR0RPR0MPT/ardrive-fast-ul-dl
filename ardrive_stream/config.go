package ardrive_stream

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

var C Config

type Config struct {
	Debug              bool     `json:"debug"`
	Host               string   `json:"host"`
	Port               int      `json:"port"`
	CachePath          string   `json:"cache-path"`
	FileMetaName       string   `json:"file-meta-name"`
	ArweaveGateway     []string `json:"arweave-gateway"`
	MaxRetries         int      `json:"max-retries"`
	MaxConcurrency     int      `json:"max-concurrency"`
	CacheTTL           int      `json:"cache-ttl"`
	Timeout            int      `json:"timeout"`
	MultiProcessFolder bool     `json:"multi-process-folder"`
}

func ConfigInit() error {
	// 默认配置
	defaultConfig := Config{
		Debug:              false,
		Host:               "",
		Port:               12888,
		CachePath:          "./cache",
		FileMetaName:       "fileMeta.json",
		ArweaveGateway:     []string{"arweave.net"},
		MaxRetries:         99999999,
		MaxConcurrency:     64,
		CacheTTL:           1440,
		Timeout:            20,
		MultiProcessFolder: false,
	}

	// 检查命令行参数
	var configPath string
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	} else {
		configPath = DefaultConfigName
	}

	// 检查配置文件是否存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// 文件不存在，创建默认配置
		file, err := json.MarshalIndent(defaultConfig, "", "  ")
		if err != nil {
			return fmt.Errorf("error creating default config: %v", err)
		}
		if err := os.WriteFile(configPath, file, 0644); err != nil {
			return fmt.Errorf("error writing config file: %v", err)
		}
		fmt.Printf("Default config created at %s\n", configPath)
		time.Sleep(99999 * time.Hour)
		os.Exit(0)
	} else {
		// 文件存在，读取配置
		data, err := os.ReadFile(configPath)
		if err != nil {
			return fmt.Errorf("error reading config file: %v", err)
		}
		if err := json.Unmarshal(data, &C); err != nil {
			return fmt.Errorf("error decoding config: %v", err)
		}
	}

	return nil
}
