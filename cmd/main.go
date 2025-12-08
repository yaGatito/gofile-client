package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

var apiKey string
var folderId string

func init() {
	apiKey = os.Getenv("GOFILE_API_KEY")
	folderId = os.Getenv("GOFILE_ACCOUNT_ID")

	apiKey = ""
	folderId = "c8d17506-8456-4c3b-84d1-ac9bfada0332"

	// if apiKey == "" || accountId == "" {
	// 	log.Fatal("GOFILE_API_KEY and GOFILE_ACCOUNT_ID environment variables must be set")
	// }
}

func main() {
	// err := sentCreateFolderRequest(folderId, "main1")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// err := sentUploadFileRequest(folderId, "./files/file1")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	err := getFileRequest("d8a2458c-dbc6-478c-a339-31be76e83e6e", "file.txt")
	// err := getFileRequest("bfb966d2-21b0-4ca9-b256-fd4943059393", "file1")
	if err != nil {
		log.Fatal(err)
	}
}

func getFileRequest(folderId string, fileName string) error {
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	url := fmt.Sprintf(
		"https://store-eu-par-1.gofile.io/download/web/%s/%s",
		url.PathEscape(folderId),
		url.PathEscape(fileName),
	)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	ct := resp.Header.Get("Content-Type")
	fmt.Println("Content-Type:", ct)
	fmt.Println(resp.Status)
	fmt.Println(string(respBody))

	return nil
}

func sentUploadFileRequest(folder string, filePath string) error {
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	body := &bytes.Buffer{}

	writer := multipart.NewWriter(body)

	err := writer.WriteField("folderId", folder)
	if err != nil {
		log.Fatal(err)
	}

	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		log.Fatal(err)
	}

	fileToSent, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer fileToSent.Close()

	_, err = io.Copy(part, fileToSent)
	if err != nil {
		log.Fatal(err)
	}
	writer.Close()

	req, err := http.NewRequest(http.MethodPost, "https://upload.gofile.io/uploadfile", body)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
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
