package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var apiKey string
var accountId string

func init() {
	apiKey = os.Getenv("GOFILE_API_KEY")
	accountId = os.Getenv("GOFILE_ACCOUNT_ID")

	apiKey = "saFX3OyrQCoAPSmY8DCjtabEmHMj2jVJ"
	accountId = "c8d17506-8456-4c3b-84d1-ac9bfada0332"

	// if apiKey == "" || accountId == "" {
	// 	log.Fatal("GOFILE_API_KEY and GOFILE_ACCOUNT_ID environment variables must be set")
	// }
}

func main() {
	// err := sentCreateFolderRequest(accountId, "main")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	err := sentUploadFileRequest(accountId, "./files/file.txt")
	if err != nil {
		panic(err)
	}
}

func sentUploadFileRequest(folder string, filePath string) error {
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	body := &bytes.Buffer{}

	writer := multipart.NewWriter(body)

	err := writer.WriteField("folderId", folder)
	if err != nil {
		panic(err)
	}

	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		panic(err)
	}

	fileToSent, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer fileToSent.Close()

	_, err = io.Copy(part, fileToSent)
	if err != nil {
		panic(err)
	}
	writer.Close()

	req, err := http.NewRequest(http.MethodPost, "https://upload.gofile.io/uploadfile", body)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	fmt.Println(resp.Status)
	fmt.Println(string(respBody))

	return nil
}

type createFolderRequestBody struct {
	ParentFolderId string `json:"parentFolderId"`
	FolderName     string `json:"folderName,omitempty"`
}

func sentCreateFolderRequest(accountId string, folderName string) error {
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	body := createFolderRequestBody{
		ParentFolderId: accountId,
		FolderName:     folderName,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, "https://api.gofile.io/contents/createFolder", bytes.NewReader(jsonBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	fmt.Println(resp.Status)
	fmt.Println(string(respBody))

	return nil
}
