package exporter

type AxiosMaker struct {
}

func (a AxiosMaker) Lang() string {
	return Ts
}

func (a AxiosMaker) Make(pkg string, methods []*Method) (files []*File, err error) {
	return
}
