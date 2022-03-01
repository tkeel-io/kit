package auth

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthenticate(t *testing.T) {
	token := "Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ0a2VlbCIsImV4cCI6MTY0NjEyMTk0Nywic3ViIjoidXNyLTIzNDgwMmM5YWQwY2NjOGUxYTViYWQ0NWZiNmMifQ.EUvVpq_ITnTZgJhH1KlUIzPReCU-IUnFN5FritWn2Co3GRhvTVDboR2xM7J4T2EIjI-5eq3dZrGOsSudx86_sA"
	auth, err := Authenticate(token)
	assert.NoError(t, err)
	fmt.Printf("%+v\n", auth)
	assert.NotEmpty(t, auth)
}
