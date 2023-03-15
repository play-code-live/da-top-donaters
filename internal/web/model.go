package web

import (
	"io"
	"text/template"
)

type Template struct {
	*template.Template
}

func NewTemplateModel(t *template.Template) *Template {
	return &Template{t}
}

func (t *Template) Execute(w io.Writer, data any) error {
	return t.Template.Execute(w, data)
}
