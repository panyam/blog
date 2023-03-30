import kebabCase from './kebabCase'
import { readFile, getAllFilesRecursively } from './files'
import path from 'path'
// import dayjs from 'dayjs'
import matter from 'gray-matter'
import readingTime from 'reading-time'

// import type { Blog, DocumentTypes } from 'contentlayer/generated'
//
type Blog = any
type Author = any
type DocumentTypes = any

export function dateSortDesc(a: string, b: string) {
  if (a > b) return -1
  if (a < b) return 1
  return 0
}

export function getAllAuthors(folder?: string): Author[] {
  folder = folder || 'data/authors'
  const files = getAllFilesRecursively(folder)
  const authors: Author[] = []
  console.log('Author Files: ', files)
  files.forEach((file: string) => {
    console.log('Reading Author: ', file)
    const source = readFile(file)
    const matter = parseAuthor(file, source)
    authors.push(matter)
  })
  return authors
}

export function getAllBlogs(folder?: string): Blog[] {
  folder = folder || 'data/blog'
  const files = getAllFilesRecursively(folder)
  const results: Blog[] = []
  files.forEach((file: string) => {
    const source = readFile(file)
    const matter = parseBlog(file, source)

    // computed fields
    results.push(matter)
  })
  return results
}

export function parseAuthor(fullpath: string, source: string): any {
  const results = matter(source)
  const data = { ...results.data }
  data.body = {
    raw: results.content,
  }
  data._raw = {
    sourceFilePath: fullpath,
    sourceFileDir: path.dirname(fullpath),
    contentType: 'mdx',
    sourceFileName: fullpath.split('/').pop(),
    flattenedPath: fullpath.replace(/\.[^/.]+$/, ''),
  }
  data.slug = data.slug || data._raw.flattenedPath.replace(/^.+?(\/)/, '')
  return data
}

export function parseBlog(fullpath: string, source: string): any {
  console.log('Reading Post: ', fullpath)
  const results = matter(source)
  const data = { ...results.data }
  data.date = data.date || Date.now()
  if (data.date) {
    if (typeof data.date === 'string') {
      data.date = Date.parse(data.date)
    }
    if (typeof data.date !== 'number') {
      data.date = data.date.getTime() as number
    }
  }
  if (data.lastmod) {
    if (typeof data.lastmod === 'string') {
      data.lastmod = Date.parse(data.lastmod)
    }
    if (typeof data.lastmod !== 'number') {
      data.lastmod = data.lastmod.getTime() as number
    }
  }
  // data.dateString = dayjs(date).format('MMMM D, YYYY')
  data.readingTime = readingTime(source).text
  data._raw = {
    sourceFilePath: fullpath,
    sourceFileDir: path.dirname(fullpath),
    contentType: 'mdx',
    sourceFileName: fullpath.split('/').pop(),
    flattenedPath: fullpath.replace(/\.[^/.]+$/, ''),
  }
  data.slug = data.slug || data._raw.flattenedPath.split('/').slice(2).join('/') // replace(/^.+?(\/)/, '')
  data.isDir = false
  /*
  console.log('Read Slug: ', data.slug, data._raw)
  if (data.slug.endsWith('/index/')) {
    data.slug = data.slug.substring(0, data.slug.length - ('/index/'.length - 1))
    data.isDir = true
  } else if (data.slug.endsWith('/index')) {
    data.slug = data.slug.substring(0, data.slug.length - ('/index'.length - 1))
    data.isDir = true
  }
  */
  data.body = {
    raw: results.content,
  }
  return data
}

export function sortedBlogPost(allBlogs: Blog[]) {
  return allBlogs.sort((a, b) => dateSortDesc(a.date, b.date))
}

export async function getAllTags(allBlogs: Blog[]) {
  const tagCount: Record<string, number> = {}
  // Iterate through each post, putting all found tags into `tags`
  allBlogs.forEach((file) => {
    if (file.tags && file.draft !== true) {
      file.tags.forEach((tag) => {
        const formattedTag = kebabCase(tag)
        if (formattedTag in tagCount) {
          tagCount[formattedTag] += 1
        } else {
          tagCount[formattedTag] = 1
        }
      })
    }
  })

  return tagCount
}

const DEFAULT_AUTHOR = 'authors/default'
export function getPostAuthor(post: Blog): Author {
  const authorList = post.authors || [DEFAULT_AUTHOR]
  const allAuthors = getAllAuthors()
  return authorList.map((author) => {
    return allAuthors.find((p) => p.slug === author) || {}
  })
}
