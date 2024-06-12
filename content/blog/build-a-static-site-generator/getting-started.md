---
title: 'Build a static site generator - Getting Started'
date: 2024-05-28T11:29:10AM
tags: ['static site generator', 'go', 'ssg']
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

func (s *Site) Init() {
  // Ensure our content and output directories have absolute paths
	s.ContentRoot = gut.ExpandUserPath(s.ContentRoot)
	s.OutputDir = gut.ExpandUserPath(s.OutputDir)
  s.Load()
}

func (s *Site) Load() *Site {
  // TODO - Builds the entire site
	log.Println("TODO - Building site...")
  return s
}
```

With this we can now create a site:

```go
package main

import (
	"net/http"

	"github.com/gorilla/mux"
	gut "github.com/panyam/goutils/utils" // just a simple utility library
)

type Site struct {
  // Definition of Site from above
}

var site = Site{
	ContentRoot: "../content",
	OutputDir:   "../output",
	PathPrefix:  "/ssgdemo",
	HtmlTemplates: []string{
		"templates/*.html",
	},
	StaticFolders: []string{
		"/static/", "static",
	},
}

func main() {
	site.Init() // Initialize it
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

1. We create a "Site" instance and set some key variables - ContentRoot, PathPrefix, OutputDir etc.  The "PathPrefix" is interesting.  Instead of serving our site at the root (eg `http://<hostname>`), we want to serve it from the "/ssgdemo" prefix - ie `http://<hostname>/ssgdemo` (our hostname here would be `localhost:8888`).
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
3. The GetRouter method creates an internal filesRouter which serves each of the `http_prefix1..N` paths with the corresponding `local_folder1..N` folder.  Finally it serves the `OutputDir` also as a static folder from the site's `PathPrefix`.  This way we could serve multiple `static` url-paths from different folders (eg `/static`, `/media`, `/docs` etc).

Start this with `cd cmd ; air` (we are using [air](https://github.com/air-verse/air) to enable live-reloads) and opening [our site](http://localhost:8888/ssgdemo/) should greet us with the iconic `Hello World`.  

## Walking the walk

We already have a lot of content in the "content" folder.  Now it is time to walk through it and generate html pages for them.  Let us do just that.  It is here one starts to *really* appreciate the Go standard library.  Let us look at the amazing [path/filepath](https://pkg.go.dev/path/filepath) module.  It has a lot of little helper methods for dealing with files and folders.  One of them is the [WalkDir](https://pkg.go.dev/path/filepath#WalkDir) method that allows us to recursively go through a folder and all its children.   The documentation at pkg.go.dev is quite detailed so do familiarize yourself with how WalkDir works.   We could collect all our content starting at the "ContentRoot" folder with:

```go
func (s *Site) WalkFiles(offset int, count int) (foundFiles []string) {
	// keep a map of files encountered and their statuses
	filepath.WalkDir(s.ContentRoot, func(fullpath string, info os.DirEntry, err error) error {
		if err != nil {
			// just print err related to the path and stop scanning
			// if this err means something else we can do other things here
			log.Println("Error in path: ", info, err)
			return err
		}

		if info.IsDir() {
      // this is a directory so return a nil error to indicate
      // this directory should be recursed into
			return nil
		}

    out = append(out, fullpath)
    
		return nil    // no errors found
	})
	if offset > 0 {
		foundFiles = foundFiles[offset:]
	}
	if count > 0 {
		foundFiles = foundFiles[:count]
	}
	return foundFiles
}
```

This users the filepath.WalkDir method and gets all files recursived starting at the site's `ContentRoot` folder.   It also allows "slicing" within this list.


But we can do slightly better as an interface.   We will be dealing *a lot* with files but it may be helpful to have an abstraction over them.   Let us go with the `Resource` (may be an Asset or Content could be a better name?.  Naming *is* hard).   Let us use the following definition for a Resource (which in our case is a wrapper over files but more as we will see):

```go

const (
	ResourceStatePending = iota
	ResourceStateLoaded
	ResourceStateDeleted
	ResourceStateNotFound
	ResourceStateFailed
)

type Resource struct {
	Site     *Site  // Site this resource belongs to
	FullPath string // Unique URI/Path

	// Info about the resource
	info os.FileInfo

	// Created timestamp on disk
	CreatedAt time.Time

	// Updated time stamp on disk
	UpdatedAt time.Time

	// Loaded, Pending, NotFound, Failed
	State int

	// Any errors with this resource
	Error error
}
```

Here instead of directly dealing with files (via their paths) we will create some metadata over it so we can get some helpers (like extensions, any errors in loading, their load state, some timestamps).   We can also add more to this as we evolve our requirements (which we shall see soon).  With this we can update our Walker method:

```go
type ResourceFilterFunc func(res *Resource) bool
type ResourceSortFunc func(a *Resource, b *Resource) bool

func (s *Site) ListResources(filterFunc ResourceFilterFunc,
	sortFunc ResourceSortFunc,
	offset int, count int) []*Resource {
	var foundResources []*Resource
	// keep a map of files encountered and their statuses
	filepath.WalkDir(s.ContentRoot, func(fullpath string, info os.DirEntry, err error) error {
		if err != nil {
			// just print err related to the path and stop scanning
			// if this err means something else we can do other things here
			log.Println("Error in path: ", info, err)
			return err
		}

		if info.IsDir() {
      // just recurse into directories as before
      // we should probably have a checker if this is a valid dir to recurse into (eg is it 
			return nil  
		}

		res := s.GetResource(fullpath)

		if filterFunc == nil || filterFunc(res) {
			foundResources = append(foundResources, res)
		}

		return nil
	})
	if sortFunc != nil {
		sort.Slice(foundResources, func(idx1, idx2 int) bool {
			ent1 := foundResources[idx1]
			ent2 := foundResources[idx2]
			return sortFunc(ent1, ent2)
		})
	}
	if offset > 0 {
		foundResources = foundResources[offset:]
	}
	if count > 0 {
		foundResources = foundResources[:count]
	}
	return foundResources
}
```

We have a done a few things now:

1. Instead of working directly with files, we have introduced a Resource construct that keeps track of a file in our content root.  
2. Our List method also accepts functions to "filter", "sort" and "slice" resources by *some* criteria (this will be useful in the future).
3. We convert the fullpath of a visited file with the `GetResource` method on the site.  This gives us an opportunity to cache resource entries along with all its state etc.

Our GetResource implementation is simple, we keep encountered resources in a private map:

```go
type Site struct {
  ...
  
	resources map[string]*Resource
}

func (s *Site) GetResource(fullpath string) *Resource {
	res, found := s.resources[fullpath]
	if res == nil || !found {
		res = &Resource{
			Site:      s,
			FullPath:  fullpath,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			State:     ResourceStatePending,
		}
		s.resources[fullpath] = res
	}
  
	// Try to load it too
	res.Load()
	if res.Info() == nil {
		log.Println("Resource info is null: ", res.FullPath)
	}

	return res
}
```

Now we see two new interesting methods on the Resource - `Load` and `Info`.   The Info method is straightforward - it is a wrapper over `os.Stat` to cache load states and errors.


```go
func (r *Resource) Info() os.FileInfo {
	if r.info == nil {
		r.info, r.Error = os.Stat(r.FullPath)
		if r.Error != nil {
			r.State = ResourceStateFailed
			log.Println("Error Getting Info: ", r.FullPath, r.Error)
		}
	}
	return r.info
}
```

The `Load` method on the Resource type is more interesting - and more involved.  So we will defer this for now and come back to it.  Sorry!

```go
func (r *Resource) Load() *Resource {
  // TODO - All the interesting content processing bits
  return r
}
```

And with that we can complete our site's Load method:

```go
func (s *Site) Load() *Site {
	foundResources := s.ListResources(nil, nil, 0, 0)
	for _, res := range foundResources {
	  s.Rebuild(res)
  }
	return s
}

// The main method that builds a resource into the output folder
func (s *Site) Rebuild(res *Resource) {
	// var errors []error
	srcpath := res.FullPath
	if strings.HasSuffix(srcpath, ".md") ||
		strings.HasSuffix(srcpath, ".mdx") ||
		strings.HasSuffix(srcpath, ".html") ||
		strings.HasSuffix(srcpath, ".htm") {
		destpath := res.DestPathFor()
		outres := s.GetResource(destpath)
		if outres != nil {
			outres.EnsureDir()
			outfile, err := os.Create(outres.FullPath)
			if err != nil {
				log.Println("Error writing to: ", outres.FullPath, err)
				return
			}
			defer outfile.Close()

			// Copy the file over if it is a .md or a .html file
		}
	}
}
```

Our `Load` method is simple - it first loads *all* the resources, and then rebuilds each resource.  The `Rebuild` method - for now - copies any file with the appropriate extension (md, mdx, html, htm) into a particular location within the site's `OutputDir` folder.  The most interesting method here is the `DestPathFor` on the resource - that tells us what target path it should be copied over as:

```go
```

The EnsureDir method is simple and just ensures that a given resource's parent directory exists (creating them if needed):

```go
[func](func) (r *Resource) EnsureDir() {
	dirname := filepath.Dir(r.FullPath)
	if err := os.MkdirAll(dirname, 0755); err != nil {
		log.Println("Error creating dir: ", dirname, err)
	}
}
```
