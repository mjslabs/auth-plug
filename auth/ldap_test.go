package auth

import (
	"errors"
	"reflect"
	"testing"

	"github.com/jtblin/go-ldap-client"
	"github.com/oleiade/reflections"
	"github.com/stretchr/testify/assert"
)

// MockLDAP lets us test LDAP code without an actual connection to a server
type MockLDAP ldap.LDAPClient

var testUserInfo = "testing"

// Authenticate mocks the method of the same name from LDAPClient
func (l MockLDAP) Authenticate(username string, password string) (bool, map[string]string, error) {
	if username == "testy" && password == "tester" {
		Cfg.Serv.Fields = populateFields("mapstructure", User{})
		userData := map[string]string{}
		for _, field := range Cfg.Serv.Fields {
			userData[field] = field + testUserInfo
		}
		return true, userData, nil
	} else if username == "gimmie" && password == "anerror" {
		return false, map[string]string{}, errors.New("the server can't be found")
	}
	return false, map[string]string{}, nil
}

// Close mocks the method of the same name from LDAPClient
func (l MockLDAP) Close() {
}

// Connect mocks the method of the same name from LDAPClient
func (l MockLDAP) Connect() error {
	return nil
}

// GetGroupsOfUser mocks the method of the same name from LDAPClient
func (l MockLDAP) GetGroupsOfUser(username string) ([]string, error) {
	return []string{}, nil
}

// init is a stand-in for Initialize()
func init() {
	Cfg = Configuration{
		JWTMethod: "HS512",
		Serv: ServerAttributes{
			Conn: MockLDAP{},
		},
	}
}

func TestLDAP(t *testing.T) {
	t.Run("Validate_errored", testErrored)
	t.Run("Validate_failed", testFailed)
	t.Run("Validate_success", testSuccess)
	t.Run("Connect", testConnect)
	t.Run("Initialize_success", testInitialize)
}

func testErrored(t *testing.T) {
	// Failed login
	valid, _, err := ValidateLogin("gimmie", "anerror")
	assert.Equal(t, valid, false)
	assert.NotNil(t, err)
}

func testFailed(t *testing.T) {
	// Failed login
	valid, _, err := ValidateLogin("test", "hello")
	assert.Equal(t, err, nil)
	assert.Equal(t, valid, false)
}

func testSuccess(t *testing.T) {
	// Correct login
	valid, user, _ := ValidateLogin("testy", "tester")
	assert.Equal(t, valid, true)
	r := reflect.TypeOf(user)
	for i := 0; i < r.NumField(); i++ {
		tag := r.Field(i).Tag.Get("mapstructure")
		if data, _ := reflections.GetField(user, r.Field(i).Name); tag != "" && data != tag+testUserInfo {
			t.Errorf("LDAP data not mapped into User structure properly: %s, %s", data, tag)
		}
	}
}

func testConnect(t *testing.T) {
	assert.NotPanics(t, func() { Cfg.Serv.Conn.Connect() })
}
