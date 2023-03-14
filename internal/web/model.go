package web

import (
	"io"
	"text/template"
)

type TemplateData struct {
}

type Template struct {
	*template.Template
}

func NewTemplateModel(t *template.Template) *Template {
	return &Template{t}
}

func (t *Template) Execute(w io.Writer, data *TemplateData) error {
	return t.Template.Execute(w, data)
}
