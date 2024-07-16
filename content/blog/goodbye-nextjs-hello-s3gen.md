---
title: 'Goodbye NextJS, Hello S3Gen'
date: 2024-05-28T11:29:10AM
tags: ['static site generator', 'go', 'ssg', 's3gen']
draft: false
authors: ['Sri Panyam']
summary: A brief journey on how I got off NextJS for hosting this blog and built a static site generator to host it instead.
page: PostPage
location: "BodyView.ContentView"
---

## History

This site is *very* old.   It has taken a few interesting paths:

In the beginning (about 25 years ago) it was a bunch of hand crafted html and js files - on Geocities!  This was a great time.   HTML/CSS/JS complexity was low.   Tables were a great choice for layouts (this humble back-end dev still thinks so at the risk of being laughed at).   You would edit a file, and upload the files (via ftp) or rsync to the hosted provider.   Then Geocities went bust and we were fending for ourselves.

Rise of Wordpress made it so much easy to create snazy themed blogs for oneself.  This site was no exception.  It moved to wordpress.com for a long time.  Infact on three different hosting providers (hosted wordpress, bluehost with a wordpress CMS server and one more that no longer resides in memory).  This was great for a long time.  Wysiwyg editor made it easy not to worry about layouts and formatting etc and choice of themes was pretty nice.

Wordpress had its problems.   The editing experience just felt clunky.   Also the type of supported content was limited.   You could embed images, videos and content (ofcoures).  However around around 6-7 years ago as I was becoming a lot more active with System Design preparation (for FAANG companies) and helping others with their preparation, I was looking for a content platform that could one day host very dynamic almost app-like content in a Blog form.   While I wrote a few blog posts (Il be surfacing them back up soon), I was struggling with supporting drawings (design drawings etc) and mathematical content (formulas etc).   I was also looking to host custom apps like system simulators on the page etc.   A "standard" editor provided by wordpress like sites was not cutting it.

Here I made the switch to building the site in NextJS.   The main advantages here were I could author all my pages as markdown - ie .md files.  Actually NextJS's plugin system allowed authoring in Extended Markdown (MDX) that ahd a larger and richer ecosystem of plugins and lot more options for plugability.   At the same time I had also moved my [Carnatic Music Notation](https://notations.us) website from server rendered pages on ExpressJS to also using NextJS and it was quite a liberating experience.  I could build as many custom components as I wanted (not that I had much of need for it beyond custom code embeding features - which we shall talk about soon).

