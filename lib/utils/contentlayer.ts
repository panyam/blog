import kebabCase from '@/lib/utils/kebabCase'
import getAllFilesRecursively from './files'
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
  console.log('Blog Files: ', files)
  return []
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
