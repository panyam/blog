---
title: 'Build a static site generator (SSG) - Getting Started'
date: 2024-05-28T11:29:10AM
tags: ['static site generator', 'go']
draft: false
authors: ['Sri Panyam']
summary: Let us get started on implementing our static site genrator.
page: PostPage
location: "BodyView.ContentView"
---

## Background

This is the second part of our [static site generator series](../).  For our goals and requirements see [part 1](../requirements) We are building out a simple SSG library.   So let start with a "Site" type in (`cmd/sample.go`) to represent some info about our Site, where files are stored and more:

```go
type Site struct {
	// ContentRoot is the root of all your pages.
	// One structure we want to place is use folders to emphasis url structure too
	// so a site with mysite.com/a/b/c/d
	// would have the file structure:
	// <ContentRoot>/a/b/c/d.{md,html}
	ContentRoot string

	// Final output directory where resources are generated/published to
	OutputDir string

	// The http path prefix the site is prefixed in,
	// The site could be served from a subpath in the domain, eg:
	// eg
	//		myblog.com								=> PathPrefix = "/"
	//		myblog.com/blog						=> PathPrefix = "/blog"
	//		myblog.com/blogs/blog1		=> PathPrefix = "/blogs/blog1"
	//
	// There is no restriction on this.  There could be other routes (eg /products page)
	// that could be served by a different router all together in parallel to say /blog.
	// This is only needed so that the generator knows where to "root" the blog in the URL
	PathPrefix string

	// A list of folders where static files could be served from along with their
	// http path prefixes
	StaticFolders []string

	// A list of GLOBs that will point to several html templates our generator will parse and use
	HtmlTemplates []string

	// A list of GLOBs that will point to several text templates our generator will parse and use
	TextTemplates []string
}
```

With this we can now create a site:

```go
package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Site struct {
  // Definition of Site from above
}

var site = Site{
	ContentRoot: "../content",
	OutputDir:   "/Users/sri/personal/golang/blog/build",
	PathPrefix:  "/ourprefix",
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
  
	srv := &http.Server{
		Handler: web.session.LoadAndSave(web.router),
		Addr:    ":8888"
	}
	log.Fatal(srv.ListenAndServe())
}

```

The code above does a few things:

1. We create a "Site" instace and set some key variables - ContentRoot, PathPrefix, OutputDir etc.  The "PathPrefix" is interesting.  Instead of serving our site at the root (eg http://localhost), we want to serve it from the "/ourprefix" prefix - ie "http://localhost/ourprefix"
2. In our main function we create a router (using the gorilla mux router) and register the site to be the http handler at `PathPrefix`.
3. The site is not yet a valid [`http.HandlerFunc`](https://pkg.go.dev/net/http#HandlerFunc) implementation as it lacks the right methods for this interface.   We will now implement those.


Let us add the method to implement the [`http.HandlerFunc`](https://pkg.go.dev/net/http#HandlerFunc) interface.  It needs a single function:

```go
func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request)
```

and our implementation is simple:

```go
func (s *Site) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.GetRouter().ServeHTTP(w, r)
}

// Returns the underlying mux.Router serving static and generated files
func (s *Site) GetRouter() *mux.Router {
	if s.filesRouter == nil {
		s.filesRouter = mux.NewRouter()

		// Setup local/static paths
		for i := 0; i < len(s.StaticFolders); i += 2 {
			path, folder := s.StaticFolders[i], s.StaticFolders[i+1]
			log.Printf("Adding static route: %s -> %s", path, folder)
			s.filesRouter.PathPrefix(path).Handler(
				http.StripPrefix(path, http.FileServer(http.Dir(folder))))
		}

		// Serve everything else from the

		// Now add the file loader/handler for the "published" dir
		fileServer := http.FileServer(http.Dir(s.OutputDir))
		realHandler := http.StripPrefix("/", fileServer)
		x := s.filesRouter.PathPrefix("/")
		x.Handler(realHandler)
	}
	return s.filesRouter
}
```

Also add a `filesRouter` field to the Site type:

```go
type Site struct {
  ...

	// The router that serves the static resources and the built output files
	filesRouter *mux.Router
}
```

Explaining the above code:

1. Our ServeHTTP method simply calls the underlying GetRouter()'s ServeHTTP method.  The router itself can be used independently if ever needed.
2. Not mentioend earlier, the StaticFolders is a list of strings of the form:
    ```
    (http_prefix1, local_folder1, http_prefix2, local_folder2, ...)
    ```
3. The GetRouter method creates an internal filesRouter which serves each of the `http_prefix1..N` paths with the corresponding `local_folder1..N` folder.  Finally it serves the `OutputDir` also as a static folder from the site's `PathPrefix`
