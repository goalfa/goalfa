package exporter

type Maker interface {
	Lang() string
	Make(pkg string, methods []*Method) (files []*File, err error)
}
