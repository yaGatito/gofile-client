package main

import (
	"context"
	gofile "gofile/client"
	"io"
	"log"
	"os"
	"time"
)

var apiKey string
var folderId string

func init() {
	apiKey = os.Getenv("GOFILE_API_KEY")
	folderId = os.Getenv("GOFILE_ACCOUNT_ID")

	apiKey = "saFX3OyrQCoAPSmY8DCjtabEmHMj2jVJ"
	folderId = "c8d17506-8456-4c3b-84d1-ac9bfada0332"

	// if apiKey == "" || accountId == "" {
	// 	log.Fatal("GOFILE_API_KEY and GOFILE_ACCOUNT_ID environment variables must be set")
	// }
}

func main() {
	client, err := gofile.NewClient(apiKey, nil)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// createFolderResponse, err := client.CreateFolder(ctx, "root", "HUESOSESICK")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.Default().Println("Response received:", createFolderResponse)

	// ctx, cancel = context.WithTimeout(context.Background(), 15*time.Second)
	// defer cancel()
	// file, err := os.Open("./files/video.mp4")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// postFileResponse, err := client.UploadFile(ctx, "ee6d30df-92f5-4fe7-b174-165c9b838efb", "./files/video.mp4", file)
	// log.Default().Println("Response received:", postFileResponse)

	reader, err := client.DownloadFile(ctx, "store5", "8402ba65-6dd4-4ef3-8178-907d4c58b9f3", "video.mp4")
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()
	bytes, err := io.ReadAll(reader)

	log.Default().Println("File get:", len(bytes), "bytes")
}
