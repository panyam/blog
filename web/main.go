package web

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/alexedwards/scs/v2"
	"github.com/felixge/httpsnoop"
	"github.com/gorilla/mux"
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

	web.setupLocalDev(web.router)

	// Enable the endpoint for API
	web.setupStaticRoutes()

	//setup basic pages
	web.setupPages(web.router)
}

func (web *BlogWeb) setupStaticRoutes() {
	web.router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))
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