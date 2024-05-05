package health

type Module struct {
	Api *Api
}

func NewModule() *Module {
	api := NewApi()
	return &Module{api}
}
