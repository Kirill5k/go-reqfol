package proxy

import (
	"github.com/go-resty/resty/v2"
	"kirill5k/reqfol/internal/config"
	"net/http"
	"strings"
)

type Client interface {
	Send(req RequestMetadata) *ResponseMetadata
}

type restyClient struct {
	client *resty.Client
}

func NewRestyClient(conf config.Client) Client {
	client := resty.New().
		SetTransport(&http.Transport{
			MaxIdleConns:        conf.MaxIdleConns,
			MaxIdleConnsPerHost: conf.MaxIdleConnsPerHost,
			IdleConnTimeout:     conf.IdleConnTimeout,
		}).
		SetTimeout(conf.Timeout).
		SetRetryCount(conf.RetryCount).
		SetRetryWaitTime(conf.RetryWaitTime).
		SetRetryMaxWaitTime(conf.RetryMaxWaitTime).
		AddRetryCondition(func(response *resty.Response, err error) bool {
			return response.StatusCode() == http.StatusInternalServerError || response.StatusCode() == http.StatusRequestTimeout ||
				(err != nil && strings.Contains(err.Error(), "Client.Timeout"))
		})
	return restyClient{client: client}
}

func (rc restyClient) Send(req RequestMetadata) *ResponseMetadata {
	return &ResponseMetadata{}
}
