package internal

import (
	"embed"
	"html/template"
	"io/fs"
)

type TemplateStore interface {
	GetIndex() (*template.Template, error)
	GetSlide() (*template.Template, error)
}

type InBuiltTemplateStore struct {
	fs fs.FS
}

//go:embed templates
var inBuiltTemplates embed.FS

func NewInBuiltTemplateStore() (InBuiltTemplateStore, error) {
	return InBuiltTemplateStore{fs: inBuiltTemplates}, nil
}

func (s InBuiltTemplateStore) GetIndex() (*template.Template, error) {
	tmpl, err := template.ParseFS(s.fs, "templates/index.tmpl")
	if err != nil {
		return nil, err
	}

	return tmpl, nil
}

func (s InBuiltTemplateStore) GetSlide() (*template.Template, error) {
	tmpl, err := template.ParseFS(s.fs, "templates/googleslides.tmpl")
	if err != nil {
		return nil, err
	}

	return tmpl, nil
}
