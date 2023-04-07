import siteMetadata from '../data/siteMetadata.js'
import fs from 'fs'
import path from 'path'
import ContentService from '../lib/utils/contentservice'

async function generate() {
  const contentsvc = await new ContentService().setup()
  const escaper = await import('./htmlEscaper')
  // RSS for blog post
  const escape = escaper.escape

  const generateRssItem = (post, siteMetadata) => `
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

  const generateRss = (posts, page, siteMetadata) => `
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

  const allBlogs = await contentsvc.getAllBlogs()
  if (allBlogs.length > 0) {
    const rss = generateRss(allBlogs, 'feed.xml', siteMetadata)
    fs.writeFileSync('./public/feed.xml', rss)
  }

  // RSS for tags
  if (allBlogs.length > 0) {
    const tags = await contentsvc.getAllTags(allBlogs)
    for (const tag of Object.keys(tags)) {
      const filteredPosts = allBlogs.filter(
        (post) => post.draft !== true && post.tags.map((t) => contentsvc.getSlug(t)).includes(tag)
      )
      const rss = generateRss(filteredPosts, `tags/${tag}/feed.xml`, siteMetadata)
      const rssPath = path.join('public', 'tags', tag)
      await fs.promises.mkdir(rssPath, { recursive: true })
      await fs.promises.writeFile(path.join(rssPath, 'feed.xml'), rss)
    }
  }
}

generate()
