package main

import (
	"bytes"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"html/template"
	htmpl "html/template"
	ttmpl "text/template"

	"github.com/adrg/frontmatter"
	"github.com/gorilla/mux"
	gut "github.com/panyam/goutils/utils" // just a simple utility library
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"go.abhg.dev/goldmark/anchor"
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

func (s *Site) ListFiles(offset int, count int) (foundFiles []string) {
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

		foundFiles = append(foundFiles, fullpath)

		return nil // no errors found
	})
	if offset > 0 {
		foundFiles = foundFiles[offset:]
	}
	if count > 0 {
		foundFiles = foundFiles[:count]
	}
	return foundFiles
}

func (s *Site) Load() *Site {
	foundFiles := s.ListFiles(0, 0)
	for _, res := range foundFiles {
		s.Rebuild(res)
	}
	return s
}

// Rebuild a given src
func (s *Site) Rebuild(srcpath string) {
	destpath := s.EnsureDestDir(srcpath)
	if destpath == "" {
		return
	}

	outfile, err := os.Create(destpath)
	if err != nil {
		log.Println("Error writing to: ", destpath, err)
		return
	}
	defer outfile.Close()

	if strings.HasSuffix(srcpath, ".md") || strings.HasSuffix(srcpath, ".mdx") {
		// if we have a path of the form <ContentRoot>x/y/z.md,
		// then render it into <OutputDir>/x/y/z/index.html
		s.RenderMarkdown(srcpath, outfile)
	} else if strings.HasSuffix(srcpath, ".html") || strings.HasSuffix(srcpath, ".htm") {
		// if we have a path of the form <ContentRoot>x/y/z.html,
		// then render it into <OutputDir>/x/y/z/index.html
		s.RenderHtml(srcpath, outfile)
	} else {
		s.Copy(srcpath, outfile)
	}
}

// The EnsureDestDir method takes a file <ContentRoot>/a/b/c.<ext> and returns its "final form"
// in the <OutputDir>.   If the <ext> is .md or .mdx or .htm or .html, then this file is rendered
// into <OutputDir>/a/b/c/index.html
//
// However if c.<ext> is already of the form index.<ext> then another folder is not created.
// Eg if our srcpath was <ContentRoot>/a/b/c/index.<ext>, then it would be output into
// <OutputDir>/a/b/c/index.html
//
// Additionally the parent folders of the target file are also created if they do not exist
func (s *Site) EnsureDestDir(srcpath string) (destpath string) {
	// is this already an "index.<ext>" file?
	isIndex := strings.HasPrefix(srcpath, "index.") || strings.HasPrefix(srcpath, "_index.")
	needsIndex := strings.HasSuffix(srcpath, ".md") || strings.HasSuffix(srcpath, ".mdx") || strings.HasSuffix(srcpath, ".html") || strings.HasSuffix(srcpath, ".htm")

	respath, found := strings.CutPrefix(srcpath, s.ContentRoot)
	if found {
		return
	}
	if isIndex {
		destpath = filepath.Join(s.OutputDir, filepath.Dir(respath), "index.html")
	} else if needsIndex {
		ext := filepath.Ext(respath)
		rem := respath[:len(respath)-len(ext)]
		destpath = filepath.Join(s.OutputDir, rem, "index.html")
	} else {
		// we have a static file so copy as is
		destpath = filepath.Join(s.OutputDir, respath)
	}

	if destpath != "" {
		// make sure dir exists
		dirname := filepath.Dir(destpath)
		if err := os.MkdirAll(dirname, 0755); err != nil {
			log.Println("Error creating dir: ", dirname, err)
		}
	}
	return
}

// Loads a file given its path and returns a frontmatter (as a map) and the content byte array.
func (s *Site) LoadFile(srcpath string) (fm map[string]any, content []byte, err error) {
	f, err := os.Open(srcpath)
	if err == nil {
		fm = make(map[string]any)
		content, err = frontmatter.Parse(f, fm)
	}
	return
}

func (s *Site) RenderMarkdown(srcpath string, writer io.Writer) {
	fm, content, err := s.LoadFile(srcpath)
	if err != nil {
		log.Println("error loading file: ", err)
		return
	}

	// TODO - error checking
	tmplname := fm["template"].(string)
	template := s.GetHtmlTemplate()
	params := map[string]any{
		"FrontMatter": fm,
		"Content":     content,
	}
	template.ExecuteTemplate(writer, tmplname, params)
}

func (s *Site) RenderHtml(srcpath string, outfile io.Writer) {
	// TODO
}

func (s *Site) Copy(srcpath string, outfile io.Writer) {
	// TODO
}

func (s *Site) GetHtmlTemplate() (tmpl *htmpl.Template) {
	tmpl = htmpl.New("SiteHtmlTemplate").Funcs(s.DefaultFuncMap())
	for _, templatesDir := range s.HtmlTemplates {
		slog.Info("Loaded HTML Template: ", "templatesDir", templatesDir)
		t, err := tmpl.ParseGlob(templatesDir)
		if err != nil {
			slog.Error("Error parsing templates glob: ", "templatesDir", templatesDir, "error", err)
		} else {
			tmpl = t
			slog.Info("Loaded HTML Template: ", "templatesDir", templatesDir)
		}
	}
	return tmpl
}

func (s *Site) GetTextTemplate() (tmpl *ttmpl.Template) {
	tmpl = ttmpl.New("SiteHtmlTemplate").Funcs(s.DefaultFuncMap())
	for _, templatesDir := range s.TextTemplates {
		slog.Info("Loaded Text Template: ", "templatesDir", templatesDir)
		t, err := tmpl.ParseGlob(templatesDir)
		if err != nil {
			slog.Error("Error parsing templates glob: ", "templatesDir", templatesDir, "error", err)
		} else {
			tmpl = t
			slog.Info("Loaded HTML Template: ", "templatesDir", templatesDir)
		}
	}
	return tmpl
}

// Define all your custom functions here
func (s *Site) DefaultFuncMap() htmpl.FuncMap {
	return htmpl.FuncMap{
		"RenderMarkdown": func(fm map[string]any, content []byte) (template.HTML, error) {
			// TODO - handle errors
			mdTemplate, _ := s.GetTextTemplate().Parse(string(content))
			finalmd := bytes.NewBufferString("")
			mdTemplate.Execute(finalmd, map[string]any{
				"Site":        s,
				"FrontMatter": fm,
			})

			md := goldmark.New(
				goldmark.WithExtensions(
					extension.GFM,
					extension.Typographer,
					highlighting.NewHighlighting(
						highlighting.WithStyle("monokai"),
					),
					&anchor.Extender{},
				),
				goldmark.WithParserOptions(
					parser.WithAutoHeadingID(),
					parser.WithASTTransformers(),
				),
				goldmark.WithRendererOptions(
					html.WithHardWraps(),
					html.WithXHTML(),
					html.WithUnsafe(),
				),
			)
			var buf bytes.Buffer
			if err := md.Convert(finalmd.Bytes(), &buf); err != nil {
				slog.Error("error converting md: ", "error", err)
				return template.HTML(""), err
			}
			return template.HTML(buf.Bytes()), nil
		},
		// And any other methods
	}
}
