type Blog = any
type Author = any
type DocumentTypes = any
const DEFAULT_AUTHOR = 'authors/default'

export function dateSortDesc(a: string, b: string) {
  if (a > b) return -1
  if (a < b) return 1
  return 0
}

export default class ContentService {
  // All the dynamic imports we need
  imports = {} as any
  async setup(): Promise<this> {
    this.imports['fs'] = await import('fs')
    this.imports['path'] = await import('path')
    this.imports['matter'] = await import('gray-matter')
    this.imports['readingTime'] = await import('reading-time')
    // const slugger = await import('github-slugger')
    return this
  }

  readFile(file: string) {
    return this.imports.fs.readFileSync(file, 'utf-8')
  }

  getAllFilesRecursively(dirPath: string, arrayOfFiles?: string[]) {
    const fs = this.imports.fs
    const path = this.imports.path
    const files = fs.readdirSync(dirPath)
    // console.log('__dirname, Dir, Files: ', __dirname, dirPath, files)

    arrayOfFiles = arrayOfFiles || []
    for (const file of files) {
      if (fs.statSync(dirPath + '/' + file).isDirectory()) {
        arrayOfFiles = this.getAllFilesRecursively(dirPath + '/' + file, arrayOfFiles)
      } else if (arrayOfFiles) {
        arrayOfFiles.push(path.join(dirPath, '/', file))
      }
    }
    return arrayOfFiles
  }

  getAllAuthors(folder?: string): Author[] {
    folder = folder || 'data/authors'
    const files = this.getAllFilesRecursively(folder)
    const authors: Author[] = []
    // console.log('Author Files: ', files)
    for (const file of files) {
      // console.log('Reading Author: ', file)
      const source = this.readFile(file)
      const matter = this.parseAuthor(file, source)
      authors.push(matter)
    }
    return authors
  }

  getAllBlogs(folder?: string): Blog[] {
    folder = folder || 'data/blog'
    const files = this.getAllFilesRecursively(folder)
    const results: Blog[] = []
    for (const file of files) {
      const source = this.readFile(file)
      const matter = this.parseBlog(file, source)

      // computed fields
      results.push(matter)
    }
    return results
  }

  parseAuthor(fullpath: string, source: string): Promise<any> {
    const path = this.imports.path
    const matter = this.imports.matter
    const results = matter.default(source)
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

  parseBlog(fullpath: string, source: string): Promise<any> {
    // console.log('Reading Post: ', fullpath)
    const matter = this.imports.matter
    const path = this.imports.path
    const results = matter.default(source)
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
    const readingTime = this.imports.readingTime
    data.readingTime = readingTime.default(source).text
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

  sortedBlogPost(allBlogs?: Blog[]) {
    allBlogs = allBlogs || this.getAllBlogs()
    return allBlogs.sort((a, b) => dateSortDesc(a.date, b.date))
  }

  getSlug(tag): string {
    return tag
    return this.imports.slugger.slug(tag)
  }

  /**
   * Gets all blogs that contain either atleast one or all tags in the given tag
   * list.  This is controlled by the matchAllTags bool flag.
   */
  getBlogsByTags(blogs: Blog[] | null, tags: string[], matchAllTags = false) {
    blogs = blogs || this.getAllBlogs()
    if (blogs.length == 0) return []
    return blogs.filter((post) => {
      if (post.draft) return false
      const postTags = post.tags.map((t) => this.getSlug(t))
      for (const tag of tags) {
        if (postTags.includes(tag)) {
          if (!matchAllTags) return true
        } else if (matchAllTags) {
          return false
        }
      }
      // all tags matched
      return true
    })
  }

  getAllTags(allBlogs?: Blog[]): Record<string, number> {
    allBlogs = allBlogs || this.getAllBlogs()
    const tagCount: Record<string, number> = {}
    // Iterate through each post, putting all found tags into `tags`
    for (const file of allBlogs) {
      if (file.tags && file.draft !== true) {
        for (const tag of file.tags) {
          const formattedTag = this.getSlug(tag)
          if (formattedTag in tagCount) {
            tagCount[formattedTag] += 1
          } else {
            tagCount[formattedTag] = 1
          }
        }
      }
    }

    return tagCount
  }

  getPostAuthor(post: Blog): Promise<Author> {
    const authorList = post.authors || [DEFAULT_AUTHOR]
    const allAuthors = this.getAllAuthors()
    return authorList.map((author) => {
      return allAuthors.find((p) => p.slug === author) || {}
    })
  }
}
