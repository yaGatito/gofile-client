package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"time"
	
	"github.com/yaGatito/gofile-client"
)

func main() {
	usecase()
}

func usecase() {
	apiKey := os.Getenv("GOFILE_API_KEY")
	if apiKey == "" {
		log.Fatal("GOFILE_API_KEY environment variable must be set")
	}

	// A simple use case example.
	loggerTag := "[GOFILE-CLIENT] "
	logger := log.New(os.Stdout, loggerTag, log.Ldate|log.Ltime|log.Lmicroseconds|log.LUTC)
	client, err := gofile.New(apiKey, &http.Client{Timeout: 15 * time.Second}, logger)
	if err != nil {
		logger.Fatal("Failed to create client:", err)
	}

	// Setting context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Creating a "NewFolder" folder.
	folderName := "NewFolder"
	createFolderResponse, err := client.CreateFolder(ctx, gofile.RootFolder, folderName)
	if err != nil {
		logger.Fatal("Failed to create folder:", err)
	}
	logger.Println("Response received:", createFolderResponse)

	// Opening file
	fileName := "sample.txt"
	file, err := os.Open("./" + fileName)
	if err != nil {
		logger.Fatal("Failed to create folder:", err)
	}

	// Uploading file
	postFileResponse, err := client.UploadFile(ctx, createFolderResponse.Data.Id, fileName, file)
	if err != nil {
		logger.Fatal("Failed to upload file:", err)
	}
	logger.Println("UploadFileResponseBody received:", postFileResponse)

	// Getting file details in order to download it later
	getFileInfoResponse, err := client.GetFileInfo(ctx, postFileResponse.Data.Id)
	if err != nil {
		logger.Fatal("Failed to retrieve file details:", err)
	}
	logger.Println("GetFileInfoResponseBody received:", getFileInfoResponse)

	// Getting file's io.ReadCloser from response
	reader, err := client.DownloadFile(ctx, getFileInfoResponse.Data.ServerSelected, getFileInfoResponse.Data.Id, getFileInfoResponse.Data.Name)
	if err != nil {
		logger.Fatal("Failed to download file:", err)
	}
	defer reader.Close()

	bytes, err := io.ReadAll(reader)
	if err != nil {
		logger.Fatal("Failed to read file from response:", err)
	}

	logger.Println("File get:", len(bytes), "bytes")
	logger.Println("File content:", string(bytes))
}
