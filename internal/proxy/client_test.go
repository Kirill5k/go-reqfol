package proxy

import (
	"github.com/stretchr/testify/require"
	"io"
	"kirill5k/reqfol/internal/config"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const successResponse = `{"message": "Hello, World!"}`
const errorResponse = `{"message": "Internal Server Error"}`

func TestRestyClient_Send_ReturnsOk(t *testing.T) {
	requestPath := ""
	requestBody := make([]byte, 0)
	requestHeaders := make(map[string]string)
	requestQueryParams := make(map[string]string)
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			requestPath = r.URL.Path
			for pn, pv := range r.URL.Query() {
				requestQueryParams[pn] = pv[0]
			}
			if body, err := io.ReadAll(r.Body); err == nil {
				requestBody = body[:]
			}
			for hk, hv := range r.Header {
				requestHeaders[hk] = hv[0]
			}
			w.Header().Set("Content-Type", "application/json")
			_, err := w.Write([]byte(successResponse))
			require.NoError(t, err)
		}))
	defer server.Close()

	client := NewRestyClient(clientConfig())

	request := RequestMetadata{
		Method:      "POST",
		Url:         server.URL + "/hello/world",
		Headers:     map[string]string{"Foo": "bar", "User-Agent": "test"},
		QueryParams: map[string]string{"param1": "value"},
		Body:        []byte(`{"body": "requestBody"}`),
	}
	response, err := client.Send(request)

	require.NoError(t, err)
	require.Contains(t, request.Url, requestPath)
	require.Equal(t, request.Body, requestBody)
	require.Subset(t, requestHeaders, request.Headers)
	require.Equal(t, request.QueryParams, requestQueryParams)
	require.Equal(t, "application/json", response.ContentType)
	require.Equal(t, 200, response.StatusCode)
	require.Equal(t, []byte(successResponse), response.Body)
}

func TestRestyClient_Send_RetriesOnerror(t *testing.T) {
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			requests++
			w.Header().Set("Content-Type", "application/json")
			if requests < 2 {
				w.WriteHeader(500)
				_, err := w.Write([]byte(errorResponse))
				require.NoError(t, err)
			} else {
				w.WriteHeader(200)
				_, err := w.Write([]byte(successResponse))
				require.NoError(t, err)
			}
		}))
	defer server.Close()

	client := NewRestyClient(clientConfig())

	request := RequestMetadata{
		Method:      "POST",
		Url:         server.URL + "/hello/world",
		Headers:     map[string]string{"Foo": "bar", "User-Agent": "test"},
		QueryParams: map[string]string{"param1": "value"},
		Body:        []byte(`{"body": "requestBody"}`),
	}
	response, err := client.Send(request)

	require.NoError(t, err)
	require.Equal(t, "application/json", response.ContentType)
	require.Equal(t, 200, response.StatusCode)
	require.Equal(t, []byte(successResponse), response.Body)
}

func clientConfig() *config.Client {
	return &config.Client{
		MaxIdleConns:        1,
		MaxIdleConnsPerHost: 1,
		IdleConnTimeout:     1,
		Timeout:             1 * time.Minute,
		RetryCount:          5,
		RetryWaitTime:       100 * time.Millisecond,
		RetryMaxWaitTime:    100 * time.Millisecond,
	}
}
