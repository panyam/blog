package web

import (
	"html/template"
	"log"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/panyam/blog/web/components"
	s3core "github.com/panyam/s3gen/core"
)

var site = s3core.Site{
	ContentRoot: "./data",
	OutputDir:   "/Users/sri/personal/golang/blog/published",
	PathPrefix:  "",
	LazyLoad:    true,
	HtmlTemplates: []string{
		"templates/*.html",
	},
	StaticFolders: []string{
		"/static/", "public/static",
	},
}

var siteMetadata = &components.SiteMetadata{
	HeaderTitle: "Buildmage",
	Title:       "Buildmage",
	Description: "My personal blog highlighting adventures in building random things on the side",
	Author:      "Sriram Panyam",
}
var headerNavLinks = []components.HeaderNavLink{
	{Href: "/blog", Title: "Blog"},
	{Href: "/tags", Title: "Tags"},
	{Href: "/projects", Title: "Projects"},
	{Href: "/about", Title: "About"},
}

// This should be mirroring how we are setting up our app.yaml
func (web *BlogWeb) setupPages(router *mux.Router) {
	site.SiteMetadata = siteMetadata
	site.Init().Load().StartWatching()

	// Here we want to point just to the root of our blog and let it get served
	// For now we will serve via a router but then take the same router to
	// publish them for static serving too
	router.PathPrefix(site.PathPrefix).Handler(http.StripPrefix(site.PathPrefix, &site))
}

func (web *BlogWeb) RenderView(v components.View, w http.ResponseWriter, r *http.Request, templateName string) {
	ctx := &components.Context{
		Writer:       w,
		Template:     web.Template,
		SiteMetadata: siteMetadata,
	}
	// w.WriteHeader(http.StatusOK)
	v.InitContext(ctx, nil)
	err := v.ValidateRequest(w, r)
	if err == nil {
		err = ctx.Render(v, templateName)
	}
	if err != nil {
		slog.Error("Render Error: ", "err", err)
		http.Error(w, err.Error(), 500)
		// c.Abort()
	}
}

func CustomFuncMap() template.FuncMap {
	return template.FuncMap{
		"RenderView": func(v components.View) template.URL {
			log.Println("V: ", v)
			v.RenderResponse()
			return template.URL(" ")
		},
	}
}
