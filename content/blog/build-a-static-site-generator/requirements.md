---
title: 'Build a static site generator (SSG) - Requirements'
date: 2024-05-25T11:29:10AM
tags: ['static site generator', 'go']
draft: false
authors: ['Sri Panyam']
summary: First part of the static site generator series detailing the requirement and current setup.
page: PostPage
location: "BodyView.ContentView"
---

## History

This is the [building a static site genrator series](../) for hosting this blog!

This site is *very* old.   It has taken a few interesting paths:

In the beginning of time it was a bunch of hand crafted (partly with love) html and js files - on Geocities!  This was a great time.   HTML/CSS/JS complexity was low.   Tables were a great choice for layouts (this humble back-end still thinks so at the risk of being laughed at).   You would edit a file, and upload the files (via ftp) or rsync to the hosted provider.   Then Geocities went bust and we were fending for ourselves.

Rise of Wordpress made it so much easy to create snazy themed blogs for oneself.  This site was no exception.  It moved to wordpress.com for a long time.  Infact on three different hosting providers (hosted wordpress, bluehost with a wordpress CMS server and one more that no longer resides in memory).  This was great for a long time.  Wysiwyg editor made it easy not to worry about layouts and formatting etc and choice of themes was pretty nice.

Wordpress had its problems.   The editing experience just felt clunky.   Also the type of supported content was limited.   You could embed images, videos and content (ofcoures).  However around around 6-7 years ago as I was becoming a lot more active with System Design preparation (for FAANG companies) and helping others with their preparation, I was looking for a content platform that could one day host very dynamic almost app-like content in a Blog form.   While I wrote a few blog posts (Il be surfacing them back up soon), I was struggling with supporting drawings (design drawings etc) and mathematical content (formulas etc).   I was also looking to host custom apps like system simulators on the page etc.   A "standard" editor provided by wordpress like sites was not cutting it.

Here I made the switch to building the site in NextJS.   The main advantages here were I could author all my pages as markdown - ie .md files.  Actually NextJS's plugin system allowed authoring in Extended Markdown (MDX) that ahd a larger and richer ecosystem of plugins and lot more options for plugability.   At the same time I had also moved my [Carnatic Music Notation](https://notations.us) website from server rendered pages on ExpressJS to also using NextJS and it was quite a liberating experience.  I could build as many custom components as I wanted (not that I had much of need for it beyond custom code embeding features - which we shall talk about soon).

I had gotten a bit busy and stopped writing for a while (both here as well as working on Notations).  And when I tried to get back I was having a few common problems across all my Node/Next apps.   Dependency problems.  For some reason Id see wierd dependency breakages where some package would be deprecated or be broken.  For example it was a nightmare migrating NextJS to the next version as its dependency (React) at some point in time was not updated at the same time freezing NextJS.   There were several such dependencies across all the libraries.   Plus the build phase itself was pretty slow (often taking 10-15s on an Mac M1).   And then there was the bloat.   Each of these "distributions" was around a Gig when uploaded.  At this point I also started learning about htmx and the idea of going back to Server side generation first and *then* adding JS when needed was very appealing (as opposed to the otherway around in the React/Angular ecosystems).   All these got me thinking why not move to a static site generator (SSG) like Hugo or Jekyll but .... it is JUST a static site.  Why do I need a new tool for it.   Are static sites just not "build" tools to convert your content into html pages?  Thus began this journey of just creating my own SSG instead of depending on conventions imposed by these tools.   Yes Id have my own conventions but they are mine!  Now you can have yours too and they would be yours!


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

We are building out a simple SSG library.   In this part we simply lay out the requirements and our site structure.  In the [next part](../getting-started) we will start building it out.
