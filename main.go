package main

import (
	"kirill5k/reqfol/internal/config"
	"kirill5k/reqfol/internal/health"
	"kirill5k/reqfol/internal/proxy"
	"kirill5k/reqfol/internal/server"
	"log"
)

func main() {
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
		log.Fatalf("failed to start http server: %s", err)
	}
}
