package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthenticate(t *testing.T) {
	token := "Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ0a2VlbCIsImV4cCI6MTY0NjEyMzQ2NSwic3ViIjoidXNyLTY1OTg3MGMzZTY5OTNlODgxMWMxOGRhNmM2YWEifQ.A3vgIXU2exn66uhDb3ANrtVBiz9yrn-y18f9HisaIhAUy208FkgFUmHtTdbZ3rgEDanRuaoBtwNGBkaA3ZOFXA"
	auth, err := Authenticate(token)
	assert.NoError(t, err)
	assert.NotEmpty(t, auth)
}
