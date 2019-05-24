package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mjslabs/auth-plug/auth"
)

func healthGet(c echo.Context) error {
	status := http.StatusOK
	statusMsg := "OK"

	// Check LDAP connection
	if err := auth.Cfg.Serv.Conn.Connect(); err != nil {
		status = http.StatusServiceUnavailable
		statusMsg = err.Error()
	}

	return c.JSON(status, map[string]string{
		"status":  statusMsg,
		"version": version,
	})
}
