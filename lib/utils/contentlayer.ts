
import kebabCase from '@/lib/utils/kebabCase'
import { readFile, getAllFilesRecursively } from './files'
import path from 'path'
import dayjs from 'dayjs';
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
  console.log('Author Files: ', files)
  return []
}

export function getAllBlogs(folder?: string): Blog[] {
  folder = folder || 'data/blog'
  const files = getAllFilesRecursively(folder)
  const results: Blog[] = []
  files.forEach((file: string) => {
    console.log('Reading Post: ', file)
    const source = readFile(file)
    const matter = load_matter(file, source)

    // computed fields
    results.push(matter)
  })
  return results
}

export function load_matter(fullpath: string, source: string): any {
  const results = matter(source)
  const data = { ...results.data }
  let date = data.date
  if (typeof date !== 'number') {
    date = date.getTime() as number
  }
  data.dateString = dayjs(date).format('MMMM D, YYYY')
  data.readingTime = readingTime(source).text
  data.date = date
  data._raw = {
    sourceFilePath: fullpath,
    sourceFileDir: path.dirname(fullpath),
    contentType: 'mdx',
    sourceFileName: fullpath.split('/').pop(),
    flattenedPath: fullpath.replace(/\.[^/.]+$/, ''),
  }
  data.slug = data.slug || data._raw.flattenedPath.replace(/^.+?(\/)/, '')
  data.body = {
    raw: results.content,
  }

  console.log('data: ', data)
  console.log('_raw: ', data._raw)
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
