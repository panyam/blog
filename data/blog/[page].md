---
title: BuildMage Blog
date: 2024-03-26T00:00:00AM
draft: false
images: []
authors: ['Sri Panyam']
page: HomePage
location: "BodyView.ContentView"
maxPosts: 20
---

{{ if eq paramname nil }}
{{ end }}

Sometimes we just want a parametric page - where we have a "semi" dynamic value but really provided during build time.
Eg we could have a 1000 blog posts but we may only want to show 25 of htem in a page - so we will need 40 of these pages
that all follow the same template.   So we want to create templatized pages like:

```
/path/a/b/c/[paramname].md
```

Here "paramname" determines the name of the parameter we use and the square brackets indicate that we shouldnt just create a 
"normal" output for this but instead create something like:

```
/outputfolder/<a>/index.html
/outputfolder/<b>/index.html
/outputfolder/<c>/index.html
...
```

So how many things lika "a", "b", "c" are needed should be calculateable by who?  Note this is an ordered list.

Option 1: From this md "return" some kind of function that tell show many pages there are.  Something like (in the front matter):


```
numPages: `ceil(len(.AllPages) / 20)`

or use normal templates:

numPages: {{ ceil(len(.AllPages) / 20) }}

and then evaluate the template
```

But this means we now have to start treating this in an interpretable way - can get messy.

Option 2: Generate pages dynamically.

Eg as a page is generated, see what are all the other pages linked from it and see if they are part of the site and generate them too (DFS).

This would need us to do a couple of things:

1. Parse the final output html and see links
2. Add part of the content processor that checks the AST for these link nodes
3. Simplest - annotate it in our template itself, eg we may be doing something like showing a "Next Page" or "Prev Page" link, here we could do:

```
<a href=/pages/{{n - 1}}>Prev</a>     <a href=/pages/{{n + 1}}>Next</a>
```

Instead we could do something like:

```
{{ genlink href="/pages/{{n - 1}} label = "Prev Page" }}

and same for next
```

"genlink" could be a separate template that notifies the parent what are generated pages from this template.

This also feels tricky as each page only knows about its immediate neighbors.

Option 3: Like option2 - but once

There is no reason the numPages has to be a frontmatter thing.  The page already knows everythign it needs to know, so why not add somethign like this:

```
{{ genpages ITEM1 ITEM2 ITEM3 .... ITEMN }}
```

Putting this anywhere in the page tells the caller that this page has to be run "N" times = each time with the paramname value set to ITEMX

In fact we can simplify this.   We could do:

```
{{ if paramname is None }}
  {{ genpages ITEM1 ITEM2 ITEM3 .... ITEMN }}
  // the caller knows that when paramname is none - it should *NOT* write to an output file but just to a null buffer
  // and look for "genpages" function being called
{{ else }}

  // Here put the code for how to render page with ITEMX
  // the caller knows that when it calls this with ITEMX it is alread writing it to /outputfolder/ITEMX/index.html

{{ end }}
```
