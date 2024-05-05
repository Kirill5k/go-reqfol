package main

import (
	"kirill5k/reqfol/internal/health"
	"kirill5k/reqfol/internal/server"
	"log"
)

func main() {
	conf := server.Config{Port: 8080}
	srv := server.NewEchoServer(&conf)

	apis := []server.RouteRegister{
		health.NewModule().Api,
	}
	for _, api := range apis {
		api.RegisterRoutes(srv)
	}

	if err := srv.Start(); err != nil {
		log.Fatalf("failed to start http server: %s", err)
	}
}
