package web

import (
	"html/template"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/panyam/blog/web/components"
)

// This should be mirroring how we are setting up our app.yaml
func (web *BlogWeb) setupPages(router *mux.Router) {
	funcmap := template.FuncMap{
		"renderView": func(v components.View) template.URL {
			v.RenderResponse()
			return template.URL(" ")
		},
	}
	var err error
	web.Template, err = template.New("hello").Funcs(funcmap).ParseGlob("./web/templates/*.*")
	if err != nil {
		panic(err)
	}

	// router.SetFuncMap(funcmap)
	// router.LoadHTMLGlob("./web/templates/*.*")
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		web.RenderView(&components.HomePage{}, w, r)
	}).Methods("GET")
}

func (web *BlogWeb) RenderView(v components.View, w http.ResponseWriter, r *http.Request) {
	ctx := &components.Context{
		Writer:   w,
		Template: web.Template,
	}
	// w.WriteHeader(http.StatusOK)
	v.InitContext(ctx, nil)
	err := v.ValidateRequest(w, r)
	if err == nil {
		err = v.RenderResponse()
	}
	if err != nil {
		slog.Error("Render Error: ", "err", err)
		http.Error(w, err.Error(), 500)
		// c.Abort()
	}
}
