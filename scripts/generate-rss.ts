import { writeFileSync, mkdirSync } from 'fs'
import path from 'path'
// import GithubSlugger from 'github-slugger'
import { slug } from 'github-slugger'
import { getAllBlogs, getAllTags } from '../lib/utils/contentlayer'
import { escape } from './htmlEscaper'
import siteMetadata from '../data/siteMetadata.js'
// import { getAllBlogs, getAllTags } from '../lib/utils/contentlayer'

const generateRssItem = (post) => `
  <item>
    <guid>${siteMetadata.siteUrl}/blog/${post.slug}</guid>
    <title>${escape(post.title)}</title>
    <link>${siteMetadata.siteUrl}/blog/${post.slug}</link>
    ${post.summary && `<description>${escape(post.summary)}</description>`}
    <pubDate>${new Date(post.date).toUTCString()}</pubDate>
    <author>${siteMetadata.email} (${siteMetadata.author})</author>
    ${post.tags && post.tags.map((t) => `<category>${t}</category>`).join('')}
  </item>
`

const generateRss = (posts, page = 'feed.xml') => `
  <rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom">
    <channel>
      <title>${escape(siteMetadata.title)}</title>
      <link>${siteMetadata.siteUrl}/blog</link>
      <description>${escape(siteMetadata.description)}</description>
      <language>${siteMetadata.language}</language>
      <managingEditor>${siteMetadata.email} (${siteMetadata.author})</managingEditor>
      <webMaster>${siteMetadata.email} (${siteMetadata.author})</webMaster>
      <lastBuildDate>${new Date(posts[0].date).toUTCString()}</lastBuildDate>
      <atom:link href="${siteMetadata.siteUrl}/${page}" rel="self" type="application/rss+xml"/>
      ${posts.map(generateRssItem).join('')}
    </channel>
  </rss>
`

async function generate() {
  const allBlogs = getAllBlogs()
  // RSS for blog post
  if (allBlogs.length > 0) {
    const rss = generateRss(getAllBlogs())
    writeFileSync('./public/feed.xml', rss)
  }

  // RSS for tags
  // TODO: use AllTags from contentlayer when computed docs is ready
  if (allBlogs.length > 0) {
    const tags = await getAllTags(allBlogs)
    for (const tag of Object.keys(tags)) {
      const filteredPosts = getAllBlogs().filter(
        (post) => post.draft !== true && post.tags.map((t) => slug(t)).includes(tag)
      )
      const rss = generateRss(filteredPosts, `tags/${tag}/feed.xml`)
      const rssPath = path.join('public', 'tags', tag)
      mkdirSync(rssPath, { recursive: true })
      writeFileSync(path.join(rssPath, 'feed.xml'), rss)
    }
  }
}

generate()
