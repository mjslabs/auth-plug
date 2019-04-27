package login

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mjslabs/auth-plug/auth"
)

// Function we want to use to authenticate the user
var authenticator = auth.ValidateLogin

// Function for generating the authentication token
var generator = auth.JWTCreateToken

// Post implements the HTTP POST method on /login.
// This takes a username and password and sends back
// a JWT on validation of the supplied credentials.
func Post(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	if username == "" || password == "" {
		return echo.ErrBadRequest
	}

	authed, user, err := authenticator(username, password)
	if !authed && err == nil {
		return echo.ErrUnauthorized
	} else if !authed && err != nil {
		return err
	}

	token, err := generator(user)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]string{
		"token": token,
	})
}
