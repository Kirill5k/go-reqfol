package server

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/ziflex/lecho/v3"
	"net/http"
	"os"
)

type Config struct {
	Port int
}

type Server interface {
	Start() error
	PrefixRoute(prefix string)
	AddRoute(method, path string, handler echo.HandlerFunc)
}

type RouteRegister interface {
	RegisterRoutes(server Server)
}

type echoServer struct {
	port       int
	echo       *echo.Echo
	routeGroup *echo.Group
}

func (s *echoServer) Start() error {
	if err := s.echo.Start(fmt.Sprintf(":%d", s.port)); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *echoServer) AddRoute(method, path string, handler echo.HandlerFunc) {
	if s.routeGroup != nil {
		s.routeGroup.Add(method, path, handler)
	} else {
		s.echo.Add(method, path, handler)
	}
}

func (s *echoServer) PrefixRoute(prefix string) {
	s.routeGroup = s.echo.Group(prefix)
}

func NewEchoServer(config *Config) Server {
	e := echo.New()
	logger := lecho.New(
		os.Stdout,
		lecho.WithLevel(log.DEBUG),
		lecho.WithTimestamp(),
		lecho.WithCaller(),
	)
	e.Logger = logger
	e.Use(middleware.RequestID())
	e.Use(lecho.Middleware(lecho.Config{Logger: logger}))
	return &echoServer{config.Port, e, nil}
}
