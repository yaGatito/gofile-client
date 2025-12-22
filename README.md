
# Unofficial GoFile API Client (Go)

[![Go Reference](https://pkg.go.dev/badge/github.com/yaGatito/gofile.svg)](https://pkg.go.dev/github.com/yaGatito/gofile-client)
[![Go Report Card](https://goreportcard.com/badge/github.com/yaGatito/gofile)](https://goreportcard.com/report/github.com/yaGatito/gofile-client)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

Unofficial Go client for interacting with the [gofile.io](https://gofile.io) file hosting service.

## Features

- Upload files (streaming multipart, no full buffering)
- Create folders
- Download files
- Retrieve file metadata
- Automatic caching of account and root folder IDs
- Concurrency-safe client

## Installation

```bash
go get github.com/yaGatito/gofile-client
```

## Usage

```go
package main

import (
    "context"
    "log"
    "os"

    "github.com/yaGatito/gofile"
)

func main() {
    ctx := context.Background()

    client := gofile.New("your-api-key", nil, nil)
    if client == nil {
        log.Fatal("Failed to create GoFile client, empty API key?")
    }

    file, err := os.Open("example.txt")
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    resp, err := client.UploadFile(ctx, "root", "example.txt", file)
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("File uploaded: %+v", resp)
}
```

## API Overview

### Interface

```go
type Gofile interface {
    GetFileInfo(ctx context.Context, fileId string) (GetFileInfoResponseBody, error)
    DownloadFile(ctx context.Context, server, fileId, fileName string) (io.ReadCloser, error)
    CreateFolder(ctx context.Context, parentFolderId, newFolderName string) (CreateFolderResponseBody, error)
    UploadFile(ctx context.Context, folderId, fileName string, fileReader io.ReadCloser) (UploadFileResponseBody, error)
}
```

### Concurrency

Client methods are safe for concurrent use.  
Account and root folder IDs are cached internally.

## Known Limitations

- GET endpoints temporarily unavailable for non-Premium users
- Requires `X-Website-Token` header for some downloads
- Some content moved to cold storage and requires Premium account import

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Disclaimer

This is an unofficial client based on observed API behavior and is not affiliated with gofile.io. Use at your own risk.
