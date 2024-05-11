package proxy

import (
	"github.com/labstack/echo/v4"
	"io"
	"kirill5k/reqfol/internal/server"
	"net/http"
	"strings"
)

type Api struct {
	client Client
}

func NewApi() *Api {
	return &Api{}
}

/*TODO:
1: Sanitise headers
2: Logging
3: Interrapt on 403
*/
func (api *Api) RegisterRoutes(server server.Server) {
	handler := func(ctx echo.Context) error {

		req := newRequestMetadata(ctx.Request())
		res, err := api.client.Send(*req)
		if err != nil {
			return ctx.String(http.StatusInternalServerError, err.Error())
		}
		return ctx.Blob(res.StatusCode, res.ContentType, []byte(res.Body))
	}

	server.AddRoute("GET", "/", handler)
	server.AddRoute("POST", "/", handler)
}

func newRequestMetadata(req *http.Request) *RequestMetadata {
	headers := make(map[string]string)
	for hk, hv := range req.Header {
		headers[hk] = strings.ToLower(hv[0])
	}
	queryParams := make(map[string]string)
	for pk, pv := range req.URL.Query() {
		queryParams[pk] = pv[0]
	}

	var requestBody = ""
	if body, err := io.ReadAll(req.Body); err == nil {
		requestBody = string(body[:])
	}

	var url = req.URL.Host + req.URL.Path
	if redirectUrl, ok := headers["x-reroute-to"]; ok {
		url = redirectUrl + req.URL.Path
	}

	return &RequestMetadata{
		Method:      req.Method,
		Url:         url,
		Headers:     headers,
		QueryParams: queryParams,
		Body:        requestBody,
	}
}
