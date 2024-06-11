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
	getSortedKeys := func(pairs map[string]string) []string {
		res := make([]string, 0, len(pairs))
		for v := range pairs {
			res = append(res, v)
		}
		sort.Strings(res)
		return res
	}

	url := rm.Url
	query := ""
	queryParams := getSortedKeys(rm.QueryParams)
	for _, k := range queryParams {
		if query != "" {
			query = query + "&"
		}
		query = query + k + "=" + rm.QueryParams[k]
	}
	if query != "" && strings.Contains(url, "?") {
		url = url + "&" + query
	} else if query != "" {
		url = url + "?" + query
	}

	headerNames := getSortedKeys(rm.Headers)
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
