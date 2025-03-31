package arweave_api

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/permadao/goar"
	"github.com/permadao/goar/schema"
	"strconv"
	"time"
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
	return string(jsonData), nil
}

func ArDriveTurboCreateFolder(name, driveId, parentId, wallet string) (string, error) {
	signer, err := goar.NewSignerFromPath(wallet)
	if err != nil {
		return "", err
	}

	b, err := goar.NewBundler(signer)
	if err != nil {
		return "", err
	}

	newFolderUUID := uuid.New().String()

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
		{
			Name: "Drive-Id", Value: driveId,
		},
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

	_, err = RequestArDriveTurboTx(signer, item.Binary)
	if err != nil {
		return "", err
	}

	return newFolderUUID, nil
}
