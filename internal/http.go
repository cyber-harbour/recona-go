package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	BaseURL = "https://api.recona.io"

	authorizationHeaderName = "Authorization"
	authorizationType       = "Bearer "
	contentTypeHeaderName   = "Content-Type"
	acceptHeaderName        = "Accept"
	defaultContentType      = "application/json"

	DefaultRateLimit = 10
	DefaultBurst     = 2
)

func MakeAuthenticatedRequest(
	ctx context.Context, client *http.Client, method, url, token string, body interface{}) (*http.Response, error) {
	if client == nil {
		return nil, fmt.Errorf("request failed, http client is empty")
	}

	var reqBody io.Reader

	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set(authorizationHeaderName, authorizationType+token)
	req.Header.Set(contentTypeHeaderName, defaultContentType)
	req.Header.Set(acceptHeaderName, defaultContentType)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if resp.StatusCode >= 400 {
		var bodyBytes []byte
		if bodyBytes, err = io.ReadAll(resp.Body); err != nil {
			return nil, fmt.Errorf("API error %d: failed to read response body", resp.StatusCode)
		}
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return resp, nil
}

func DecodeJSON(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}