I had gotten a bit busy and stopped writing for a while (both here as well as working on [Notations](https://notations.us).  And when I tried to get back into it I was having a few common problems across all my Node/Next apps.   Dependency problems.  For some reason Id see wierd dependency breakages where some package would be deprecated or be broken.  For example it was a nightmare migrating NextJS to the next version (i think v13) as its dependency (React) at some point in time was not updated at the same time freezing NextJS.   There were several such dependencies across several libraries.   Plus the build phase itself was pretty slow (often taking 10-15s on an Mac M1).   And then there was the bloat.   Each of these "distributions" was around a Gig when uploaded.  At this point I also started learning about htmx and the idea of going back to Server side generation first and *then* adding JS when needed was very appealing (as opposed to the otherway around in the React/Angular ecosystems).   All these got me thinking why not move to a static site generator (SSG) like Hugo or Jekyll but .... it is JUST a static site.  Why do I need a new tool for it.   Are static sites just not "build" tools to convert your content into html pages?  Thus began this journey of just creating my own SSG instead of depending on conventions imposed by these tools.   Yes Id have my own conventions but they are mine!  Now you can have yours too and they would be yours!

## Goals

Now that Ive made a few trips around the block I settled into a few basic requirements:

* We want to be able to write html (.html) or in markdown files (.md or .mdx).  Note that even though we support the .mdx extension for now we dont need Extended Markdown support as we shall see.
* Our system will be in Go - so we can enjoy an amazing standard library as well as a very powerful text and html templating system and we will see why this is a great thing.
* Like most other popular SSGs, we want to leverage directory structures to reflect http paths (eg content for `myblog.com/a/b/c` would be triggered from `<my content root>/a/b/c.md` or `<my content root>/a/b/c/index.{html,md}`)
* We want to be able to load "data" files and use content from those in our pages.  For example we may have a json file SiteMetadata.json that has some interesting info like twitter handles, github links, site titles etc that we want to reuse in a bunch of places.
* Since we are leveraging Go htm/text templates, we want the power of customizability and as such we want our content to be first class templates that will be rendered (within a series of layouts).
* Again since we are leveraging Go htm/text templates, being able to provide custom "functions" available in templates is very desirable.   Some of these functions could be very specific to *your own site*.
* Provision for custom static files to be packaged and bundled together.
* It should be very easy to build our site into a target folder with all the html/js/css files and also serve in dev mode (including live reloading of content changes)

## Non Goals

S3Gen is in no way a replacement for mature, famous and battle hardened static site generators like Hugo or Jekyll or NextJS.   This is intended mainly for those interested in building one themselves and providing them with one particular of way of doing it.   There are several ways (would love to hear more from yall).   Another impetus for creating S3Gen was that I wanted a static site onto which I could incrementally add dynamic content - mainly because I have gotten hooked onto HTMX and a very very light weight page builder/provider is very useful and hugely fun!

## Getting Started

TL;DR Here is the link to the git repo for [this blog](https://github.com/panyam/blog/tree/byenextjs) the [simple static site generator](https://github.com/panyam/s3gen) library powering this blog.

Let us use S3Gen to build and serve the site (and later we will dive into the internals of S3Gen).   Our folder structure is:

```
|--- main.go              <--- Main entry point
|--- content/             <--- The global data and pages in .md will be here
|--- templates/           <--- All our "base" templates will be here (more on this later)
|--- static/              <--- Files to be served statically
|--- output/              <--- Folder where all static pages are built and served from
     |--- index.html      <--- A very basic test page (only for now which we will replace)
```

We have three main folders (content, templates and static) as described above and one output folder (output) where all our build artifacts are stored so we can simply serve this as a static folder.   There are other supporting files/folders (for css generation, readme etc) but those are not important for now.

## Getting Started

First let us define our site using S3Gen in `web/main.go`:

```go
var site = s3.Site{
  ContentRoot: "./content",
  OutputDir:   "./build",
  // PathPrefix:  "/ourblog",
  HtmlTemplates: []string{
    "templates/*.html",
  },
  StaticFolders: []string{
    "/static/", "static",
  },
}
```

Our Site definition is pretty self explanatory.   We created a Site with some some key attributes - ContentRoot, HtmlTemplates locations, StaticFolders, PathPrefix, OutputDir etc.  The "PathPrefix" is interesting.  Instead of serving our site at the root (eg `http://<hostname>`), we could serve it from the "/ourblog" prefix - ie `http://<hostname>/ourblog` (our hostname here would be `localhost:8888`) (if we chose to)

Now we can serve this site (we are using the gorilla mux router - but not needed) and by register the site to be the http handler at `PathPrefix`.

```go
func main() {
  flag.Parse()
  
  // Only do build etc if this is in dev.
  // In production directly serve statically from the output dir
	if os.Getenv("APP_ENV") != "production" {
		site.CommonFuncMap = TemplateFunctions()
		site.NewViewFunc = NewView
		site.Watch()
	}

  // Attach our site to be at /`PathPrefix`
  // The Site will also take care of serving static files from /`PathPrefix`/static paths
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
```

For now ignore the NewViewFunc and CommonFuncMap attributes on the site.   These will be explained soon.

The lines:

```go
	if os.Getenv("APP_ENV") != "production" {
		site.CommonFuncMap = TemplateFunctions()
		site.NewViewFunc = NewView
		site.Watch()
	}
```

simply ensures that the site is setup for live reloading (which includes a one time full-build).   Note that this only runs in NON PRODUCTION mode because in production we are more interested in serving a statically built site instead of rebuilding etc.

That is it.  The `Site` is already a valid [`http.HandlerFunc`](https://pkg.go.dev/net/http#HandlerFunc) implementation so it can be served - and we are doing just that.  The Site also ensures that its http.HandlerFunc implementation contains routes for all the pages and any static content/folders registered (eg `/static` above).

## Extending common functions

S3Gen is a static site generator library.  As such it comes with a simple (and continually expanding) "standard library" of functions.   But custom sites can pass their own functions that can be called from their own sites (and templates).  We are doing just that by setting the `Site.CommonFuncMap` attribute so that our templates have access to these.  For example our template uses the `LeafPages` function to get all pages that are "leaf" pages - ie those that have a real .md or .html page backing it.  Or the `PagesByDate` function can be used to fetch all pages ordered by creation date. And more.

## Styling and CSS

This bit was kept light and is really very specific to this blog.   S3Gen does not care about styling - so you can plug it in as you like.  In our case we are using [tailwindcss](https://tailwindcss.com/).   Our tailwind css config is defined in tailwind.config.js and we build our "final" css with the command:

```
npx tailwindcss -i ./css/tailwind.css -o ./static/css/tailwind.css
```

This resultant css file is loaded from our base template - in `templates/CommonPageHeader.html`.

## Deployment

The `Watch` method builds all artifacts onto the `OutputDir` folder.   Static folders are not copied however (they could be?).   Also (inspired by Hugo), the `OutputDir` is not first deleted.  This allows any manual addition of resources into the output directory across builds.   We are using google app engine so our app.yaml file is pretty simple:

```
runtime: go122
env_variables:
  APP_ENV: production
handlers:
- url: /static
  static_dir: static
- url: /
  static_files: output/index.html
  upload: output/index.html
- url: /(.*)/
  static_files: output/\1/index.html
  upload: output/(.*)
- url: /(.*)
  static_files: output/\1/index.html
  upload: output/(.*)
- url: .*
  script: auto
```

Just serve the static files from the static folder and everything else via our main entry file (main.go).   Here the all the url patterns other than "/static" are served by `output/*` because each of your blog post of the form `ContentDir/X.{md,html}` is compiled to `OutputDir/X/index.html`.   If you use a different convention for your site simply update this app.yaml file.  Or if you use a cdn where you can upload files directly even better!

## And more

There is a lot more to S3Gen.  Head over to the [S3Gen documentation page](https://pkg.go.dev/github.com/panyam/s3gen).  The documentation is still WIP but it will be constantly updated so you can learn more about S3Gen's capabilities and what else is planned for the future.

