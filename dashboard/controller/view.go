package controller

import (
	"fmt"
	"html/template"
	"net/http"
)

func IndexPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("view/index.html"))
	if err := tmpl.Execute(w, nil); err != nil {
		fmt.Println("[IndexPageHandler]", err)
	}
}
