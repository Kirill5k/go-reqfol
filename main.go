package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"kirill5k/reqfol/internal/config"
	"kirill5k/reqfol/internal/health"
	"kirill5k/reqfol/internal/interrupter"
	"kirill5k/reqfol/internal/proxy"
	"kirill5k/reqfol/internal/server"
	"os"
)

func main() {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, NoColor: false, TimeFormat: "2020-01-01T00:00:00.999"})

	log.Info().Msgf("Starting request-follower %s", os.Getenv("VERSION"))

	conf := config.LoadAppConfig()
	srv := server.NewEchoServer(&conf.Server)
	inter := interrupter.NewSignallingInterrupter(&conf.Interrupter)
	apis := []server.RouteRegister{
		health.NewModule(inter).Api,
		proxy.NewModule(&conf.Client, inter).Api,
	}
	for _, api := range apis {
		api.RegisterRoutes(srv)
	}

	srv.StartAndWaitForShutdown()
}
