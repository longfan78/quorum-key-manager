package client

import (
	"context"
	"fmt"

	"github.com/longfan78/quorum-key-manager/src/stores/api/types"
)

const keysPath = "keys"

func (c *HTTPClient) CreateKey(ctx context.Context, storeName, id string, req *types.CreateKeyRequest) (*types.KeyResponse, error) {
	key := &types.KeyResponse{}
	reqURL := fmt.Sprintf("%s/%s/%s", withURLStore(c.config.URL, storeName), keysPath, id)
	response, err := postRequest(ctx, c.client, reqURL, req)
	if err != nil {
		return nil, err
	}

	defer closeResponse(response)
	err = parseResponse(response, key)
	if err != nil {
		return nil, err
	}

	return key, nil
}

func (c *HTTPClient) ImportKey(ctx context.Context, storeName, id string, req *types.ImportKeyRequest) (*types.KeyResponse, error) {
	key := &types.KeyResponse{}
	reqURL := fmt.Sprintf("%s/%s/%s/import", withURLStore(c.config.URL, storeName), keysPath, id)
	response, err := postRequest(ctx, c.client, reqURL, req)
	if err != nil {
		return nil, err
	}

	defer closeResponse(response)
	err = parseResponse(response, key)
	if err != nil {
		return nil, err
	}

	return key, nil
}

func (c *HTTPClient) SignKey(ctx context.Context, storeName, id string, req *types.SignBase64PayloadRequest) (string, error) {
	reqURL := fmt.Sprintf("%s/%s/%s/sign", withURLStore(c.config.URL, storeName), keysPath, id)
	response, err := postRequest(ctx, c.client, reqURL, req)
	if err != nil {
		return "", err
	}

	defer closeResponse(response)
	return parseStringResponse(response)
}

func (c *HTTPClient) GetKey(ctx context.Context, storeName, id string) (*types.KeyResponse, error) {
	key := &types.KeyResponse{}
	reqURL := fmt.Sprintf("%s/%s/%s", withURLStore(c.config.URL, storeName), keysPath, id)

	response, err := getRequest(ctx, c.client, reqURL)
	if err != nil {
		return nil, err
	}

	defer closeResponse(response)
	err = parseResponse(response, key)
	if err != nil {
		return nil, err
	}

	return key, nil
}

func (c *HTTPClient) GetDeletedKey(ctx context.Context, storeName, id string) (*types.KeyResponse, error) {
	key := &types.KeyResponse{}
	reqURL := fmt.Sprintf("%s/%s/%s?deleted=true", withURLStore(c.config.URL, storeName), keysPath, id)

	response, err := getRequest(ctx, c.client, reqURL)
	if err != nil {
		return nil, err
	}

	defer closeResponse(response)
	err = parseResponse(response, key)
	if err != nil {
		return nil, err
	}

	return key, nil
}

func (c *HTTPClient) UpdateKey(ctx context.Context, storeName, id string, req *types.UpdateKeyRequest) (*types.KeyResponse, error) {
	key := &types.KeyResponse{}
	reqURL := fmt.Sprintf("%s/%s/%s", withURLStore(c.config.URL, storeName), keysPath, id)

	response, err := patchRequest(ctx, c.client, reqURL, req)
	if err != nil {
		return nil, err
	}

	defer closeResponse(response)
	err = parseResponse(response, key)
	if err != nil {
		return nil, err
	}

	return key, nil
}

func (c *HTTPClient) ListKeys(ctx context.Context, storeName string, limit, page uint64) ([]string, error) {
	return listRequest(ctx, c.client, fmt.Sprintf("%s/%s", withURLStore(c.config.URL, storeName), keysPath), false, limit, page)
}

func (c *HTTPClient) ListDeletedKeys(ctx context.Context, storeName string, limit, page uint64) ([]string, error) {
	return listRequest(ctx, c.client, fmt.Sprintf("%s/%s", withURLStore(c.config.URL, storeName), keysPath), true, limit, page)
}

func (c *HTTPClient) DeleteKey(ctx context.Context, storeName, id string) error {
	reqURL := fmt.Sprintf("%s/%s/%s", withURLStore(c.config.URL, storeName), keysPath, id)
	response, err := deleteRequest(ctx, c.client, reqURL)
	if err != nil {
		return err
	}

	defer closeResponse(response)
	return parseEmptyBodyResponse(response)
}

func (c *HTTPClient) DestroyKey(ctx context.Context, storeName, id string) error {
	reqURL := fmt.Sprintf("%s/%s/%s/destroy", withURLStore(c.config.URL, storeName), keysPath, id)
	response, err := deleteRequest(ctx, c.client, reqURL)
	if err != nil {
		return err
	}

	defer closeResponse(response)
	return parseEmptyBodyResponse(response)
}

func (c *HTTPClient) RestoreKey(ctx context.Context, storeName, id string) error {
	reqURL := fmt.Sprintf("%s/%s/%s/restore", withURLStore(c.config.URL, storeName), keysPath, id)
	response, err := putRequest(ctx, c.client, reqURL, nil)
	if err != nil {
		return err
	}

	defer closeResponse(response)
	return parseEmptyBodyResponse(response)
}
