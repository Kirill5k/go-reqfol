package health

import (
	"kirill5k/reqfol/internal/interrupter"
)

type Module struct {
	Api *Api
}

func NewModule(interrupter interrupter.Interrupter) *Module {
	api := NewApi(interrupter)
	return &Module{Api: api}
}
