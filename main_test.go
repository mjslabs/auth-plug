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
var testPortSuccess = "13337"

// Port to start the web server on for testing failures
var testPortFail = "65536"

func TestMain(t *testing.T) {
	t.Run("SetupRoutes", testSetupRoutes)
	// The next two functions capture main()'s stdout, implementing
	// a timeout in case something goes wrong
	t.Run("Main", testMainSuccess)
	t.Run("Main_fail", testMainFail)
}

func testMainFail(t *testing.T) {
	timer := time.After(time.Second * 10)
	old, r, w := newStdout()
	outC := make(chan string)
	go myIOCopy(r, outC)

	os.Setenv("AUTH_PORT", testPortFail)
	go main()
	time.Sleep(time.Second * 3)

	select {
	// The channel main() listens for signals on
	case quit <- syscall.SIGINT:
	case <-timer:
		t.Log("timed out starting HTTP server")
		t.Fail()
	}
	resetStdout(w, old)
	out := <-outC
	assert.Contains(t, out, "error starting server: listen tcp: address "+testPortFail)
}

func testMainSuccess(t *testing.T) {
	timer := time.After(time.Second * 10)
	old, r, w := newStdout()
	outC := make(chan string)
	go myIOCopy(r, outC)

	os.Setenv("AUTH_PORT", testPortSuccess)
	go main()
	time.Sleep(time.Second * 3)

	select {
	case quit <- syscall.SIGINT:
	case <-timer:
		t.Log("timed out starting HTTP server")
		t.Fail()
	}
	resetStdout(w, old)
	out := <-outC
	assert.Contains(t, out, "http server started on [::]:"+testPortSuccess)
}

func testSetupRoutes(t *testing.T) {
	e := echo.New()
	setupRoutes(e, routes)
	r := e.Routes()

	for _, route := range routes {
		for name, method := range route.methods {
			for _, i := range r {
				if i.Method == name && i.Path == route.path {
					// Validates that the entry in routes got properly turned into an echo route handler
					assert.Equal(t, i.Name, runtime.FuncForPC(reflect.ValueOf(method).Pointer()).Name())
				}
			}
		}
	}
}

func myIOCopy(r *os.File, outC chan string) {
	var buf bytes.Buffer
	io.Copy(&buf, r)
	outC <- buf.String()
}

func newStdout() (old, newRead, newWrite *os.File) {
	old = os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	return old, r, w
}

func resetStdout(w, old *os.File) {
	w.Close()
	os.Stdout = old
}
