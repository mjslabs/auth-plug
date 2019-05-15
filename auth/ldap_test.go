package auth

import (
	"reflect"
	"testing"

	"github.com/mjslabs/auth-plug/internal/testlib"
	"github.com/stretchr/testify/assert"
)

// init is a stand-in for Initialize()
func init() {
	Cfg = Configuration{
		JWTMethod: "HS512",
		Serv: ServerAttributes{
			Conn: testlib.MockLDAP{},
		},
	}
}

func TestLDAP(t *testing.T) {
	t.Run("Validate_errored", testErrored)
	t.Run("Validate_failed", testFailed)
	t.Run("Validate_success", testSuccess)
	t.Run("populateFields", testPopulateFields)
}

func testErrored(t *testing.T) {
	// Error on login
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
	valid, _, err := ValidateLogin("testy", "tester")
	assert.Equal(t, err, nil)
	assert.Equal(t, valid, true)
}

func testPopulateFields(t *testing.T) {
	myTagName := "mapstructure"
	myFields := populateFields(myTagName, User{})

	// Independently verify that populateFields parsed the struct correctly
	r := reflect.TypeOf(User{})
	for i := 0; i < r.NumField(); i++ {
		tag := r.Field(i).Tag.Get(myTagName)
		if tag == "" {
			continue
		}
		if inSlice := stringInSlice(tag, myFields); !inSlice {
			t.Errorf(
				"LDAP data not mapped into User struct properly: %s tag '%s' not in %s",
				myTagName, tag, myFields)
		}
	}
}

func stringInSlice(s string, list []string) bool {
	for _, item := range list {
		if item == s {
			return true
		}
	}
	return false
}
