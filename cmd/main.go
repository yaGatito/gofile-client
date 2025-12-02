package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var apiKey string
var accountId string

func init() {
	apiKey = os.Getenv("GOFILE_API_KEY")
	accountId = os.Getenv("GOFILE_ACCOUNT_ID")
	if apiKey == "" || accountId == "" {
		log.Fatal("GOFILE_API_KEY and GOFILE_ACCOUNT_ID environment variables must be set")
	}
}

func main() {
	err := sentCreateFolderRequest(accountId, "main")
	if err != nil {
		log.Fatal(err)
	}
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

	// buf := bytes.Buffer{}
	// err := json.NewEncoder(&buf).Encode(body)
	// if err != nil {
	// 	return err
	// }

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

	// map[data:map[code:3PXYcB createTime:1.764681296e+09 id:a7026f8b-49e1-4517-9e86-5a2839774a79 modTime:1.764681296e+09 name:main owner:a03efd8c-82fb-477c-9a3a-f24a7b892b23 parentFolder:c8d17506-8456-4c3b-84d1-ac9bfada0332 type:folder] status:ok]
	var respBody map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return err
	}

	fmt.Println(respBody)

	return nil
}
