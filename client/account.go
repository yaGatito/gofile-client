package gofile

import (
	"context"
	"encoding/json"
	"fmt"
)

// accountId resolves and caches the account ID associated with the API key.
//
// The value is fetched once and reused for subsequent calls.
// The method is safe for concurrent use.
func (c *GofileClient) accountId(ctx context.Context) (string, error) {
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

// rootFolderId resolves and caches the root folder ID of the account.
//
// The value is fetched once and reused for subsequent calls.
// The method is safe for concurrent use.
func (c *GofileClient) rootFolderId(ctx context.Context) (string, error) {
	c.rootFolderIdOnce.Do(func() {
		accountId, err := c.accountId(ctx)
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
