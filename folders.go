package gofile

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	rootFolderIdPlaceholderConst = "root"
	contentTypeHeader            = "Content-Type"
	applicationJsonContentType   = "application/json"
)

// CreateFolder creates a new folder under the specified parent folder.
//
// The parentFolderId may be a concrete folder identifier or the special value "root".
// When "root" is provided, the client's root folder ID is resolved automatically.
func (c *GofileClient) CreateFolder(
	ctx context.Context,
	parentFolderId,
	newFolderName string,
) (CreateFolderResponseBody, error) {

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

// createPostFolderRequest builds an HTTP POST request for creating a folder
// under the specified parent folder.
func (c *GofileClient) createPostFolderRequest(
	ctx context.Context,
	parentFolderId,
	folderName string,
) (*http.Request, error) {

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
	req.Header.Set(contentTypeHeader, applicationJsonContentType)

	return req, nil
}
