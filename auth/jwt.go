package auth

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/fatih/structs"
)

// JWTCreateToken creates a JWT with the claims configured
// in the User struct and sends it back to the client.
func JWTCreateToken(u User) (string, error) {
	method, err := jwtSigningMethodFromString(Cfg.JWTMethod)
	if method == nil || err != nil {
		return "", err
	}

	token := jwt.New(method)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(Cfg.JWTValidMinutes).Unix()
	claims["nbf"] = time.Now().Add(time.Second * -5).Unix()
	// Grab the fields from the User struct and put them into the jwt claims
	myMap := structs.Map(u)
	for claim := range myMap {
		claims[claim] = myMap[claim]
	}

	return token.SignedString([]byte(Cfg.JWTSecret))
}

func jwtSigningMethodFromString(method string) (jwt.SigningMethod, error) {
	switch method {
	case "HS256":
		return jwt.SigningMethodHS256, nil
	case "HS384":
		return jwt.SigningMethodHS384, nil
	case "HS512":
		return jwt.SigningMethodHS512, nil
	}

	return nil, errors.New("unsupported JWT signing method")
}
