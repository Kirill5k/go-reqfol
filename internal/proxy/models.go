package proxy

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
)

type RequestMetadata struct {
	Method      string
	Url         string
	Headers     map[string]string
	QueryParams map[string]string
	Body        []byte
}

func (rm *RequestMetadata) String() string {
	url := rm.Url
	query := ""
	for k, v := range rm.QueryParams {
		if query != "" {
			query = query + "&"
		}
		query = query + k + "=" + v
	}
	if query != "" && strings.Contains(url, "?") {
		url = url + "&" + query
	} else if query != "" {
		url = url + "?" + query
	}

	headerNames := make([]string, 0, len(rm.Headers))
	for h := range rm.Headers {
		headerNames = append(headerNames, h)
	}
	sort.Strings(headerNames)
	headers := ""
	for _, k := range headerNames {
		if headers != "" {
			headers = headers + ", "
		}
		headers = headers + k + ":" + rm.Headers[k]
	}
	headers = "{" + headers + "}"

	return fmt.Sprintf("%s %s %s", rm.Method, url, headers)
}

type ResponseMetadata struct {
	StatusCode  int
	Body        []byte
	ContentType string
}

func (rm *ResponseMetadata) IsForbidden() bool {
	return rm.StatusCode == http.StatusForbidden
}
