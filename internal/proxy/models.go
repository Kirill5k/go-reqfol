package proxy

type RequestMetadata struct {
	Method      string
	Url         string
	Headers     map[string]string
	QueryParams map[string]string
	Body        string
}

type ResponseMetadata struct {
	StatusCode  int
	Body        string
	ContentType string
}
