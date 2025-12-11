package gofile

import (
	"bytes"
	"context"
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

// TODO: logger support
// TODO: context.Context() support

const (
	postFolderEndpoint   = "https://api.gofile.io/contents/createFolder"
	postFileEndpoint     = "https://upload.gofile.io/uploadfile"
	accountsEndpointPart = "https://api.gofile.io/accounts/"
	contentsEndpointPart = "https://api.gofile.io/contents/"

	contentTypeHeader = "Content-Type"

	folderIdAttribute = "folderId"
	fileAttribute     = "file"

	rootFolderIdPlaceholder = "root"
)

type GofileClient struct {
	apiKey       string
	client       *http.Client
	accountId    string
	rootFolderId string
}

// NewClient creates client with provided API key. It will create a new client if the provided one is.
func NewClient(apiKey string, client *http.Client) (*GofileClient, error) {
	if client == nil {
		client = &http.Client{}
	}

	gfclient := &GofileClient{
		apiKey: apiKey,
		client: client,
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	// Get account ID
	req, err := gfclient.CreateGetIdRequest(ctx)
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
	req, err = gfclient.CreateGetAccountInfoRequest(ctx)
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
		return nil, nil, fmt.Errorf("sending request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("reading response body: %w", err)
	}

	// Check error responses
	if resp.Header.Get(contentTypeHeader) == "text/html" || bytes.HasPrefix(respBody, []byte("<!DOCTYPE html>")) {
		log.Default().Printf("Received HTML response body")
		return nil, resp, fmt.Errorf("received HTML response, possible error page")
	}

	if resp.StatusCode >= 400 {
		log.Default().Printf("Error response body: %s", string(respBody))
		return nil, resp, fmt.Errorf("received bad status: %s", resp.Status)
	}

	return respBody, resp, nil
}

func (c *GofileClient) CreatePostFolderRequest(ctx context.Context, parentFolderId, folderName string) (*http.Request, error) {
	if parentFolderId == rootFolderIdPlaceholder {
		parentFolderId = c.rootFolderId
	}
	jsonBody, err := json.Marshal(createFolderRequestBody{
		ParentFolderId: parentFolderId,
		FolderName:     folderName,
	})
	if err != nil {
		return nil, fmt.Errorf("creating marshalling response: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, postFolderEndpoint, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("creating post folder request: %w", err)
	}
	req.Header.Set(contentTypeHeader, "application/json")

	return req, nil
}

func (c *GofileClient) CreatePostFileRequest(ctx context.Context, folderId, filePath string) (*http.Request, error) {
	body := &bytes.Buffer{}

	writer := multipart.NewWriter(body)
	err := writer.WriteField(folderIdAttribute, folderId)
	if err != nil {
		return nil, fmt.Errorf("writing folder id: %w", err)
	}
	part, err := writer.CreateFormFile(fileAttribute, filepath.Base(filePath))
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

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, postFileEndpoint, body)
	if err != nil {
		return nil, fmt.Errorf("creating post file request: %w", err)
	}
	req.Header.Set(contentTypeHeader, writer.FormDataContentType())

	log.Default().Printf("Created file upload request for file %s to folder id %s", filePath, folderId)

	return req, nil
}

const getFileEndpoint = "https://%s.gofile.io/download/web/%s/%s"

func (c *GofileClient) CreateGetFileRequest(ctx context.Context, server, fileId, fileName string) (*http.Request, error) {
	url := fmt.Sprintf(getFileEndpoint, server, url.PathEscape(fileId), url.PathEscape(fileName))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating 'getFile' request: %w", err)
	}

	return req, nil
}

func (c *GofileClient) CreateGetIdRequest(ctx context.Context) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, accountsEndpointPart+"getid", nil)
	if err != nil {
		return nil, fmt.Errorf("creating 'getid' request: %w", err)
	}

	return req, nil
}

func (c *GofileClient) CreateGetAccountInfoRequest(ctx context.Context) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, accountsEndpointPart+c.accountId, nil)
	if err != nil {
		return nil, fmt.Errorf("creating 'getAccountInfo' request: %w", err)
	}

	return req, nil
}
