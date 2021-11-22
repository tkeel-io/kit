package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
)

func TestRegister(t *testing.T) {
	errUnknown := New(int(codes.Unknown), "unknown", "unknown")
	Register(errUnknown)

	// new error
	newErrUnknown := errUnknown.WithMessage("unknown something").WithMetadata(map[string]string{
		"hello": "world",
	}).(*TError)
	assert.Equal(t, "unknown something", newErrUnknown.GetMessage())
	assert.Equal(t, map[string]string{
		"hello": "world",
	}, newErrUnknown.GetMetadata())

	assert.ErrorIs(t, newErrUnknown, errUnknown)
}
