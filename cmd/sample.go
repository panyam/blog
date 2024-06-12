package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/gorilla/mux"
	gut "github.com/panyam/goutils/utils" // just a simple utility library
)

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

	// The router that serves the static resources and the built output files
	filesRouter *mux.Router

	// All resources encountered/used in our Site
	resources map[string]*Resource
}

func (s *Site) Init() {
	s.ContentRoot = gut.ExpandUserPath(s.ContentRoot)
	s.OutputDir = gut.ExpandUserPath(s.OutputDir)
	// our site's initializationa and build will go here
	s.Load()
}

var site = Site{
	ContentRoot: "../content",
	OutputDir:   "../output",
	PathPrefix:  "/ssgdemo",
	HtmlTemplates: []string{
		"../templates/*.html",
	},
	StaticFolders: []string{
		"/static/", "./static",
	},
}

func main() {
	site.Init() // Initialize it
	router := mux.NewRouter()
	router.PathPrefix(site.PathPrefix).Handler(http.StripPrefix(site.PathPrefix, &site))
	srv := &http.Server{
		Handler: router,
		Addr:    ":8888",
	}
	log.Fatal(srv.ListenAndServe())
}

func (s *Site) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// The entry point router for our site
	s.GetRouter().ServeHTTP(w, r)
}

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
			return nil
		}

		// map fullpath to a resource here
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

const (
	ResourceStatePending = iota
	ResourceStateLoaded
	ResourceStateDeleted
	ResourceStateNotFound
	ResourceStateFailed
)

/**
 * Each resource in our static site is identified by a unique path.
 * Note that resources can go through multiple transformations
 * resulting in more resources - to be converted into other resources.
 * Each resource is uniquely identified by its full path
 */
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

func (r *Resource) Load() *Resource {
	// TODO - All the interesting content processing bits
	return r
}

func (s *Site) Load() *Site {
	foundResources := s.ListResources(nil, nil, 0, 0)
	for _, res := range foundResources {
		s.Rebuild(res)
	}
	return s
}

// The main method that builds a list of resources
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

func (r *Resource) EnsureDir() {
	dirname := filepath.Dir(r.FullPath)
	if err := os.MkdirAll(dirname, 0755); err != nil {
		log.Println("Error creating dir: ", dirname, err)
	}
}

func (r *Resource) DestPathFor() (destpath string) {
	s := r.Site
	respath, found := strings.CutPrefix(r.FullPath, s.ContentRoot)
	if !found {
		log.Println("Respath not found: ", r.FullPath, s.ContentRoot)
		return ""
	}

	if r.Info().IsDir() {
		// Then this will be served with dest/index.html
		destpath = filepath.Join(s.OutputDir, respath)
	} else if r.IsIndex {
		destpath = filepath.Join(s.OutputDir, filepath.Dir(respath), "index.html")
	} else if r.NeedsIndex {
		// res is not a dir - eg it something like xyz.ext
		// depending on ext - if the ext is for a page file
		// then generate OutDir/xyz/index.html
		// otherwise OutDir/xyz.ext
		ext := filepath.Ext(respath)

		rem := respath[:len(respath)-len(ext)]

		// TODO - also see if there is a .<lang> prefix on rem after ext has been removed
		// can use that for language sites
		destpath = filepath.Join(s.OutputDir, rem, "index.html")
	} else {
		// basic static file - so copy as is
		destpath = filepath.Join(s.OutputDir, respath)
	}
	return
}
