package proxy

import (
	"fmt"
	"io"
	"kirill5k/reqfol/internal/interrupter"
	"kirill5k/reqfol/internal/server"
	"net/http"
	"regexp"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

const (
	headerXRerouteTo  = "X-Reroute-To"
	headerReloadOn403 = "X-Reload-On-403"

	invalidHeaderRegex = "(?i)^((x-reroute|x-reload)-.*|host)$"
)

type Api struct {
	client      Client
	interrupter interrupter.Interrupter
}

func NewApi(client Client, inter interrupter.Interrupter) *Api {
	return &Api{client: client, interrupter: inter}
}

func (api *Api) RegisterRoutes(server server.Server) {
	handler := func(ctx echo.Context) error {

		redirectUrl := ctx.Request().Header.Get(headerXRerouteTo)
		if redirectUrl == "" {
			log.Warn().Msgf("Missing %s header for REQUEST %s %s", headerXRerouteTo, ctx.Request().Method, ctx.Request().URL.Path)
			return ctx.String(http.StatusForbidden, fmt.Sprintf("Missing %s header", headerXRerouteTo))
		}
		headers := make(map[string]string)
		for hk, hv := range ctx.Request().Header {
			if matches, _ := regexp.MatchString(invalidHeaderRegex, strings.ToLower(hk)); !matches {
				headers[hk] = hv[0]
			}
		}
		queryParams := make(map[string]string)
		for pk, pv := range ctx.Request().URL.Query() {
			queryParams[pk] = pv[0]
		}
		requestBody := make([]byte, 0)
		if body, err := io.ReadAll(ctx.Request().Body); err == nil {
			requestBody = body[:]
		}

		req := RequestMetadata{
			Method:      ctx.Request().Method,
			Url:         redirectUrl + ctx.Request().URL.Path,
			Headers:     headers,
			QueryParams: queryParams,
			Body:        requestBody,
		}
		res, err := api.client.Send(req)
		if err != nil {
			return ctx.String(http.StatusInternalServerError, err.Error())
		}
		if header := ctx.Request().Header.Get(headerReloadOn403); header != "" && res.IsForbidden() {
			api.interrupter.Interrupt()
		}
		log.Info().Msgf("REQUEST %s RESPONSE %d", req.String(), res.StatusCode)
		return ctx.Blob(res.StatusCode, res.ContentType, res.Body)
	}

	server.AddRoute("GET", "/*", handler)
	server.AddRoute("POST", "/*", handler)
}
