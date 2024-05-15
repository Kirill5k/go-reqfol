package proxy

import (
	"kirill5k/reqfol/internal/config"
	"kirill5k/reqfol/internal/interrupter"
)

type Module struct {
	Api *Api
}

func NewModule(conf *config.Client, inter interrupter.Interrupter) *Module {
	client := NewRestyClient(conf)
	api := NewApi(client, inter)
	return &Module{Api: api}
}
