package server

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	"github.com/ziflex/lecho/v3"
	"kirill5k/reqfol/internal/config"
	"net/http"
)

type Server interface {
	Start() error
	Close() error
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
	if err := s.echo.Start(fmt.Sprintf(":%d", s.port)); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (s *echoServer) Close() error {
	return s.echo.Close()
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

func NewEchoServer(config *config.Server) Server {
	e := echo.New()
	logger := lecho.From(
		log.Logger,
		lecho.WithTimestamp(),
		lecho.WithCaller(),
	)
	e.Logger = logger
	e.Use(middleware.RequestID())
	e.Use(lecho.Middleware(lecho.Config{Logger: logger}))
	return &echoServer{config.Port, e, nil}
}
