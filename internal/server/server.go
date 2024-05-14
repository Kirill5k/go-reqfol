package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	"github.com/ziflex/lecho/v3"
	"kirill5k/reqfol/internal/config"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Server interface {
	StartAndWaitForShutdown()
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

func (s *echoServer) StartAndWaitForShutdown() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Info().Msgf("Starting Echo server on port: %d", s.port)
		if err := s.echo.Start(fmt.Sprintf(":%d", s.port)); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Err(err).Msg("Error starting Echo server")
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)
	<-sigChan

	log.Info().Msg("Shutting down Echo server")
	if err := s.echo.Shutdown(context.Background()); err != nil {
		log.Err(err).Msg("Failed to shut down Echo server")
	}
	wg.Wait()
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
