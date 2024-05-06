package proxy

type Module struct {
	Api *Api
}

func NewModule() *Module {
	api := NewApi()
	return &Module{Api: api}
}
