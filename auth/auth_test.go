package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthenticate(t *testing.T) {
	token := "Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ0a2VlbCIsImV4cCI6MTY0NjE0MDI5Miwic3ViIjoidXNyLTIzNDgwMmM5YWQwY2NjOGUxYTViYWQ0NWZiNmMifQ.n3xo5lavvWz5tBV-Gs0UPFafP69Aumfn2L38DTm_E_VVUhLG7SblTBZqgtlyjHfD5qVmJH8iIsJmy-hkAWYz4w"
	auth, err := Authenticate(token, AuthTokenURLTestRemote)
	assert.NoError(t, err)
	assert.NotEmpty(t, auth)
}
