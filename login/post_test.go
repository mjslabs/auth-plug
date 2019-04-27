package login

import (
	"errors"
	"net/http"
	"net/url"
	"regexp"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/mjslabs/auth-plug/auth"
	"github.com/mjslabs/auth-plug/internal/testlib"
	"github.com/stretchr/testify/assert"
)

func TestPost(t *testing.T) {
	authenticator = func(username, password string) (bool, auth.User, error) {
		return false, auth.User{}, errors.New("stuff blew up")
	}
	generator = func(u auth.User) (string, error) {
		return "123.abc.beef", nil
	}
	t.Run("Post_error", testError)

	authenticator = func(username, password string) (bool, auth.User, error) {
		return false, auth.User{}, nil
	}
	t.Run("Post_invalid", testInvalid)
	t.Run("Post_failed", testFailed)

	authenticator = func(username, password string) (bool, auth.User, error) {
		return true, auth.User{}, nil
	}
	t.Run("Post_valid_token", testValidToken)

	generator = func(u auth.User) (string, error) {
		return "", errors.New("argh")
	}
	t.Run("Post_invalid_token", testInvalidToken)
}

func testError(t *testing.T) {
	// Error on login
	f := make(url.Values)
	f.Set("username", "gimmie")
	f.Set("password", "anerror")
	c, _ := testlib.EchoSetup(http.MethodPost, "/login", f.Encode())
	assert.Error(t, Post(c))
}

func testInvalidToken(t *testing.T) {
	// Error on token create
	auth.Cfg.JWTMethod = "invalid"
	f := make(url.Values)
	f.Set("username", "testy")
	f.Set("password", "tester")
	c, _ := testlib.EchoSetup(http.MethodPost, "/login", f.Encode())
	assert.Error(t, Post(c))
}

func testInvalid(t *testing.T) {
	// Invalid request
	c, _ := testlib.EchoSetup(http.MethodPost, "/login", "")
	assert.Equal(t, Post(c), echo.ErrBadRequest)
}

func testFailed(t *testing.T) {
	// Failed login
	f := make(url.Values)
	f.Set("username", "test")
	f.Set("password", "hello")
	c, _ := testlib.EchoSetup(http.MethodPost, "/login", f.Encode())
	assert.Equal(t, echo.ErrUnauthorized, Post(c))
}

func testValidToken(t *testing.T) {
	// Successful login and token creation
	f := make(url.Values)
	f.Set("username", "testy")
	f.Set("password", "tester")
	c, rec := testlib.EchoSetup(http.MethodPost, "/login", f.Encode())
	assert.NoError(t, Post(c))
	assert.Regexp(t, regexp.MustCompile(`^\s*{"token":"[A-Za-z0-9-_=]+\.[A-Za-z0-9-_=]+\.?[A-Za-z0-9-_.+/=]*"}\s*$`), rec.Body.String())
}
