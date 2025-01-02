package views

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/enlistedmango/lenslocked/models"
)

func Must(t Template, err error) Template {
	if err != nil {
		panic(err)
	}
	return t
}

func Parse(filepaths ...string) (Template, error) {
	tpl, err := template.ParseFiles(filepaths...)
	if err != nil {
		return Template{}, fmt.Errorf("parsing template: %w", err)
	}
	return Template{
		htmlTpl: tpl,
	}, nil
}

type Template struct {
	htmlTpl *template.Template
}

func (t Template) Execute(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Printf("Debug: Executing template with data: %+v\n", data)
	err := t.htmlTpl.Execute(w, data)
	if err != nil {
		log.Printf("executing template: %v", err)
		http.Error(w, "There was an error executing the template.", http.StatusInternalServerError)
		return
	}
}

type TemplateData struct {
	Nav       NavigationData
	Form      *FormData
	Alert     *Alert
	Title     string
	Gallery   *models.Gallery
	Galleries []models.Gallery
}
