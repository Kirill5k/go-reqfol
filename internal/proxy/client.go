package proxy

type Client interface {
	Send(req RequestMetadata) *ResponseMetadata
}
