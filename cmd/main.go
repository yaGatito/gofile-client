package main

import (
	"encoding/json"
	"log"
	"os"
	gofile "remstor/client"
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

	// req, err := client.CreatePostFolderRequest("root", "HUESOSES")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// bytes, _, err := client.Do(req)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// var createFolderResponseData gofile.CreateFolderResponseData
	// err = json.Unmarshal(bytes, &createFolderResponseData)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.Default().Println("Folder create JSON:", string(bytes))

	// folder id 
	req, err := client.CreatePostFileRequest("ee6d30df-92f5-4fe7-b174-165c9b838efb", "./files/video.mp4")
	if err != nil {
		log.Fatal(err)
	}
	bytes, _, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	var uploadFileResponseData gofile.UploadFileResponseData
	err = json.Unmarshal(bytes, &uploadFileResponseData)
	if err != nil {
		log.Fatal(err)
	}
	log.Default().Println("File upload JSON:", string(bytes))

	req, err = client.CreateGetFileRequest(uploadFileResponseData.Data.Servers[0], uploadFileResponseData.Data.Id, uploadFileResponseData.Data.Name)
	if err != nil {
		log.Default().Println("Error: ", err)
	}
	bytes, _, err = client.Do(req)
	if err != nil {
		log.Default().Println("Error: ", err)
	}

	log.Default().Println("File get:", len(bytes), "bytes")
}
