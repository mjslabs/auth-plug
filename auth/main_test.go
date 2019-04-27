package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	t.Run("Initialize", testInitialize)
}

func testInitialize(t *testing.T) {
	Initialize()
	assert.NotEmpty(t, Cfg.JWTMethod)
	assert.NotEmpty(t, Cfg.JWTValidMinutes)
	// Test the server struct while we're here
	assert.NotEmpty(t, Cfg.Serv.Host)
	assert.NotEmpty(t, Cfg.Serv.Port)
	assert.NotEmpty(t, Cfg.Serv.UIDFieldName)
	assert.NotEmpty(t, Cfg.Serv.GIDFieldName)
	assert.NotEmpty(t, Cfg.Serv.Fields)
}
