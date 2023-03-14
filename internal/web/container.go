package web

import (
	"errors"
	"os"
	"strings"
	"text/template"
)

type TemplateContainer struct {
	templates map[string]*Template
	masterTpl *template.Template
}

func NewTemplateContainer(masterPath string) *TemplateContainer {
	return &TemplateContainer{
		masterTpl: template.Must(template.ParseGlob(masterPath)),
		templates: make(map[string]*Template),
	}
}

func (t *TemplateContainer) GetAll() map[string]*Template {
	return t.templates
}

func (t *TemplateContainer) Get(name string) (*Template, error) {
	tpl, exists := t.templates[name]
	if !exists {
		return nil, errors.New("requested templated does not exist")
	}

	return tpl, nil
}

func (t *TemplateContainer) MustGet(name string) *Template {
	tpl, err := t.Get(name)
	if err != nil {
		panic(err)
	}

	return tpl
}

func (t *TemplateContainer) Register(name, path string) error {
	tpl, err := template.Must(t.masterTpl.Clone()).ParseFiles(path)
	if err != nil {
		return err
	}

	t.templates[name] = NewTemplateModel(tpl)

	return nil
}

func (t *TemplateContainer) FindAndRegister(path string) error {
	dirEntries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, e := range dirEntries {
		if e.IsDir() {
			continue
		}

		if !strings.Contains(e.Name(), ".gohtml") {
			continue
		}

		name := strings.TrimSuffix(e.Name(), ".gohtml")
		if err = t.Register(name, path+e.Name()); err != nil {
			return err
		}
	}

	return nil
}
