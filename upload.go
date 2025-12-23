package gofile

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

const (
	folderIdAttribute = "folderId"
	fileAttribute     = "file"
)

// UploadFile uploads a file to the specified folder.
//
// The folderId may be a concrete folder identifier or the special value "root".
// When "root" is provided, the client's root folder ID is resolved automatically.
//
// The provided fileReader is fully consumed and closed by this method.
func (c *GofileClient) UploadFile(
	ctx context.Context,
	folderId, fileName string,
	fileReader io.ReadCloser,
) (UploadFileResponseBody, error) {

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

	var result UploadFileResponseBody
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return UploadFileResponseBody{}, err
	}
	return result, nil
}

// createPostFileRequest constructs a streaming multipart/form-data
// HTTP request for file upload.
//
// The request body is produced asynchronously using an io.Pipe to avoid
// buffering the entire file in memory.
//
// The provided fileReader is consumed and closed during request body generation.
func (c *GofileClient) createPostFileRequest(
	ctx context.Context,
	folderId, fileName string,
	fileReader io.ReadCloser,
) (*http.Request, error) {

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
