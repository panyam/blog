import siteMetadata from '@/data/siteMetadata'
import ListLayout from '@/layouts/ListLayout'
import { PageSEO } from '@/components/SEO'
import ContentService from '@/lib/utils/contentservice'
import { InferGetStaticPropsType } from 'next'

export const POSTS_PER_PAGE = 10

export const getStaticProps = async () => {
  const contentSvc = await new ContentService().setup()
  const posts = contentSvc.sortedBlogPost()
  const initialDisplayPosts = posts.slice(0, POSTS_PER_PAGE)
  const pagination = {
    currentPage: 1,
    totalPages: Math.ceil(posts.length / POSTS_PER_PAGE),
  }

  return {
    props: {
      initialDisplayPosts: initialDisplayPosts,
      posts: posts,
      pagination,
    },
  }
}

export default function Blog({
  posts,
  initialDisplayPosts,
  pagination,
}: InferGetStaticPropsType<typeof getStaticProps>) {
  return (
    <>
      <PageSEO title={`Blog - ${siteMetadata.author}`} description={siteMetadata.description} />
      <ListLayout
        posts={posts}
        initialDisplayPosts={initialDisplayPosts}
        pagination={pagination}
        title="All Posts"
      />
    </>
  )
}
