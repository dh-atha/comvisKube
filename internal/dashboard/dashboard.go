package dashboard

import (
	"html/template"
	"io"
)

type Page struct {
	template *template.Template
}

func New() (*Page, error) {
	tmpl, err := template.New("dashboard").Parse(pageTemplate)
	if err != nil {
		return nil, err
	}

	return &Page{template: tmpl}, nil
}

func (p *Page) Render(w io.Writer, data any) error {
	return p.template.Execute(w, data)
}
