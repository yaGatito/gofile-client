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
	"strings"
	"sync"
)

const (
	postFolderEndpoint   = "https://api.gofile.io/contents/createFolder"
	postFileEndpoint     = "https://upload.gofile.io/uploadfile"
	accountsEndpointPart = "https://api.gofile.io/accounts/"
	contentsEndpointPart = "https://api.gofile.io/contents/"

	contentTypeHeader = "Content-Type"

	folderIdAttribute = "folderId"
	fileAttribute     = "file"

	rootFolderIdPlaceholderConst = "root"
)

type GofileClient struct {
	apiKey string
	client *http.Client
	logger *log.Logger

	accountIdCached string
	accountIdOnce   sync.Once
	accountIdError  error

	rootFolderIdCached string
	rootFolderIdOnce   sync.Once
	rootFolderIdError  error
}

// NewClient creates client with provided API key. It will create a new client if the provided one is.
func NewClient(apiKey string, client *http.Client, logger *log.Logger) *GofileClient {
	if apiKey == "" {
		return nil
	}
	if client == nil {
		client = &http.Client{}
	}
	if logger == nil {
		logger = log.New(os.Stdout, "[GOFILE-CLIENT] ", log.Ldate|log.Ltime|log.Lmicroseconds|log.LUTC)
	}

	c := &GofileClient{
		apiKey: apiKey,
		client: client,
		logger: logger,
	}

	return c
}

func (c *GofileClient) CreateFolder(ctx context.Context, parentFolderId, newFolderName string) (CreateFolderResponseBody, error) {
	if parentFolderId == "" {
		return CreateFolderResponseBody{}, fmt.Errorf("parentFolderId empty")
	}
	if newFolderName == "" {
		return CreateFolderResponseBody{}, fmt.Errorf("folder name empty")
	}

	var err error
	if parentFolderId == rootFolderIdPlaceholderConst {
		parentFolderId, err = c.rootFolderId(ctx)
		if err != nil {
			return CreateFolderResponseBody{}, err
		}
	}

	req, err := c.createPostFolderRequest(ctx, parentFolderId, newFolderName)
	if err != nil {
		return CreateFolderResponseBody{}, err
	}
	resp, err := c.do(req)
	if err != nil {
		return CreateFolderResponseBody{}, err
	}
	defer resp.Body.Close()

	var ceateFolderResponseBody CreateFolderResponseBody
	err = json.NewDecoder(resp.Body).Decode(&ceateFolderResponseBody)
	if err != nil {
		return CreateFolderResponseBody{}, err
	}
	return ceateFolderResponseBody, nil
}

func (c *GofileClient) UploadFile(ctx context.Context, folderId, fileName string, fileReader io.ReadCloser) (UploadFileResponseBody, error) {
	if folderId == "" {
		return UploadFileResponseBody{}, fmt.Errorf("folderId is not specified")
	}
	if fileName == "" {
		return UploadFileResponseBody{}, fmt.Errorf("fileName is not specified")
	}
	if fileReader == nil {
		return UploadFileResponseBody{}, fmt.Errorf("fileReader is not specified")
	}

	req, err := c.createPostFileRequest(ctx, folderId, fileName, fileReader)
	if err != nil {
		return UploadFileResponseBody{}, err
	}
	resp, err := c.do(req)
	if err != nil {
		return UploadFileResponseBody{}, err
	}
	defer resp.Body.Close()

	var uploadFileResponseBody UploadFileResponseBody
	err = json.NewDecoder(resp.Body).Decode(&uploadFileResponseBody)
	if err != nil {
		return UploadFileResponseBody{}, err
	}
	return uploadFileResponseBody, nil
}

func (c *GofileClient) DownloadFile(ctx context.Context, server, fileId, fileName string) (io.ReadCloser, error) {
	if server == "" {
		return nil, fmt.Errorf("server is not specified")
	}
	if fileId == "" {
		return nil, fmt.Errorf("fileId is not specified")
	}
	if fileName == "" {
		return nil, fmt.Errorf("fileName is not specified")
	}

	req, err := c.createGetFileRequest(ctx, server, fileId, fileName)
	if err != nil {
		return nil, err
	}
	response, err := c.do(req)
	if err != nil {
		return nil, err
	}

	return response.Body, nil
}

func (c *GofileClient) GetFileInfo(ctx context.Context, fileId string) (GetFileInfoResponseBody, error) {
	req, err := c.createGetFileInfoRequest(ctx, fileId)
	if err != nil {
		return GetFileInfoResponseBody{}, err
	}
	resp, err := c.do(req)
	if err != nil {
		return GetFileInfoResponseBody{}, err
	}
	defer resp.Body.Close()

	var getFileInfoResponseBody GetFileInfoResponseBody
	err = json.NewDecoder(resp.Body).Decode(&getFileInfoResponseBody)
	if err != nil {
		return GetFileInfoResponseBody{}, err
	}
	return getFileInfoResponseBody, nil
}

func (c *GofileClient) do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	c.logger.Printf("Sending request to %s\n", req.URL.String())

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("sending request: %w", err)
	}

	// Check error responses
	if strings.HasPrefix(resp.Header.Get(contentTypeHeader), "text/html") {
		c.logger.Printf("Received HTML response body\n")
		return nil, fmt.Errorf("received HTML response, possible error page")
	}

	if resp.StatusCode >= 400 {
		bytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("received bad status: %s, body: %s", resp.Status, string(bytes))
	}

	return resp, nil
}

