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

## History

This site is *very* old.   It has taken a few interesting paths:

In the beginning of time it was a bunch of hand crafted (partly with love) html and js files - on Geocities!  This was a great time.   HTML/CSS/JS complexity was low.   Tables were a great choice for layouts (this humble back-end still thinks so at the risk of being laughed at).   You would edit a file, and upload the files (via ftp) or rsync to the hosted provider.   Then Geocities went bust and we were fending for ourselves.

Rise of Wordpress made it so much easy to create snazy themed blogs for oneself.  This site was no exception.  It moved to wordpress.com for a long time.  Infact on three different hosting providers (hosted wordpress, bluehost with a wordpress CMS server and one more that no longer resides in memory).  This was great for a long time.  Wysiwyg editor made it easy not to worry about layouts and formatting etc and choice of themes was pretty nice.

Wordpress had its problems.   The editing experience just felt clunky.   Also the type of supported content was limited.   You could embed images, videos and content (ofcoures).  However around around 6-7 years ago as I was becoming a lot more active with System Design preparation (for FAANG companies) and helping others with their preparation, I was looking for a content platform that could one day host very dynamic almost app-like content in a Blog form.   While I wrote a few blog posts (Il be surfacing them back up soon), I was struggling with supporting drawings (design drawings etc) and mathematical content (formulas etc).   I was also looking to host custom apps like system simulators on the page etc.   A "standard" editor provided by wordpress like sites was not cutting it.

Here I made the switch to building the site in NextJS.   The main advantages here were I could author all my pages as markdown - ie .md files.  Actually NextJS's plugin system allowed authoring in Extended Markdown (MDX) that ahd a larger and richer ecosystem of plugins and lot more options for plugability.   At the same time I had also moved my [Carnatic Music Notation](https://notations.us) website from server rendered pages on ExpressJS to also using NextJS and it was quite a liberating experience.  I could build as many custom components as I wanted (not that I had much of need for it beyond custom code embeding features - which we shall talk about soon).

I had gotten a bit busy and stopped writing for a while (both here as well as working on Notations).  And when I tried to get back into it I was having a few common problems across all my Node/Next apps.   Dependency problems.  For some reason Id see wierd dependency breakages where some package would be deprecated or be broken.  For example it was a nightmare migrating NextJS to the next version as its dependency (React) at some point in time was not updated at the same time freezing NextJS.   There were several such dependencies across several libraries.   Plus the build phase itself was pretty slow (often taking 10-15s on an Mac M1).   And then there was the bloat.   Each of these "distributions" was around a Gig when uploaded.  At this point I also started learning about htmx and the idea of going back to Server side generation first and *then* adding JS when needed was very appealing (as opposed to the otherway around in the React/Angular ecosystems).   All these got me thinking why not move to a static site generator (SSG) like Hugo or Jekyll but .... it is JUST a static site.  Why do I need a new tool for it.   Are static sites just not "build" tools to convert your content into html pages?  Thus began this journey of just creating my own SSG instead of depending on conventions imposed by these tools.   Yes Id have my own conventions but they are mine!  Now you can have yours too and they would be yours!


## Requirements

* We want to be able to write html (.html) or in markdown files (.md or .mdx).  Note that even though we support the .mdx extension for now we dont need Extended Markdown support as we shall see.
* Our system will be in Go - so we can enjoy an amazing standard library as well as a very powerful text and html templating system and we will see why this is a great thing.
* Like most other popular SSGs, we want to leverage directory structures to reflect http paths (eg content for `myblog.com/a/b/c` would be triggered from `<my content root>/a/b/c.md` or `<my content root>/a/b/c/index.{html,md}`)
* We want to be able to load "data" files and use content from those in our pages.  For example we may have a json file SiteMetadata.json that has some interesting info like twitter handles, github links, site titles etc that we want to reuse in a bunch of places.
* Since we are leveraging Go htm/text templates, we want the power of customizability and as such we want our content to be first class templates that will be rendered (within a series of layouts).
* Again since we are leveraging Go htm/text templates, being able to provide custom "functions" available in templates is very desirable.   Some of these functions could be very specific to *your own site*.
* Provision for custom static files to be packaged and bundled together.
* It should be very easy to build our site into a target folder with all the html/js/css files and also serve in dev mode (including live reloading of content changes)

## Getting Started

