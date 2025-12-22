package gofile

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

// do sends an HTTP request using the underlying http.Client.
//
// The method automatically attaches the Authorization header,
// logs the outgoing request, and validates the HTTP response.
//
// It returns an error if:
//   - the request fails at the transport level
//   - the response status code is >= 400
//   - the response content type indicates an HTML error page
//
// On success, the caller is responsible for closing the response body.
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
