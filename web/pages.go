package web

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

// This should be mirroring how we are setting up our app.yaml
func (web *BlogWeb) setupPages(router *mux.Router) {
	funcmap := template.FuncMap{}
	var err error
	web.Template, err = template.New("hello").Funcs(funcmap).ParseGlob("./web/templates/*.*")
	if err != nil {
		panic(err)
	}

	// router.SetFuncMap(funcmap)
	// router.LoadHTMLGlob("./web/templates/*.*")
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World")
	}).Methods("GET")
}
