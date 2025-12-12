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
	"sync"
)

// TODO: logger support

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

	accountIdCached string
	accountIdOnce   sync.Once
	accountIdError  error

	rootFolderIdCached string
	rootFolderIdOnce   sync.Once
	rootFolderIdError  error
}

// NewClient creates client with provided API key. It will create a new client if the provided one is.
func NewClient(apiKey string, client *http.Client) (*GofileClient, error) {
	if client == nil {
		client = &http.Client{}
	}

	return &GofileClient{
		apiKey: apiKey,
		client: client,
	}, nil
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

func (c *GofileClient) UploadFile(ctx context.Context, folderId, fileName string, fileReader io.Reader) (UploadFileResponseBody, error) {
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

func (c *GofileClient) do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	log.Default().Printf("Sending request to %s", req.URL.String())

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("sending request: %w", err)
	}

	// Check error responses
	if resp.Header.Get(contentTypeHeader) == "text/html" {
		log.Default().Printf("Received HTML response body")
		return nil, fmt.Errorf("received HTML response, possible error page")
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("received bad status: %s", resp.Status)
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

func (c *GofileClient) createPostFileRequest(ctx context.Context, folderId, fileName string, fileReader io.Reader) (*http.Request, error) {
	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)

	go func() {
		defer pw.Close()
		defer writer.Close()
		go func() {
			<-ctx.Done()
			pw.CloseWithError(ctx.Err())
		}()

		err := writer.WriteField(folderIdAttribute, folderId)
		if err != nil {
			log.Default().Printf("failed to write 'folder ID' into multipart body: %v", err)
			pw.CloseWithError(err)
			return
		}
		part, err := writer.CreateFormFile(fileAttribute, fileName)
		if err != nil {
			log.Default().Printf("error creating form file ofr multipart writer: %v", err)
			pw.CloseWithError(err)
			return
		}
		_, err = io.Copy(part, fileReader)
		if err != nil {
			log.Default().Printf("failed to copy file %v into multipart body", err)
			pw.CloseWithError(err)
			return
		}
		if err := writer.Close(); err != nil {
			log.Default().Printf("error closing multipart writer: %v", err)
			pw.CloseWithError(err)
			return
		}
	}()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, postFileEndpoint, pr)
	if err != nil {
		return nil, fmt.Errorf("creating post file request: %w", err)
	}
	req.Header.Set(contentTypeHeader, writer.FormDataContentType())

	log.Default().Printf("Created file upload request for file %s to folder id %s", fileName, folderId)

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

func (c *GofileClient) accountID(ctx context.Context) (string, error) {
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
		accountId, err := c.accountID(ctx)
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
