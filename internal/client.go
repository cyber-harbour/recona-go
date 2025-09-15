package internal

import (
	"context"
	"net/http"
)

type Client interface {
	MakeRequest(ctx context.Context, method, endpoint string, body interface{}) (*http.Response, error)
}
