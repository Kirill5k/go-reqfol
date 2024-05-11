package proxy

import "kirill5k/reqfol/internal/config"

type Module struct {
	Api *Api
}

func NewModule(conf *config.Client) *Module {
	client := NewRestyClient(conf)
	api := NewApi(client)
	return &Module{Api: api}
}
