package web

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/panyam/blog/web/components"
)

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

func DefaultFuncMap() template.FuncMap {
	return template.FuncMap{
		"Now": time.Now,

		"HTML": func(s string) template.HTML {
			return template.HTML(s)
		},

		"JS": func(s string) template.JS {
			return template.JS(s)
		},

		"URL": func(s string) template.URL {
			return template.URL(s)
		},

		"expandAttrs": func(attrs map[string]any) template.JS {
			out := " "
			if attrs != nil {
				for key, value := range attrs {
					val := fmt.Sprintf("%v", value)
					val = strings.Replace(val, "\"", "&quot;", -1)
					val = strings.Replace(val, "\"", "&quot;", -1)
					out += " " + fmt.Sprintf("%s = \"%s\"", key, val)
				}
			}
			return template.JS(out)
		},

		"dict": func(values ...interface{}) (map[string]interface{}, error) {
			if len(values)%2 != 0 {
				return nil, errors.New("invalid dict call")
			}
			dict := make(map[string]interface{}, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					return nil, errors.New("dict keys must be strings")
				}
				dict[key] = values[i+1]
			}
			return dict, nil
		},
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

// This should be mirroring how we are setting up our app.yaml
func (web *BlogWeb) setupPages(router *mux.Router) {
	var err error
	web.Template, err = template.New("hello").
		Funcs(DefaultFuncMap()).
		Funcs(CustomFuncMap()).
		ParseGlob("./web/components/*.html")
	if err != nil {
		panic(err)
	}

	router.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		aboutpage := &components.BasePage{
			HeaderNavLinks: headerNavLinks,
			BodyView:       &components.AboutPage{},
		}
		web.RenderView(aboutpage, w, r, "AboutPage")
	}).Methods("GET")

	// router.SetFuncMap(funcmap)
	// router.LoadHTMLGlob("./web/templates/*.*")
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		homepage := &components.BasePage{
			BaseView:       components.BaseView{},
			HeaderNavLinks: headerNavLinks,
			BodyView:       &components.HomePage{},
		}
		web.RenderView(homepage, w, r, "")
	}).Methods("GET")
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
