// Package gofile provides a minimal client for the GoFile API.
//
// The package exposes a single public client type that allows:
//   - uploading files
//   - creating folders
//   - downloading files
//   - retrieving file metadata
//
// All HTTP, caching, and request-building logic is internal to the package
// and is not part of the public API.
//
// Usage example:
//
//	client := gofile.NewClient(apiKey, nil, nil)
//	resp, err := client.UploadFile(ctx, "root", "file.txt", reader)
package gofile
