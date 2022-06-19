package exporter

import (
	"encoding/json"
	"fmt"
)

type Folder struct {
	Name      string  `json:"name"`
	Namespace string  `json:"namespace"`
	Files     []*File `json:"files"`
}

type File struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

func NewSDK(methods []*Method) *SDK {
	return &SDK{methods: methods}
}

type SDK struct {
	methods []*Method
}

func (p SDK) Make(makers map[string]Maker, lang, pkg string) ([]byte, error) {
	maker, ok := makers[lang]
	if !ok {
		return nil, fmt.Errorf("target '%s' maker not found", lang)
	}
	var methods []*Method
	for _, v := range p.methods {
		methods = append(methods, v.Fork())
	}
	files, err := maker.Make(pkg,methods)
	if err != nil {
		return nil, err
	}
	return json.Marshal(files)
}
