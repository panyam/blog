import { TagSEO } from '@/components/SEO'
import siteMetadata from '@/data/siteMetadata'
import ListLayout from '@/layouts/ListLayout'
import { InferGetStaticPropsType } from 'next'
import ContentService from '@/lib/utils/contentservice'

export async function getStaticPaths() {
  const contentSvc = await new ContentService().setup()
  const tags = await contentSvc.getAllTags()

  return {
    paths: Object.keys(tags).map((tag) => ({
      params: {
        tag,
      },
    })),
    fallback: false,
  }
}

export const getStaticProps = async (context) => {
  const tag = context.params.tag as string
  const contentSvc = await new ContentService().setup()
  const filteredPosts = contentSvc.getBlogsByTags(null, [tag])
  return { props: { posts: filteredPosts, tag } }
}

export default function Tag({ posts, tag }: InferGetStaticPropsType<typeof getStaticProps>) {
  // Capitalize first letter and convert space to dash
  const title = tag[0].toUpperCase() + tag.split(' ').join('-').slice(1)
  return (
    <>
      <TagSEO
        title={`${tag} - ${siteMetadata.title}`}
        description={`${tag} tags - ${siteMetadata.author}`}
      />
      <ListLayout posts={posts} title={title} />
    </>
  )
}
