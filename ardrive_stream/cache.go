package ardrive_stream

import (
	"encoding/json"
	"fmt"
	"github.com/ERR0RPR0MPT/ardrive-fast-ul-dl/rs_splitter"
	"log"
	"os"
	"path/filepath"
	"time"
)

var CachePath = DefaultCachePath

func CacheInit(cacheDirPath string) {
	if cacheDirPath != "" {
		CachePath = cacheDirPath
	}
	GetCachePath()

	// 读取缓存数据
	cm, err := ReadCacheFileMeta(C.FileMetaName)
	if err != nil {
		log.Fatalf("cache init error: %v", err)
	}
	cacheLock.Lock()
	fileCache = cm
	cacheLock.Unlock()

	go CacheCleaner()
	go CacheSaver()
}

func GetCachePath() (string, error) {
	if CachePath == "" {
		if C.CachePath == "" {
			CachePath = DefaultCachePath
		} else {
			CachePath = C.CachePath
		}
	}
	if _, err := os.Stat(CachePath); os.IsNotExist(err) {
		err := os.MkdirAll(CachePath, 0755)
		if err != nil {
			log.Fatal("cache path error")
			return "", err
		}
	}
	return CachePath, nil
}

func GetCacheFolderPath(folderId string) (string, error) {
	folderPath := filepath.Join(CachePath, folderId)
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		err := os.MkdirAll(folderPath, 0755)
		if err != nil {
			log.Fatal("cache folder path error")
			return "", err
		}
	}
	return folderPath, nil
}

func ReadCacheFileMeta(fileMetaName string) (map[string]CachedMeta, error) {
	if C.FileMetaName == "" {
		C.FileMetaName = FileMetaName
	}
	d, err := GetCachePath()
	if err != nil {
		return nil, err
	}
	fd := filepath.Join(d, fileMetaName)

	// 尝试读取文件内容
	data, err := os.ReadFile(fd)
	if err != nil {
		// 如果文件不存在，则创建并写入初始数据
		if os.IsNotExist(err) {
			initialData := make(map[string]CachedMeta)
			// 将初始数据序列化为 JSON
			data, err = json.Marshal(initialData)
			if err != nil {
				return nil, err
			}
			// 创建文件并写入初始数据
			if err := os.WriteFile(fd, data, 0644); err != nil {
				return nil, err
			}
			return initialData, nil
		}
		return nil, err // 处理其他错误
	}

	// 反序列化 JSON 数据到 CachedMeta
	var cachedMeta map[string]CachedMeta
	if err := json.Unmarshal(data, &cachedMeta); err != nil {
		return nil, err
	}

	return cachedMeta, nil
}

func ReadCacheFolderFileInfoData(folderId string) (rs_splitter.FileInfoManifest, error) {
	// 获取缓存路径
	cachePath, err := GetCachePath()
	if err != nil {
		return rs_splitter.FileInfoManifest{}, err
	}

	// 构造要读取的文件信息文件路径
	fileInfoPath := filepath.Join(cachePath, folderId, rs_splitter.DefaultFileInfoJsonName)

	// 读取文件内容
	data, err := os.ReadFile(fileInfoPath)
	if err != nil {
		return rs_splitter.FileInfoManifest{}, err
	}

	// 反序列化 JSON 数据到 FileInfo 列表
	var fileInfos rs_splitter.FileInfoManifest
	if err := json.Unmarshal(data, &fileInfos); err != nil {
		return rs_splitter.FileInfoManifest{}, err
	}

	return fileInfos, nil
}

func SaveCacheFolderFileInfoData(folderId string, fileInfoManifest rs_splitter.FileInfoManifest) error {
	// 获取缓存路径
	cachePath, err := GetCachePath()
	if err != nil {
		return err
	}

	// 构造要保存的文件信息文件路径
	fileInfoPath := filepath.Join(cachePath, folderId, rs_splitter.DefaultFileInfoJsonName)

	// 序列化 FileInfoManifest 为 JSON
	data, err := json.MarshalIndent(fileInfoManifest, "", "  ")
	if err != nil {
		return err
	}

	// 确保文件夹存在
	if err := os.MkdirAll(filepath.Dir(fileInfoPath), os.ModePerm); err != nil {
		return err
	}

	// 写入文件内容
	err = os.WriteFile(fileInfoPath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func DeleteCacheFolderData(folderId string) error {
	// 获取缓存路径
	cachePath, err := GetCachePath()
	if err != nil {
		return err
	}

	// 构造要删除的文件夹路径
	folderPath := filepath.Join(cachePath, folderId)

	// 删除文件夹及其内容
	err = os.RemoveAll(folderPath)
	if err != nil {
		return err
	}

	return nil
}

func ReadCacheFolderBinData(folderId string, index string) ([]byte, error) {
	// 获取缓存路径
	cachePath, err := GetCachePath()
	if err != nil {
		return nil, err
	}

	// 构造要读取的文件路径
	filePath := filepath.Join(cachePath, folderId, index)

	// 读取文件内容
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func SaveCacheFileMeta(cachedMeta map[string]CachedMeta, fileMetaName string) error {
	if C.FileMetaName == "" {
		C.FileMetaName = FileMetaName // 默认文件名
	}

	// 获取缓存路径
	cachePath, err := GetCachePath()
	if err != nil {
		return err
	}

	// 构造要保存的文件路径
	filePath := filepath.Join(cachePath, fileMetaName)

	// 序列化 CachedMeta 为 JSON
	data, err := json.MarshalIndent(cachedMeta, "", "  ")
	if err != nil {
		return err
	}

	// 确保文件夹存在
	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return err
	}

	// 写入文件内容
	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func SaveCacheFolderBinData(folderId string, index string, data []byte) error {
	// 获取缓存路径
	cachePath, err := GetCachePath()
	if err != nil {
		return err
	}

	// 构造要保存的文件路径
	filePath := filepath.Join(cachePath, folderId, index)

	// 确保文件夹存在
	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return err
	}

	// 写入文件内容
	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func CacheCleaner() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		<-ticker.C
		for folderId, cm := range fileCache {
			if time.Now().After(cm.Expiration) {
				DeleteCacheFolderData(folderId)
				cacheLock.Lock()
				delete(fileCache, folderId)
				cacheLock.Unlock()
			}
		}
	}
}

func CacheSaver() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		<-ticker.C
		//每10秒钟保存一次 CachedMeta 到 ./cache/fileMeta.json
		if err := SaveCacheFileMeta(fileCache, C.FileMetaName); err != nil {
			// 处理错误，例如打印日志
			fmt.Println("Error saving cache file meta:", err)
		}
	}
}
