package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

var site = Site{
	ContentRoot: "./content",
	OutputDir:   "/Users/sri/personal/golang/blog/build",
	PathPrefix:  "/published",
	HtmlTemplates: []string{
		"templates/*.html",
	},
	StaticFolders: []string{
		"/static/", "static",
	},
}

func main() {
	router := mux.NewRouter()
	router.PathPrefix(site.PathPrefix).Handler(http.StripPrefix(site.PathPrefix, &site))
}
