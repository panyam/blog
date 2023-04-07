import { PageSEO } from '@/components/SEO'
import siteMetadata from '@/data/siteMetadata'
import ListLayout from '@/layouts/ListLayout'
import { POSTS_PER_PAGE } from '../../blog'
import { InferGetStaticPropsType } from 'next'
import ContentService from '@/lib/utils/contentservice'

export const getStaticPaths = async () => {
  const contentSvc = await new ContentService().setup()
  const totalPosts = contentSvc.getAllBlogs()
  const totalPages = Math.ceil(totalPosts.length / POSTS_PER_PAGE)
  const paths = Array.from({ length: totalPages }, (_, i) => ({
    params: { page: (i + 1).toString() },
  }))

  return {
    paths,
    fallback: false,
  }
}

export const getStaticProps = async (context) => {
  const {
    params: { page },
  } = context
  const contentSvc = await new ContentService().setup()
  const posts = contentSvc.sortedBlogPost()
  const pageNumber = parseInt(page as string)
  const initialDisplayPosts = posts.slice(
    POSTS_PER_PAGE * (pageNumber - 1),
    POSTS_PER_PAGE * pageNumber
  )
  const pagination = {
    currentPage: pageNumber,
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

export default function PostPage({
  posts,
  initialDisplayPosts,
  pagination,
}: InferGetStaticPropsType<typeof getStaticProps>) {
  return (
    <>
      <PageSEO title={siteMetadata.title} description={siteMetadata.description} />
      <ListLayout
        posts={posts}
        initialDisplayPosts={initialDisplayPosts}
        pagination={pagination}
        title="All Posts"
      />
    </>
  )
}
