package controller

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
)

func IndexPageHandler(w http.ResponseWriter, r *http.Request) {
	cwd, _ := os.Getwd()
	path := filepath.Join(cwd, "view/index.html")
	tmpl := template.Must(template.ParseFiles(path))
	if err := tmpl.Execute(w, nil); err != nil {
		fmt.Println("[IndexPageHandler]", err)
	}
}
