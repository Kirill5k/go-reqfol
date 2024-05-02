package main

import (
	"kirill5/reqfol/internal/server"
	"log"
)

func main() {
	conf := server.Cofing{}
	srv := server.NewEchoServer(&conf)

	if err := srv.Start(); err != nil {
		log.Fatalf("failed to start http server: %s", err)
	}
}
