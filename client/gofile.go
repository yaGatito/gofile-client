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

// TODO: logger support
// TODO: context.Context() support
// TODO: check header data format and if it is HTML - ignore body and close immediately

const (
	postFolderEndpoint   = "https://api.gofile.io/contents/createFolder"
	postFileEndpoint     = "https://upload.gofile.io/uploadfile"
	accountsEndpointPart = "https://api.gofile.io/accounts/"
	contentsEndpointPart = "https://api.gofile.io/contents/"

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

	// Get account ID
	req, err := gfclient.CreateGetIdRequest()
	if err != nil {
		return nil, fmt.Errorf("creating 'getid' request: %w", err)
	}
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

	// Get root folder ID
	req, err = gfclient.CreateGetAccountInfoRequest()
	if err != nil {
		return nil, fmt.Errorf("creating 'getAccountInfo' request: %w", err)
	}
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

	// // Get contents info to validate account
	// req, err = gfclient.CreateGetContentsInfoRequest()
	// if err != nil {
	// 	return nil, fmt.Errorf("creating 'getContentsInfo' request: %w", err)
	// }
	// body, _, err = gfclient.Do(req)
	// if err != nil {
	// 	return nil, fmt.Errorf("sending 'getContentsInfo' request: %w", err)
	// }
	// var getContentsInfoResp getContentsInfoResponseData
	// err = json.Unmarshal(body, &getContentsInfoResp)
	// if err != nil {
	// 	return nil, fmt.Errorf("unmarshalling 'getContentsInfo' response: %w", err)
	// }
	// fmt.Println(getContentsInfoResp)

	// Validate received data
	if gfclient.rootFolderId == "" || gfclient.accountId == "" {
		return nil, fmt.Errorf("invalid account data received")
	}

	return gfclient, nil
}

func (c *GofileClient) Do(req *http.Request) ([]byte, *http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	log.Default().Printf("Sending request to %s", req.URL.String())

	resp, err := c.client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	if resp.Header.Get("Content-Type") == "text/html" {
		log.Default().Printf("Received HTML response body")
		return nil, resp, fmt.Errorf("received HTML response, possible error page")
	}

	if resp.StatusCode >= 400 {
		log.Default().Printf("Error response body: %s", string(respBody))
		return nil, resp, fmt.Errorf("received bad status: %s", resp.Status)
	}

	return respBody, resp, nil
}

func (c *GofileClient) CreatePostFolderRequest(parentFolderId, folderName string) (*http.Request, error) {
	if parentFolderId == "root" {
		parentFolderId = c.rootFolderId
	}
	jsonBody, err := json.Marshal(createFolderRequestBody{
		ParentFolderId: parentFolderId,
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

func (c *GofileClient) CreatePostFileRequest(folderId, filePath string) (*http.Request, error) {
	body := &bytes.Buffer{}

	writer := multipart.NewWriter(body)
	err := writer.WriteField("folderId", folderId)
	if err != nil {
		return nil, fmt.Errorf("writing folder id: %w", err)
	}
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return nil, fmt.Errorf("create form file: %w", err)
	}
	fileToSent, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}
	defer fileToSent.Close()
	_, err = io.Copy(part, fileToSent)
	if err != nil {
		return nil, fmt.Errorf("copying file: %w", err)
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, postFileEndpoint, body)
	if err != nil {
		return nil, fmt.Errorf("creating post file request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	log.Default().Printf("Created file upload request for file %s to folder id %s", filePath, folderId)

	return req, nil
}

const getFileEndpoint = "https://%s.gofile.io/download/web/%s/%s"

func (c *GofileClient) CreateGetFileRequest(server, fileId, fileName string) (*http.Request, error) {
	url := fmt.Sprintf(getFileEndpoint, server, url.PathEscape(fileId), url.PathEscape(fileName))

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Cookie", "accountToken="+c.apiKey)

	return req, nil
}

func (c *GofileClient) CreateGetIdRequest() (*http.Request, error) {
	req, err := http.NewRequest(http.MethodGet, accountsEndpointPart+"getid", nil)
	if err != nil {
		return nil, fmt.Errorf("creating 'getid' request: %w", err)
	}

	return req, nil
}

func (c *GofileClient) CreateGetAccountInfoRequest() (*http.Request, error) {
	req, err := http.NewRequest(http.MethodGet, accountsEndpointPart+c.accountId, nil)
	if err != nil {
		return nil, fmt.Errorf("creating 'getAccountInfo' request: %w", err)
	}

	return req, nil
}

// func (c *GofileClient) CreateGetContentsInfoRequest() (*http.Request, error) {
// 	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s%s%s", contentsEndpointPart, c.rootFolderId, contentsQueryWT), nil)
// 	if err != nil {
// 		return nil, fmt.Errorf("creating 'getContentsInfo' request: %w", err)
// 	}

// 	return req, nil
// }
