package proxy

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"io"
	"kirill5k/reqfol/internal/server"
	"net/http"
)

const (
	HeaderXRerouteTo = "X-Reroute-To"
)

type Api struct {
	client Client
}

func NewApi(client Client) *Api {
	return &Api{client: client}
}

/*
TODO:
1: Sanitise headers
2: Logging
3: Interrapt on 403
*/
func (api *Api) RegisterRoutes(server server.Server) {
	handler := func(ctx echo.Context) error {

		redirectUrl := ctx.Request().Header.Get(HeaderXRerouteTo)
		if redirectUrl == "" {
			return ctx.String(http.StatusForbidden, fmt.Sprintf("Missing %s header", HeaderXRerouteTo))
		}
		headers := make(map[string]string)
		for hk, hv := range ctx.Request().Header {
			headers[hk] = hv[0]
		}
		queryParams := make(map[string]string)
		for pk, pv := range ctx.Request().URL.Query() {
			queryParams[pk] = pv[0]
		}
		var requestBody = ""
		if body, err := io.ReadAll(ctx.Request().Body); err == nil {
			requestBody = string(body[:])
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
		return ctx.Blob(res.StatusCode, res.ContentType, []byte(res.Body))
	}

	server.AddRoute("GET", "/*", handler)
	server.AddRoute("POST", "/*", handler)
}