func (c *GofileClient) createPostFolderRequest(ctx context.Context, parentFolderId, folderName string) (*http.Request, error) {
	if parentFolderId == "" {
		return nil, fmt.Errorf("empty parentFolderId provided")
	}
	if folderName == "" {
		return nil, fmt.Errorf("empty folderName provided")
	}
	var requestBody = createFolderRequestBody{
		ParentFolderId: parentFolderId,
		FolderName:     folderName,
	}
	var err error
	jsonBody, err := json.Marshal(requestBody)
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

func (c *GofileClient) createPostFileRequest(ctx context.Context, folderId, fileName string, fileReader io.ReadCloser) (*http.Request, error) {
	bodyReader, bodyWriter := io.Pipe()
	writer := multipart.NewWriter(bodyWriter)

	go func() {
		defer bodyWriter.Close()
		err := writer.WriteField(folderIdAttribute, folderId)
		if err != nil {
			c.logger.Printf("failed to write 'folder ID' into multipart body: %v\n", err)
			bodyWriter.CloseWithError(err)
			return
		}
		part, err := writer.CreateFormFile(fileAttribute, fileName)
		if err != nil {
			c.logger.Printf("error creating form file for multipart writer: %v\n", err)
			bodyWriter.CloseWithError(err)
			return
		}
		_, err = io.Copy(part, fileReader)
		if err != nil {
			c.logger.Printf("failed to copy file into multipart body error: %v\n", err)
			bodyWriter.CloseWithError(err)
			return
		}
		err = fileReader.Close()
		if err = writer.Close(); err != nil {
			c.logger.Printf("closing resources error: %v\n", err)
			bodyWriter.CloseWithError(err)
			return
		}
	}()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, postFileEndpoint, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("creating post file request: %w", err)
	}
	req.Header.Set(contentTypeHeader, writer.FormDataContentType())

	c.logger.Printf("Created file upload request for file %s to folder id %s\n", fileName, folderId)

	return req, nil
}

const getFileEndpoint = "https://%s.gofile.io/download/web/%s/%s"

func (c *GofileClient) createGetFileRequest(ctx context.Context, server, fileId, fileName string) (*http.Request, error) {
	url := fmt.Sprintf(getFileEndpoint, server, url.PathEscape(fileId), url.PathEscape(fileName))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating 'getFile' request: %w", err)
	}
	return req, nil
}

func (c *GofileClient) createGetFileInfoRequest(ctx context.Context, fileId string) (*http.Request, error) {
	url := fmt.Sprintf("%s%s", contentsEndpointPart, url.PathEscape(fileId))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	req.Header.Set("X-Website-Token", "4fd6sg89d7s6")
	if err != nil {
		return nil, fmt.Errorf("creating 'getFile' request: %w", err)
	}
	return req, nil
}

func (c *GofileClient) createGetIdRequest(ctx context.Context) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, accountsEndpointPart+"getid", nil)
	if err != nil {
		return nil, fmt.Errorf("creating 'getid' request: %w", err)
	}

	return req, nil
}

func (c *GofileClient) createGetAccountInfoRequest(ctx context.Context, accountId string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, accountsEndpointPart+accountId, nil)
	if err != nil {
		return nil, fmt.Errorf("creating 'getAccountInfo' request: %w", err)
	}

	return req, nil
}

func (c *GofileClient) accountId(ctx context.Context) (string, error) {
	c.accountIdOnce.Do(func() {
		req, err := c.createGetIdRequest(ctx)
		if err != nil {
			c.accountIdError = fmt.Errorf("creating 'getid' request: %w", err)
			return
		}
		resp, err := c.do(req)
		if err != nil {
			c.accountIdError = fmt.Errorf("sending 'getid' request: %w", err)
			return
		}
		defer resp.Body.Close()

		var getIdResp getIdResponseData
		err = json.NewDecoder(resp.Body).Decode(&getIdResp)
		if err != nil {
			c.accountIdError = fmt.Errorf("unmarshalling 'getid' response: %w", err)
			return
		}

		c.accountIdCached = getIdResp.Data.Id
		if c.accountIdCached == "" {
			c.accountIdError = fmt.Errorf("empty accountId error")
			return
		}
	})

	return c.accountIdCached, c.accountIdError
}

func (c *GofileClient) rootFolderId(ctx context.Context) (string, error) {
	c.rootFolderIdOnce.Do(func() {
		accountId, err := c.accountId(ctx)
		if err != nil {
			c.rootFolderIdError = fmt.Errorf("failed to get 'accountID': %w", err)
			return
		}

		req, err := c.createGetAccountInfoRequest(ctx, accountId)
		if err != nil {
			c.rootFolderIdError = fmt.Errorf("failed to create 'getAccountInfo' request: %w", err)
			return
		}
		resp, err := c.do(req)
		if err != nil {
			c.rootFolderIdError = fmt.Errorf("failed to send 'getAccountInfo' request: %w", err)
			return
		}
		defer resp.Body.Close()

		var getAccountInfoResp getAccountInfoResponseData
		err = json.NewDecoder(resp.Body).Decode(&getAccountInfoResp)
		if err != nil {
			c.rootFolderIdError = fmt.Errorf("failed to unmarshal 'getAccountInfo' response: %w", err)
			return
		}

		c.rootFolderIdCached = getAccountInfoResp.Data.RootFolder
		if c.rootFolderIdCached == "" {
			c.rootFolderIdError = fmt.Errorf("empty root folder error")
			return
		}
	})

	return c.rootFolderIdCached, c.rootFolderIdError
}
