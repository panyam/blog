import { OurMDXRemote } from '@/components/MDXComponents'
import { renderPostContent } from '@/lib/utils/renderer'
import { InferGetStaticPropsType } from 'next'
import ContentService from '@/lib/utils/contentservice'

const DEFAULT_LAYOUT = 'AuthorLayout'

export const getStaticProps = async () => {
  const contentSvc = await new ContentService().setup()
  const allAuthors = contentSvc.getAllAuthors()
  const author = allAuthors.find((p) => p.slug === 'authors/default')
  const mdxSource = await renderPostContent(author.body.raw)
  return { props: { author, mdxSource } }
}

export default function About({
  author,
  mdxSource,
}: InferGetStaticPropsType<typeof getStaticProps>) {
  return (
    <OurMDXRemote source={mdxSource} layout={author.layout || DEFAULT_LAYOUT} content={author} />
  )
}
