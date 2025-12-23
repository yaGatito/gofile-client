package main

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/yaGatito/gofile-client"
)

func main() {
	uploadUsecase()
}

func uploadUsecase() {
	ctx := context.Background()

	client, err := gofile.New("your-api-key", nil, nil)
	if err != nil {
		log.Fatal("Failed to create client:", err)
	}

	fileName := "sample.txt"
	file, err := os.Open("./" + fileName)
	if err != nil {
		log.Fatal("Failed to open file:", err)
	}

	_, err = client.UploadFile(ctx, "your-folder-id", fileName, file)
	if err != nil {
		log.Fatal("Failed to upload file:", err)
	}
}

func uploadAndGetUsecase() {
	ctx := context.Background()

	client, err := gofile.New("your-api-key", nil, nil)
	if err != nil {
		log.Fatal("Failed to create client:", err)
	}

	fileName := "sample.txt"
	file, err := os.Open("./" + fileName)
	if err != nil {
		log.Fatal("Failed to open file:", err)
	}

	uploadFile, err := client.UploadFile(ctx, "your-folder-id", fileName, file)
	if err != nil {
		log.Fatal("Failed to upload file:", err)
	}

	// Reusing file ID in order to download a file
	getFileInfoResponse, err := client.GetFileInfo(ctx, "4fd6sg89d7s6", uploadFile.Data.Id)
	if err != nil {
		log.Fatal("Failed to retrieve file details:", err)
	}

	reader, err := client.DownloadFile(ctx, getFileInfoResponse.Data.ServerSelected, getFileInfoResponse.Data.Id, getFileInfoResponse.Data.Name)
	if err != nil {
		log.Fatal("Failed to download file:", err)
	}
	defer reader.Close()
}
