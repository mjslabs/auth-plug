package verify

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Get implements the HTTP GET method on /verify.
// This just serves as an OK message assuming the JWT
// middleware let's the request go through to here
func Get(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}
