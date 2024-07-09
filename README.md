
# My Personal Blog

This was started by clonning the amazing [Tailwind starter blog](https://tailwind-nextjs-starter-blog.vercel.app/)).   Once cloned I made the following changes:

* Removed Contentlayer -  Contentlayer was adding unnecessary complexity just to load files - instead going to a simple abstraction (lib/util/contentservice) was easier and dependency free.  For exampel Contentlayer in Next needed next-contentlayer but this had a hard dependency on Next V12.   Upgrades were breaking.  Getting rid of contentlayer also got my build size down by about 20%.
* Removed Preact - Experts can talk more about this.  This was also not working well with React v18 (at as of 4/2022) due to hydrateRoot being replaced by hydrate method that Preact seemed to depend on.   May circile back on this later.
* Added a CodeEmbed plugin for embedding github sources (at particular tags etc) as snippets during build time.

## Dev Tips

When changing css/tailwind.css or adding new tailwind components, do the following:

```
npx tailwindcss -i ./css/tailwind.css -o ./static/css/tailwind.css
```
