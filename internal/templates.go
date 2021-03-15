package internal

import (
	"bytes"
	"embed"
	"html/template"
	"io/fs"
	"net/http"
)

type IndexData struct {
	Title  string
	Author string
	Links  []Link
	Groups []TalkGroup
}

type TemplateStore interface {
	GetIndex() (*template.Template, error)
	GenerateIndex(IndexData) ([]byte, error)
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

func (s InBuiltTemplateStore) GenerateIndex(data IndexData) ([]byte, error) {
	var buf bytes.Buffer

	template, err := s.GetIndex()
	if err != nil {
		return []byte{}, err
	}

	if err := template.Execute(&buf, data); err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil
}

func (s InBuiltTemplateStore) WriteIndex(w http.ResponseWriter, data IndexData) error {
	return nil
}

func (s InBuiltTemplateStore) GetSlide() (*template.Template, error) {
	tmpl, err := template.ParseFS(s.fs, "templates/googleslides.tmpl")
	if err != nil {
		return nil, err
	}

	return tmpl, nil
}
