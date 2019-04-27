package main

import (
	"bytes"
	"io"
	"os"
	"reflect"
	"runtime"
	"syscall"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// Port to start the web server on for testing
var testPort = "13337"

func TestMain(t *testing.T) {
	t.Run("SetupRoutes", testSetupRoutes)
	t.Run("Main", testMain)
}

func testMain(t *testing.T) {
	timer := time.After(time.Second * 10)
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	outC := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	os.Setenv("AUTH_PORT", testPort)
	go main()
	time.Sleep(time.Second * 3)
	select {
	case quit <- syscall.SIGINT:
	case <-timer:
		t.Log("timed out starting HTTP server")
		t.Fail()
	}

	// back to normal state
	w.Close()
	os.Stdout = old
	out := <-outC
	assert.Contains(t, out, "http server started on [::]:"+testPort)
}

func testSetupRoutes(t *testing.T) {
	e := echo.New()
	setupRoutes(e, routes)
	r := e.Routes()

	for _, route := range routes {
		for name, method := range route.methods {
			for _, i := range r {
				if i.Method == name && i.Path == route.path {
					assert.Equal(t, i.Name, runtime.FuncForPC(reflect.ValueOf(method).Pointer()).Name())
				}
			}
		}
	}
}