TL;DR Here is the link to the git repo for [this blog](https://github.com/panyam/blog) the [simple static site generator](https://github.com/panyam/s3gen) library powering this blog.

Now let us see how to actually build up to this step by step.   Our folder structure is:

```
|--- content/             <--- The global data and pages in .md will be here
|--- templates/           <--- All our "base" templates will be here (more on this later)
|--- static/              <--- Files to be served statically
|--- output/              <--- Folder where all static pages are built and served from
     |--- index.html      <--- A very basic test page (only for now which we will replace)
|--- cmd/sample.go        <--- The code samples built in this post
```

We have three main folders (content, templates and static) as described above and one output folder (build) where all our build artifacts are stored so we can simply serve this as a static folder.   The code samples in this blog will be in the `cmd/sample.go` folder and can run with `go run cmd/sample.go`.

Our goal is to build out a simple SSG library.

## Getting Started

First start with a "Site" type in (`cmd/sample.go`) to represent some info about our Site, where files are stored and more:

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
  //    myblog.com                => PathPrefix = "/"
  //    myblog.com/blog            => PathPrefix = "/blog"
  //    myblog.com/blogs/blog1    => PathPrefix = "/blogs/blog1"
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

func (s *Site) Load() {
  // TODO
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

Once we have this we can simply process the files and render them into the output folder.  We modify the Site's `Load` method to perform this:


```go
func (s *Site) Load() *Site {
  foundFiles := s.ListFiles(nil, nil, 0, 0)
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
```


The method `EnsureDestDir` finds the corresponding output file path for a given input file path (in the `ContentRoot` folder).  Our rules are pretty simple (and here depending on your own site and conventions you could adopt a more comprehensive set of rules):

* If our source file is of the form `<ContentRoot>/a/b/c.<ext>` where c is not "index" or "_index", then the destination path is `<OutputDir>/a/b/c/index.html`.  This will ensure that this page will be rendered at the url path `http://<mysite>/a/b/c/`.
* If c *is* "index" or "_index", then our destination path is just `<OutputDir>/a/b/index.html` to be served at `http://<mysite>/a/b/`.
* Finally this method also ensures that the parent folder (in OutputDir) of the destination path is also created if it does not exist.

```go
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


func (s *Site) RenderMarkdown(srcpath string, outfile io.Writer) {
  // TODO
}

func (s *Site) RenderHtml(srcpath string, outfile io.Writer) {
  // TODO
}

func (s *Site) Copy(srcpath string, outfile io.Writer) {
  // TODO
}
```

Next we shall look at the rendering of specific file types.

## Rendering

In the previous section we chose to render either .md or .html files.  Any other file is just copied to its respective `OutputDir` location.   We will look at each of the rendering types now.

### Common Structure 

First the common convention.   As is common with static pages, each of our source files (.md or .html) carries [`FrontMatter`](https://markdoc.dev/docs/frontmatter).   Frontmatter lets you add page level metadata to documents.   It is a block of YAML, or TOML or JSON content delimited by "---" before and after the block.  For example:

```markdown
---
title: My post title
summary: This is a simple post.
---

Rest of our content begins and continues here.
```

So for each of our .md or .html files we will seperate the Frontmatter from the actual page content and *then* render the page content (by also passing the Frontmatter as a variable to it - this way the render and/or template can use page specific info during rendering if needed).  The following method does just that:

```go
// Loads a file given its path and returns a frontmatter (as a map) and the content byte array.
func (s *Site) LoadFile(srcpath string) (fm map[string]any, content []byte, err error) {
	f, err := os.Open(srcpath)
	if err == nil {
		fm = make(map[string]any)
		content, err = frontmatter.Parse(f, fm)
	}
	return
}
```

### Rendering Hierarchy

Both our .md and .html files have a common structure as discussed above.  How do we render the content?   For example our page could be richly themed with headers, footers, side bars, SEO metadata etc and finally somewhere in the page may be a "content div" that hosts our page (the actual article).   After all the idea is we want to have a uniform look and feel for all pages in our site without having to repeat the layout again and again.   For now we will use a very basic convention - and you can try various flavors of this for your own site and its conventions.


1. Once we load a file we shall look for a variable in its Frontmatter called "template".  Remember in our "templates" folder we can have several html templates for different parts of our site (Projects section, Blog posts, Authors section etc).  By using this variable we can select which template to render into.
2. Each template can have an arbitrary hierarchy of html content.  The template itself may include/reference/use other templates.  So data can be hierarchical.   When we render a template we usually pass it a rendering-context dictionary (more on this later).   Our post content (following the frontmatter) needs to be rendered at some custom hierarchy.   We can use another variable "fieldpath" to indicate the nested key in the context to use to pass the content via.   An example will make this clear.

Let us assume we have the following "BasePage" template (defined in `templates/BasePage.html`):

```html
{{ "{{ define \"BasePage\" }}" }}
  <html>
    <head>
    </head>
    <body>
      <header><h2>{{"{{.FrontMatter.Title}}"}}</h2></header>
      <article>
        {{ "{{ RenderMarkdown .Content .FrontMatter }} " }}
      </article>
      <footer>Copyright 2024 Bruce Wayne</footer>
    </body>
  </html>
{{ "{{ end }}" }}
```

The name of this template is "BasePage" so one of our posts might have the following frontmatter:

```markdown
---
Title: My page title
Summary: Here is a summary for this page.
template: BasePage
---

Did you know this is just normal markdown which will be rendered into html.
```

Our RenderMarkdown method would need to pass the frontmatter as the "FrontMatter" parameter and the post content via the "Content" parameter so it is picked up by the above template correctly:

```go
func (s *Site) RenderMarkdown(srcpath string, outfile io.Writer) {
  fm, content, err := s.LoadFile(srcpath)
  if err != nil {
    log.Println("error loading file: ", err)
    return
  }
  
  // TODO - error checking
  tmplname := fm["template"].(string)
  template := s.GetHtmlTemplate()
  params := map[string]any {
    "Site": s,
    "FrontMatter": fm,
    "Content": content,
  }
  template.ExecuteTemplate(writer, tmplname, params)
}
```

That's it.   The right template is selected, the site's "Parsed Template" is invoked with this templatename along with the FrontMatter and page content and written to the output file.

The "GetTemplate" method is the place for you to create you `html/template.Template` instances.  Here we can also add any custom functions you need etc (and these would be visible to the templates to use as they see fit).   Here is a very basic implementation of the template creator and a method to add any custom function needed:

```go
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
```

The `RenderHtml` method is very similar.
