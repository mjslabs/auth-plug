package verify

import (
	"net/http"
	"testing"

	"github.com/mjslabs/auth-plug/internal/testlib"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	c, rec := testlib.EchoSetup(http.MethodGet, "/verify", "")
	if assert.NoError(t, Get(c)) {
		assert.Equal(t, "OK", rec.Body.String())
	}
}
