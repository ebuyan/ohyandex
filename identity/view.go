package identity

import (
	"net/http"
	"text/template"
)

type View struct{}
type Errors struct{ Error string }

func (v View) Render(w http.ResponseWriter, view string, err error) {
	errors := Errors{}
	if err != nil {
		errors.Error = err.Error()
	}

	parsedTemplate, _ := template.ParseFiles("static/" + view + ".html")
	err = parsedTemplate.Execute(w, errors)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
