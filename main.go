package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/felixge/httpsnoop"
	"github.com/gorilla/mux"
	s3 "github.com/panyam/s3gen"
)

var (
	addr = flag.String("addr", DefaultGatewayAddress(), "Address where the http grpc gateway endpoint is running")
)

var site = s3.Site{
	OutputDir:   "./output",
	ContentRoot: "./content",
	TemplateFolders: []string{
		"./templates",
	},
	StaticFolders: []string{
		"/static/", "static",
		"/demos/", "demos",
	},
	DefaultBaseTemplate: s3.BaseTemplate{
		Name:   "BasePage.html",
		Params: map[any]any{"BodyTemplateName": "BaseBody"},
	},
	/*
		GetTemplate: func(res *s3.Resource, out *s3.PageTemplate) {
			relpath := res.RelPath()
			if strings.HasPrefix(relpath, "/blog/") {
				out.Params = map[any]any{"BodyTemplateName": "PostSimple"}
			}
		},
	*/
}

func main() {
	flag.Parse()
	if os.Getenv("APP_ENV") != "production" {
		site.Watch()
	}

	// Attach our site to be at /`PathPrefix`
	// The site will also take care of serving static files from /`PathPrefix`/static paths
	router := mux.NewRouter()
	router.PathPrefix(site.PathPrefix).Handler(http.StripPrefix(site.PathPrefix, &site))

	srv := &http.Server{
		Handler: withLogger(router),
		Addr:    *addr,
		// Good practice: enforce timeouts for servers you create!
		// WriteTimeout: 15 * time.Second,
		// ReadTimeout:  15 * time.Second,
	}
	log.Printf("Serving Gateway endpoint on %s:", *addr)
	log.Fatal(srv.ListenAndServe())
}

func DefaultGatewayAddress() string {
	gateway_addr := os.Getenv("BLOG_WEB_PORT")
	if gateway_addr != "" {
		return gateway_addr
	}
	return ":8080"
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
