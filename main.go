package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/mjslabs/auth-plug/auth"
	"github.com/mjslabs/auth-plug/login"
	"github.com/mjslabs/auth-plug/verify"
	"github.com/mmcloughlin/professor"
)

var version = "undefined"

type methodMap map[string]func(c echo.Context) error

// RouteDef defines the structure of a route to register
type RouteDef struct {
	// path in URL
	path string
	// map of HTTP method to handler func
	methods methodMap
	// protected behind the JWT middleware
	restricted bool
}

// Routes and their handlers
var routes = []RouteDef{
	{
		path: "/login",
		methods: methodMap{
			"POST": login.Post,
		},
	},
	{
		path: "/verify",
		methods: methodMap{
			"GET": verify.Get,
		},
		restricted: true,
	},
	{
		path: "/health",
		methods: methodMap{
			"GET": healthGet,
		},
	},
}

// When sigint is sent to this, the server will attempt a clean shutdown
var quit = make(chan os.Signal)

func main() {
	// Echo instance
	e := echo.New()

	// pprof launcher
	if os.Getenv("AUTH_PROFILE") != "" {
		professor.Launch(os.Getenv("AUTH_PROFILE"))
	}

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Logger.SetLevel(log.INFO)

	// Initialize auth subsystem which gives access to auth.Cfg
	auth.Initialize()

	// Register routes
	setupRoutes(e, routes)

	// Start server
	go func() {
		s := &http.Server{
			Addr:         os.Getenv("AUTH_IP") + ":" + os.Getenv("AUTH_PORT"),
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
			IdleTimeout:  5 * time.Second,
		}

		if err := e.StartServer(s); err != nil {
			e.Logger.Errorf("error starting server: %s", err)
		}
	}()

	// Graceful server shutdown
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	auth.Cfg.Serv.Conn.Close()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

// Set up route definitions based on the info above
func setupRoutes(e *echo.Echo, routes []RouteDef) {
	for _, routeDef := range routes {
		r := e.Group(routeDef.path)
		if routeDef.restricted {
			r.Use(middleware.JWTWithConfig(middleware.JWTConfig{
				SigningKey:    []byte(auth.Cfg.JWTSecret),
				SigningMethod: auth.Cfg.JWTMethod,
			}))
		}
		for methodName, methodFunc := range routeDef.methods {
			switch methodName {
			case "GET":
				r.GET("", methodFunc)
			case "POST":
				r.POST("", methodFunc)
			}
		}
	}
}
