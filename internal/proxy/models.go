package proxy

type RequestMetadata struct {
	Method      string
	url         string
	Headers     map[string]string
	QueryParams map[string]string
	body        string
}

type ResponseMetadata struct {
	StatusCode  int
	Body        string
	ContentType string
}
