import { writeFileSync } from 'fs'
import prettier from 'prettier'
import siteMetadata from '../data/siteMetadata.js'
import { globSync } from 'glob'
import ContentService from '../lib/utils/contentservice'

async function generate() {
  const prettierConfig = await prettier.resolveConfig('./.prettierrc.js')
  const svc = await new ContentService().setup()
  const allBlogs = svc.getAllBlogs()
  const contentPages = allBlogs
    .filter((x) => !x.draft && !x.canonicalUrl)
    .map((x) => `/${x._raw.flattenedPath}`)
  const pages = globSync(['pages/*.{js|tsx}', 'public/tags/**/*.xml'], {
    ignore: ['pages/_*.{js|tsx}', 'pages/api', 'pages/404.{js|tsx}'],
  })

  const sitemap = `
        <?xml version="1.0" encoding="UTF-8"?>
        <urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
            ${pages
              .concat(contentPages)
              .map((page) => {
                const path = page
                  .replace('pages/', '/')
                  .replace('public/', '/')
                  .replace('.js', '')
                  .replace('.mdx', '')
                  .replace('.md', '')
                  .replace('/feed.xml', '')
                const route = path === '/index' ? '' : path
                return `
                        <url>
                            <loc>${siteMetadata.siteUrl}${route}</loc>
                        </url>
                    `
              })
              .join('')}
        </urlset>
    `

  const formatted = prettier.format(sitemap, {
    ...prettierConfig,
    parser: 'html',
  })

  writeFileSync('public/sitemap.xml', formatted)
}

generate()
