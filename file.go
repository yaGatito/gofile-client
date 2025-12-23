package gofile

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// GetFileInfo retrieves metadata information for the specified file.
func (c *GofileClient) GetFileInfo(ctx context.Context, websiteToken, fileId string) (GetFileInfoResponseBody, error) {
	req, err := c.createGetFileInfoRequest(ctx, websiteToken, fileId)
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

// DownloadFile downloads a file from the specified GoFile server.
//
// The caller is responsible for closing the returned ReadCloser.
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

// createGetFileRequest builds an HTTP GET request for getting a file
// with the specified server, fieldId, name.
func (c *GofileClient) createGetFileRequest(ctx context.Context, server, fileId, fileName string) (*http.Request, error) {
	url := fmt.Sprintf(getFileEndpoint, server, url.PathEscape(fileId), url.PathEscape(fileName))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating 'getFile' request: %w", err)
	}
	return req, nil
}
