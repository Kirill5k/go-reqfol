package proxy

import "net/http"

type RequestMetadata struct {
	Method      string
	Url         string
	Headers     map[string]string
	QueryParams map[string]string
	Body        []byte
}

type ResponseMetadata struct {
	StatusCode  int
	Body        []byte
	ContentType string
}

func (rm *ResponseMetadata) IsForbidden() bool {
	return rm.StatusCode == http.StatusForbidden
}
