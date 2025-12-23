package gofile

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

const websiteTokenHeader = "X-Website-Token"

// createGetFileInfoRequest builds an HTTP GET request for retrieving
// detailed metadata for the specified file ID.
func (c *GofileClient) createGetFileInfoRequest(ctx context.Context, wsToken, fileId string) (*http.Request, error) {
	url := fmt.Sprintf("%s%s", contentsBaseURL, url.PathEscape(fileId))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	req.Header.Set(websiteTokenHeader, wsToken)
	if err != nil {
		return nil, fmt.Errorf("creating 'getFile' request: %w", err)
	}
	return req, nil
}

// createGetIdRequest builds an HTTP GET request for retrieving
// the account ID associated with the API key in use.
func (c *GofileClient) createGetIdRequest(ctx context.Context) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, accountsBaseURL+"getid", nil)
	if err != nil {
		return nil, fmt.Errorf("creating 'getid' request: %w", err)
	}

	return req, nil
}

// createGetAccountInfoRequest builds an HTTP GET request for retrieving
// account metadata for the specified account ID.
func (c *GofileClient) createGetAccountInfoRequest(ctx context.Context, accountId string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, accountsBaseURL+accountId, nil)
	if err != nil {
		return nil, fmt.Errorf("creating 'getAccountInfo' request: %w", err)
	}

	return req, nil
}
