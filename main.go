package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"kirill5k/reqfol/internal/config"
	"kirill5k/reqfol/internal/health"
	"kirill5k/reqfol/internal/proxy"
	"kirill5k/reqfol/internal/server"
	"os"
)

func main() {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, NoColor: false, TimeFormat: "2006-01-02T15:04:05.999"})

	log.Info().Msg("Starting request-follower")

	conf := config.LoadAppConfig()
	srv := server.NewEchoServer(&conf.Server)

	apis := []server.RouteRegister{
		health.NewModule().Api,
		proxy.NewModule(&conf.Client).Api,
	}
	for _, api := range apis {
		api.RegisterRoutes(srv)
	}

	if err := srv.Start(); err != nil {
		log.Fatal().Err(err).Msgf("Failed to start server on port %d", conf.Server.Port)
	}
}
