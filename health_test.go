package main

import (
	"net/http"
	"testing"

	"github.com/mjslabs/auth-plug/auth"
	"github.com/mjslabs/auth-plug/internal/testlib"
	"github.com/stretchr/testify/assert"
)

func TestMainHealth(t *testing.T) {
	t.Run("Success", testMainHealthSuccess)
	t.Run("Failure", testMainHealthFailure)
}

func testMainHealthSuccess(t *testing.T) {
	setAuthCfgWithLDAPHost("success.com")
	c, rec := testlib.EchoSetup(http.MethodPost, "/health", "")
	err := healthGet(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func testMainHealthFailure(t *testing.T) {
	setAuthCfgWithLDAPHost(testlib.LDAPHostFail)
	c, rec := testlib.EchoSetup(http.MethodPost, "/health", "")
	err := healthGet(c)
	assert.Equal(t, http.StatusServiceUnavailable, rec.Code)
	assert.NoError(t, err)
}

func setAuthCfgWithLDAPHost(host string) {
	auth.Cfg = auth.Configuration{
		Serv: auth.ServerAttributes{
			Conn: testlib.MockLDAP{
				Host: host,
			},
		},
	}
}
