package arweave_api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/permadao/goar"
	"github.com/permadao/goar/utils"
	"io"
	"net/http"
)

const (
	DefaultArDriveTurboURL = "https://upload.ardrive.io/v1"
	DefaultUserAgent       = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36"
)

type RespArDriveTurboTxStruct struct {
	Id                  string   `json:"id"`
	Timestamp           int64    `json:"timestamp"`
	Winc                string   `json:"winc"`
	Version             string   `json:"version"`
	DeadlineHeight      int64    `json:"deadlineHeight"`
	DataCaches          []string `json:"dataCaches"`
	FastFinalityIndexes []string `json:"fastFinalityIndexes"`
	Public              string   `json:"public"`
	Signature           string   `json:"signature"`
	Owner               string   `json:"owner"`
}

func RequestArDriveTurboTx(signer *goar.Signer, binary []byte) (RespArDriveTurboTxStruct, error) {
	xNonce := uuid.New().String()
	xSigRaw, err := utils.Sign([]byte(xNonce), signer.PrvKey)
	if err != nil {
		return RespArDriveTurboTxStruct{}, err
	}
	xSig := utils.Base64Encode(xSigRaw)
	xPubKey := utils.Base64Encode(signer.PubKey.N.Bytes())

	// 创建请求
	req, err := http.NewRequest("POST", DefaultArDriveTurboURL+"/tx", bytes.NewBuffer(binary))
	if err != nil {
		return RespArDriveTurboTxStruct{}, err
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("X-Public-Key", xPubKey)
	req.Header.Set("X-Nonce", xNonce)
	req.Header.Set("X-Signature", xSig)
	req.Header.Set("User-Agent", DefaultUserAgent)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Origin", "https://app.ardrive.io")
	req.Header.Set("Referer", "https://app.ardrive.io/")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return RespArDriveTurboTxStruct{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return RespArDriveTurboTxStruct{}, fmt.Errorf("resp status code error: %v", resp.Status)
	}

	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return RespArDriveTurboTxStruct{}, err
	}

	// 解析到对象
	var respStruct RespArDriveTurboTxStruct
	if err := json.Unmarshal(body, &respStruct); err != nil {
		return RespArDriveTurboTxStruct{}, fmt.Errorf("error decoding respStruct: %v", err)
	}

	return respStruct, nil
}
