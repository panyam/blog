
We want to keep this file in a way so we can jot down quick notes so that when we actually write the blog it can look like we did this in steps.

So the way to do it is keep the following open:

1. Two side by side firefox windows - left window shows buildmage.com, right shows "current" blog during progress.
2. Tag of our code base at each "checkpoint" and in our static/public/images/gohtmx/TAG.png that shows what was the screenshot of the above two FF windows *at* that tag

At this point we could just go to the tag and see what the state of the repo was compared to "before" it (start is BEFORE_GOHTMX_MIGRATION) and describe then in our step by step blog

Or jot down her as we go (TAG as our headline from "so far we have this" kinda thing)

GOHTMX_STEP1 - Create base project to show a serve hello world

1. Create the basic files - get from diff from prev tag
2. Describe each file
3. key callouts on style:
  a. using web/* folder instead of putting all .go files at root
  b. web/templates folder to contain our templates
  c. serve hello world

GOHTMX_STEP2 - <Goal>

We have a bunch of things to do now:

1. Thing of which pages we have:
  a. Post listing page (paginated)
  b. Detail Post page
  c. Ability to add new pages
  d. /tags etc
  e. About page
2. Serve the listing/front page
3. Ensure our tailwind css from nextjs is ported and useable here
4. Folder structure for templates
5. Talk about components - and why we need some thing composeable even if on BE
6. Need for MDX!

Something incremental?   NextJS still used for page details - but we render them as templates first?
Or go directly to our own pages and FU NextJS - no need for htmx yet?

May be step 2 should just be - to get the template rendered to show our theme

Considerations

Ok this list is growing so let us keep scope down for step 2 - just show the front 

also things are looking pretty hopeless.   Our NextJS had tons of "components" and it was hard to make head or tail out of what was getting called how etc.

Looks like NextJS (v12) used the following convention:

pages are rendered via top level uRLs - eg

```
/               =====>    ./pages/index.tsx
/about          =====>  ./pages/about.tsx
/blog/<slug>        =====>  ./pages/blog/[...slug].tsx      // the acutal post given by slug
/blog/page/pagenum  =====>  ./pages/blog/page/[page].tsx    // paginated listing - by page num
```

and so on

But next js just takes the "body" content from these pages and renders it inside what ever the template is
in the `pages/_document.tsx` master template.

Then we have `_app.tsx` for client side rendering.  Frankly we dont need both and here is where our setup kicks in!

But we are rushing ahead - let us serve a basic page to serve the "listing" page

1. take our `_document` and convert to templates/index.html
