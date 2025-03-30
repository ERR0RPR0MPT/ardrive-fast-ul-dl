package arweave_api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/permadao/goar"
	"github.com/permadao/goar/schema"
	"github.com/permadao/goar/utils"
	"net/http"
	"strconv"
	"time"
)

const (
	DefaultArDriveTurboURL = "https://upload.ardrive.io/v1"
)

type FolderStruct struct {
	Name string `json:"name"`
}

func GetFolderData(name string) (string, error) {
	folder := FolderStruct{Name: name}
	jsonData, err := json.Marshal(folder)
	if err != nil {
		return "", err
	}
	//encoded := base64.StdEncoding.EncodeToString(jsonData)
	return string(jsonData), nil
}

func CreateFolder(name, driveId, parentId, wallet string) (string, error) {
	signer, err := goar.NewSignerFromPath(wallet)
	if err != nil {
		return "", err
	}

	b, err := goar.NewBundler(signer)
	if err != nil {
		return "", err
	}

	newFolderUUID := uuid.New().String()
	xNonce := uuid.New().String()
	xSigRaw, err := utils.Sign([]byte(xNonce), signer.PrvKey)
	if err != nil {
		return "", err
	}
	xSig := utils.Base64Encode(xSigRaw)
	xPubKey := utils.Base64Encode(signer.PubKey.N.Bytes())

	tags := []schema.Tag{
		{
			Name: "Content-Type", Value: "application/json",
		},
		{
			Name: "ArFS", Value: "0.14",
		},
		{
			Name: "Entity-Type", Value: "folder",
		},
		//{
		//	Name: "Drive-Id", Value: driveId,
		//},
		{
			Name: "Folder-Id", Value: newFolderUUID,
		},
		{
			Name: "Parent-Folder-Id", Value: parentId,
		},
		{
			Name: "App-Name", Value: "Arweave-Explorer",
		},
		{
			Name: "App-Platform", Value: "CLI",
		},
		{
			Name: "App-Version", Value: "0.0.1",
		},
		{
			Name: "Unix-Time", Value: strconv.FormatInt(time.Now().Unix(), 10),
		},
	}

	folderData, _ := GetFolderData(name)

	item, err := b.CreateAndSignItem([]byte(folderData), "", "", tags)
	if err != nil {
		return "", err
	}

	// 创建请求
	req, err := http.NewRequest("POST", DefaultArDriveTurboURL+"/tx", bytes.NewBuffer(item.Binary))
	if err != nil {
		return "", err
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("X-Public-Key", xPubKey)
	req.Header.Set("X-Nonce", xNonce)
	req.Header.Set("X-Signature", xSig)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Origin", "https://app.ardrive.io")
	req.Header.Set("Referer", "https://app.ardrive.io/")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br, zstd")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,ru;q=0.7,zh-TW;q=0.6")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("resp status code error: %v", resp.Status)
	}

	return newFolderUUID, nil
}
