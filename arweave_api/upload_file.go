package arweave_api

import (
	"github.com/permadao/goar"
	"github.com/permadao/goar/schema"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func ArDriveTurboUploadFile(path, wallet string) (string, error) {
	name := filepath.Base(path)
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return ArDriveTurboUploadData(data, name, wallet)
}

func ArDriveTurboUploadData(data []byte, name, wallet string) (string, error) {
	signer, err := goar.NewSignerFromPath(wallet)
	if err != nil {
		return "", err
	}

	b, err := goar.NewBundler(signer)
	if err != nil {
		return "", err
	}

	tags := []schema.Tag{
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
		{
			Name: "Content-Type", Value: getMimeType(name),
		},
	}

	item, err := b.CreateAndSignItem(data, "", "", tags)
	if err != nil {
		return "", err
	}

	resp, err := RequestArDriveTurboTx(signer, item.Binary)
	if err != nil {
		return "", err
	}

	return resp.Id, nil
}
