package testlib

import (
	"net/http/httptest"
	"strings"

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
