package proxy

type Api struct {
	client *Client
}

func NewApi() *Api {
	return &Api{}
}
