package gofile

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
)

// TODO: context.Context() support
// TODO: check header data format and if it is HTML - ignore body and close immediately

const (
	postFolderEndpoint = "https://api.gofile.io/contents/createFolder"
	postFileEndpoint   = "https://upload.gofile.io/uploadfile"
	accountsEndpoint   = "https://api.gofile.io/accounts/"
	contentsEndpoint   = "https://api.gofile.io/contents/"

	// TODO: made it dynamic switching according to info provided by CreateGetContentsInfoRequest()
	getFileEndpoint = "https://cold-na-phx-4.gofile.io/download/web/"

	// TODO: remove when API service is fixed
	contentsQueryWT = "?wt=4fd6sg89d7s6"
)

type GofileClient struct {
	apiKey       string
	client       *http.Client
	accountId    string
	rootFolderId string
}

func NewClient(apiKey string, client *http.Client) (*GofileClient, error) {
	if client == nil {
		client = &http.Client{}
	}

	gfclient := &GofileClient{
		apiKey: apiKey,
		client: client,
	}

	req, err := gfclient.CreateGetIdRequest()

	body, _, err := gfclient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("sending 'getid' request: %w", err)
	}

	var getIdResp getIdResponseData
	err = json.Unmarshal(body, &getIdResp)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling 'getid' response: %w", err)
	}

	gfclient.accountId = getIdResp.Data.Id

	req, err = gfclient.CreateGetAccountInfoRequest()
	body, _, err = gfclient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("sending 'getAccountInfo' request: %w", err)
	}

	var getAccountInfoResp getAccountInfoResponseData
	err = json.Unmarshal(body, &getAccountInfoResp)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling 'getAccountInfo' response: %w", err)
	}

	gfclient.rootFolderId = getAccountInfoResp.Data.RootFolder

	if gfclient.rootFolderId == "" || gfclient.accountId == "" {
		return nil, fmt.Errorf("invalid account data received")
	}

	return gfclient, nil
}

func (c *GofileClient) Do(req *http.Request) ([]byte, *http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode >= 400 {
		return nil, resp, fmt.Errorf("received bad status: %s", resp.Status)
	}

	return respBody, resp, nil
}

type createFolderRequestBody struct {
	ParentFolderId string `json:"parentFolderId"`
	FolderName     string `json:"folderName,omitempty"`
}

func (c *GofileClient) CreatePostFolderRequest(accountId, folderName string) (*http.Request, error) {
	jsonBody, err := json.Marshal(createFolderRequestBody{
		ParentFolderId: accountId,
		FolderName:     folderName,
	})
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodPost, postFolderEndpoint, bytes.NewReader(jsonBody))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

func (c *GofileClient) CreatePostFileRequest(folderName, filePath string) (*http.Request, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	defer writer.Close()

	err := writer.WriteField("folderId", folderName)
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

	req, err := http.NewRequest(http.MethodPost, postFileEndpoint, body)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	return req, nil
}

func (c *GofileClient) CreateGetFileRequest(folderId, fileName string) (*http.Request, error) {
	url := fmt.Sprintf("%s%s/%s", getFileEndpoint, url.PathEscape(folderId), url.PathEscape(fileName))

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	return req, nil
}

type getIdResponseData struct {
	Status string `json:"status"`
	Data   struct {
		Id    string `json:"id"`
		Tier  string `json:"tier"`
		Email string `json:"email"`
	} `json:"data"`
}

func (c *GofileClient) CreateGetIdRequest() (*http.Request, error) {
	req, err := http.NewRequest(http.MethodGet, accountsEndpoint+"getid", nil)
	if err != nil {
		return nil, fmt.Errorf("creating 'getid' request: %w", err)
	}

	return req, nil
}

type getAccountInfoResponseData struct {
	Status string `json:"status"`
	Data   struct {
		RootFolder string `json:"rootFolder"`
		Stats      struct {
			FolderCount int `json:"folderCount"`
			FileCount   int `json:"fileCount"`
			Storage     int `json:"storage"`
		} `json:"statsCurrent"`
		Email string `json:"email"`
	} `json:"data"`
}

func (c *GofileClient) CreateGetAccountInfoRequest() (*http.Request, error) {
	req, err := http.NewRequest(http.MethodGet, accountsEndpoint+c.accountId, nil)
	if err != nil {
		return nil, fmt.Errorf("creating 'getAccountInfo' request: %w", err)
	}

	return req, nil
}

func (c *GofileClient) CreateGetContentsInfoRequest() (*http.Request, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s%s%s", contentsEndpoint, c.rootFolderId, contentsQueryWT), nil)
	if err != nil {
		return nil, fmt.Errorf("creating 'getContentsInfo' request: %w", err)
	}

	return req, nil
}
