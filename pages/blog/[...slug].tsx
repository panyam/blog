import PageTitle from '@/components/PageTitle'
import { OurMDXRemote } from '@/components/MDXComponents'
import { sortedBlogPost } from '@/lib/utils/contentlayer'
import { renderPostContent } from '@/lib/utils/renderer'
import { InferGetStaticPropsType } from 'next'
import { getPostAuthor, getAllBlogs } from 'lib/utils/contentlayer'

const DEFAULT_LAYOUT = 'PostLayout'

export async function getStaticPaths() {
  const allBlogs = getAllBlogs()
  const paths = allBlogs.map((p) => ({ params: { slug: p.slug.split('/') } }))
  return {
    paths: paths,
    fallback: false,
  }
}

export const getStaticProps = async ({ params }) => {
  const slug = (params.slug as string[]).join('/')
  const sortedPosts = sortedBlogPost(getAllBlogs())
  const postIndex = sortedPosts.findIndex((p) => p.slug === slug)

  const prevContent = sortedPosts[postIndex + 1] || null
  const prev = prevContent || null
  const nextContent = sortedPosts[postIndex - 1] || null
  const next = nextContent || null
  const post = sortedPosts.find((p) => p.slug === slug)
  const authorDetails = getPostAuthor(post)
  const mdxSource = await renderPostContent(post.body.raw)
  return {
    props: {
      post,
      mdxSource,
      authorDetails,
      prev,
      next,
    },
  }
}

export default function Blog({
  post,
  mdxSource,
  authorDetails,
  prev,
  next,
}: InferGetStaticPropsType<typeof getStaticProps>) {
  return (
    <>
      {'draft' in post && post.draft !== true ? (
        // <OurMDXProvider
        <OurMDXRemote
          source={mdxSource}
          layout={post.layout || DEFAULT_LAYOUT}
          content={post}
          toc={post.toc}
          authorDetails={authorDetails}
          prev={prev}
          next={next}
        />
      ) : (
        <div className="mt-24 text-center">
          <PageTitle>
            Under Construction{' '}
            <span role="img" aria-label="roadwork sign">
              🚧
            </span>
          </PageTitle>
        </div>
      )}
    </>
  )
}
