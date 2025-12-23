
# Unofficial GoFile API Client (Go)

[![Go Reference](https://pkg.go.dev/badge/github.com/yaGatito/gofile.svg)](https://pkg.go.dev/github.com/yaGatito/gofile-client)
[![Go Report Card](https://goreportcard.com/badge/github.com/yaGatito/gofile)](https://goreportcard.com/report/github.com/yaGatito/gofile-client)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

Unofficial Go client for interacting with the [gofile.io](https://gofile.io) file hosting service.

## Features

- Upload files
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
func uploadUsecase() {
    ctx := context.Background()

	client, err := gofile.New("your-api-key", nil, nil)
	if err != nil {
		log.Fatal("Failed to create client:", err)
	}

	fileName := "sample.txt"
	file, err := os.Open("./" + fileName)
	if err != nil {
		log.Fatal("Failed to open file:", err)
	}

	postFileResponse, err := client.UploadFile(ctx, "your-folder-id", fileName, file)
	if err != nil {
		log.Fatal("Failed to upload file:", err)
	}
	log.Println("UploadFileResponseBody received:", postFileResponse)
}
```

## API Overview

### Interface

```go
type Gofile interface {
    GetFileInfo(ctx context.Context, websiteToken, fileId string) (GetFileInfoResponseBody, error)
    DownloadFile(ctx context.Context, server, fileId, fileName string) (io.ReadCloser, error)
    CreateFolder(ctx context.Context, parentFolderId, newFolderName string) (CreateFolderResponseBody, error)
    UploadFile(ctx context.Context, folderId, fileName string, fileReader io.ReadCloser) (UploadFileResponseBody, error)
}
```

## Known Limitations

- Check traffic and storage limitations: [gofile.io/myprofile](https://gofile.io/myprofile).
- Uploaded content may be moved to cold storage if inactive for a long time and requires importing into a Premium account to access.
- Requires `X-Website-Token` header to download a file until it moved to cold storage.
- GET endpoints unavailable for non-Premium users.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Disclaimer

This is an unofficial client based on observed API behavior and is not affiliated with gofile.io. Use at your own risk.
