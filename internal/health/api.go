package health

import (
	"github.com/labstack/echo/v4"
	"kirill5k/reqfol/internal/server"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

type Api struct {
	startupTime time.Time
	ipAddress   string
	appVersion  string
}

func NewApi() *Api {
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
		startupTime: time.Time{},
		ipAddress:   getIpaddress(),
		appVersion:  os.Getenv("VERSION"),
	}
}

func (api *Api) RegisterRoutes(server server.Server) {
	server.AddRoute("GET", "/health/status", func(ctx echo.Context) error {
		status := StatusUp(api.startupTime, api.ipAddress, api.appVersion)
		return ctx.JSON(http.StatusOK, status)
	})
}
