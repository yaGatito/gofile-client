package gofile

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
)

// Gofile defines the public contract for interacting with the GoFile API.
//
// Implementations must be safe for concurrent use.
type Gofile interface {
	GetFileInfo(ctx context.Context, fileId string) (GetFileInfoResponseBody, error)
	DownloadFile(ctx context.Context, server, fileId, fileName string) (io.ReadCloser, error)
	CreateFolder(ctx context.Context, parentFolderId, newFolderName string) (CreateFolderResponseBody, error)
	UploadFile(ctx context.Context, folderId, fileName string, fileReader io.ReadCloser) (UploadFileResponseBody, error)
}

var _ Gofile = &GofileClient{}

// GofileClient is a reusable, concurrency-safe client for interacting
// with the GoFile API.
//
// A client instance caches account and root folder identifiers internally
// and may be used concurrently by multiple goroutines.
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

// New creates a new GofileClient using the provided API key.
//
// If httpClient is nil, http.DefaultClient is used.
// If logger is nil, a default logger writing to stdout is created.
//
// The function returns nil if apiKey is empty.
func New(apiKey string, client *http.Client, logger *log.Logger) (Gofile, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("empty apiKey")
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

	return c, nil
}
