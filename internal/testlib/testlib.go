package testlib

import (
	"errors"
	"net/http/httptest"
	"strings"

	"github.com/c0sco/go-ldap-client"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// EchoSetup gives an echo Context and resp recorder to work with in test functions
func EchoSetup(method, path, reader string) (c echo.Context, rec *httptest.ResponseRecorder) {
	e := echo.New()
	e.Use(middleware.Logger())
	req := httptest.NewRequest(method, path, strings.NewReader(reader))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	return c, rec
}

// MockLDAP lets us test LDAP code without an actual connection to a server
type MockLDAP ldap.LDAPClient

// MockLDAPFields is for field mapping tests in ValidateLogin
var MockLDAPFields = []string{"uid", "gecos"}

// LDAPHostFail - when the LDAP host is set to this, Connect() will fail
var LDAPHostFail = "fail.com"

// Authenticate mocks the method of the same name from LDAPClient
func (l MockLDAP) Authenticate(username string, password string) (bool, map[string]string, error) {
	if username == "testy" && password == "tester" {
		return true, map[string]string{}, nil
	} else if username == "gimmie" && password == "anerror" {
		return false, map[string]string{}, errors.New("the server can't be found")
	}
	// go-ldap-client returns an error with this string on failed login
	return false, map[string]string{}, errors.New("Invalid Credentials")
}

// Close mocks the method of the same name from LDAPClient
func (l MockLDAP) Close() {
}

// Connect mocks the method of the same name from LDAPClient
func (l MockLDAP) Connect() error {
	if l.Host == LDAPHostFail {
		return errors.New("Connection to server failed")
	}
	return nil
}

// GetGroupsOfUser mocks the method of the same name from LDAPClient
func (l MockLDAP) GetGroupsOfUser(username string) ([]string, error) {
	return []string{}, nil
}
