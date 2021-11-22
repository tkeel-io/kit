package http

import (
	"context"
	"net/http"
)

const ContextHTTPHeaderKey = "http_header"

func GetHeader(ctx context.Context) http.Header {
	h := ctx.Value(ContextHTTPHeaderKey)
	header, ok := h.(http.Header)
	if !ok {
		return nil
	}
	return header
}
