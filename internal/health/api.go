package health

import (
	"github.com/labstack/echo/v4"
	"kirill5k/reqfol/internal/interrupter"
	"kirill5k/reqfol/internal/server"
	"log"
	"net"
	"net/http"
	"os"
)

type Api struct {
	interrupter interrupter.Interrupter
	ipAddress   string
	appVersion  string
}

func NewApi(interrupter interrupter.Interrupter) *Api {
	getIpaddress := func() string {
		conn, err := net.Dial("udp", "8.8.8.8:80")
		if err != nil {
			log.Fatal(err)
		}
		defer func(conn net.Conn) {
			err := conn.Close()
			if err != nil {
				log.Fatal(err)
			}
		}(conn)

		localAddr := conn.LocalAddr().(*net.UDPAddr)

		return localAddr.String()
	}

	return &Api{
		interrupter: interrupter,
		ipAddress:   getIpaddress(),
		appVersion:  os.Getenv("VERSION"),
	}
}

func (api *Api) RegisterRoutes(server server.Server) {
	server.AddRoute("GET", "/health/status", func(ctx echo.Context) error {
		status := StatusUp(api.interrupter.StartupTime(), api.ipAddress, api.appVersion)
		return ctx.JSON(http.StatusOK, status)
	})

	server.AddRoute("DELETE", "/health/status", func(ctx echo.Context) error {
		isInterrupted := api.interrupter.Interrupt()
		if isInterrupted {
			status := StatusDown(api.interrupter.StartupTime(), api.ipAddress, api.appVersion)
			return ctx.JSON(http.StatusServiceUnavailable, status)
		} else {
			status := StatusUp(api.interrupter.StartupTime(), api.ipAddress, api.appVersion)
			return ctx.JSON(http.StatusForbidden, status)
		}
	})
}
