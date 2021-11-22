package http

import (
	"context"
	"net/http"
)

var contextHTTPHeaderKey = struct{}{}

func HeaderFromContext(ctx context.Context) http.Header {
	h := ctx.Value(contextHTTPHeaderKey)
	header, ok := h.(http.Header)
	if !ok {
		return nil
	}
	return header
}

func ContextWithHeader(ctx context.Context, h http.Header) context.Context {
	return context.WithValue(ctx, contextHTTPHeaderKey, h)
}
