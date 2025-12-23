package main

import (
	"context"
	"io"
	"net/http"

	gfclient "github.com/yaGatito/gofile-client"

	"log"
	"os"
	"time"
)

// func main() {
// 	ctx := context.Background()

// 	client := gfclient.New("your-api-key", nil, nil)
// 	if client == nil {
// 		log.Fatal("Failed to create GoFile client, empty API key?")
// 	}

// 	file, err := os.Open("example.txt")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer file.Close()

// 	resp, err := client.UploadFile(ctx, "root", "example.txt", file)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	log.Printf("File uploaded: %+v", resp)
// }

var apiKey string
var folderId string

func init() {
	apiKey = os.Getenv("GOFILE_API_KEY")
	folderId = os.Getenv("GOFILE_ACCOUNT_ID")

	if apiKey == "" || folderId == "" {
		log.Fatal("GOFILE_API_KEY and GOFILE_ACCOUNT_ID environment variables must be set")
	}
}

func main() {
	// A simple use case example.
	logger := log.New(os.Stdout, "[GOFILE-CLIENT] ", log.Ldate|log.Ltime|log.Lmicroseconds|log.LUTC)
	client, err := gfclient.New(apiKey, &http.Client{Timeout: 15 * time.Second}, logger)
	if err != nil {
		logger.Fatal("failed to create client:", err)
	}
	ctx := context.Background()

	// Creating a "NewFolder" folder.
	createFolderResponse, err := client.CreateFolder(ctx, "root", "NewFolder")
	if err != nil {
		log.Fatal(err)
	}
	logger.Println("Response received:", createFolderResponse)

	// Setting context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Opening file
	file, err := os.Open("./files/file.txt")
	if err != nil {
		log.Fatal(err)
	}

	// Uploading file
	postFileResponse, err := client.UploadFile(ctx, createFolderResponse.Data.Id, "./files/file.txt", file)
	if err != nil {
		log.Fatal(err)
	}
	logger.Println("UploadFileResponseBody received:", postFileResponse)

	// Getting file details in order to download it later
	getFileInfoResponse, err := client.GetFileInfo(ctx, postFileResponse.Data.Id)
	if err != nil {
		log.Fatal(err)
	}
	logger.Println("GetFileInfoResponseBody received:", getFileInfoResponse)

	// Getting file's io.ReadCloser from response
	reader, err := client.DownloadFile(ctx, getFileInfoResponse.Data.ServerSelected, getFileInfoResponse.Data.Id, getFileInfoResponse.Data.Name)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()
	bytes, err := io.ReadAll(reader)

	logger.Println("File get:", len(bytes), "bytes")
	logger.Println("File string content:", string(bytes))
}
