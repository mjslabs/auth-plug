package auth

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

var (
	mockUser = User{
		Email:    "test@test.com",
		Realname: "Testy Tester",
		Username: "testy",
		Groups:   []string{},
	}
)

var defaultJWTSecret = "testing12345678"
var defaultJWTMethod = "HS512"

func TestJWT(t *testing.T) {
	t.Run("Create", testCreate)
	t.Run("SigningMethods", testSigningMethods)
}

func testCreate(t *testing.T) {
	Cfg.JWTSecret = defaultJWTSecret
	token, err := JWTCreateToken(mockUser)
	if assert.NoError(t, err) {
		assert.Regexp(t, regexp.MustCompile(`^[A-Za-z0-9-_=]+\.[A-Za-z0-9-_=]+\.?[A-Za-z0-9-_.+/=]*$`), token)
	}

	tokenObj, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok || token.Header["alg"] != defaultJWTMethod {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(defaultJWTSecret), nil
	})

	if _, ok := tokenObj.Claims.(jwt.MapClaims); !ok || !tokenObj.Valid {
		t.Errorf("failed with error: %s", err)
	}
}

func testSigningMethods(t *testing.T) {
	assert.Equal(t, func() jwt.SigningMethod { s, _ := jwtSigningMethodFromString("HS256"); return s }(), jwt.SigningMethodHS256)
	assert.Equal(t, func() jwt.SigningMethod { s, _ := jwtSigningMethodFromString("HS384"); return s }(), jwt.SigningMethodHS384)
	assert.Equal(t, func() jwt.SigningMethod { s, _ := jwtSigningMethodFromString("HS512"); return s }(), jwt.SigningMethodHS512)
}
