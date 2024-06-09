package web

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/alexedwards/scs/v2"
	"github.com/felixge/httpsnoop"
	"github.com/gorilla/mux"
	gohttp "github.com/panyam/goutils/http"
	s3 "github.com/panyam/s3gen/core"
)

func DefaultGatewayAddress() string {
	gateway_addr := os.Getenv("BLOG_WEB_PORT")
	if gateway_addr != "" {
		return gateway_addr
	}
	return ":8080"
}

type BlogWeb struct {
	// TODO - turn this over to slicer to manage clients
	addr     string
	router   *mux.Router
	BaseUrl  string
	Template *template.Template
	session  *scs.SessionManager
}

func NewWebApp(addr string) (web *BlogWeb, err error) {
	web = &BlogWeb{
		addr:    addr,
		session: scs.New(), //scs.NewCookieManager("u46IpCV9y5Vlur8YvODJEhgOY8m9JVE4"),
	}
	return
}

func (web *BlogWeb) Start() {
	// Run the router
	web.setupRouter()
	// web.router.Run(web.addr)
	srv := &http.Server{
		Handler: withLogger(web.session.LoadAndSave(web.router)),
		Addr:    web.addr,
		// Good practice: enforce timeouts for servers you create!
		// WriteTimeout: 15 * time.Second,
		// ReadTimeout:  15 * time.Second,
	}
	log.Printf("Serving Gateway endpoint on %s:", web.addr)
	log.Fatal(srv.ListenAndServe())
}

func (web *BlogWeb) setupRouter() {
	web.router = mux.NewRouter()

	// Setup local dev by allowing CORS
	web.router.Use(gohttp.CORS)

	//setup basic pages
	web.setupSite(web.router)
}

func withLogger(handler http.Handler) http.Handler {
	// the create a handler
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		// pass the handler to httpsnoop to get http status and latency
		m := httpsnoop.CaptureMetrics(handler, writer, request)
		// printing exracted data
		log.Printf("http[%d]-- %s -- %s\n", m.Code, m.Duration, request.URL.Path)
	})
}

var site = s3.Site{
	ContentRoot: "./data",
	OutputDir:   "/Users/sri/personal/golang/blog/build",
	PathPrefix:  "/published",
	LazyLoad:    true,
	HtmlTemplates: []string{
		"templates/*.html",
	},
	StaticFolders: []string{
		"/static/", "static",
	},
}

// This should be mirroring how we are setting up our app.yaml
func (web *BlogWeb) setupSite(router *mux.Router) {
	site.NewViewFunc = web.NewView
	site.CommonFuncMap = web.TemplateFunctions()

	site.Init().Load().StartWatching()

	// Here we want to point just to the root of our blog and let it get served
	// For now we will serve via a router but then take the same router to
	// publish them for static serving too
	router.PathPrefix(site.PathPrefix).Handler(http.StripPrefix(site.PathPrefix, &site))
}
