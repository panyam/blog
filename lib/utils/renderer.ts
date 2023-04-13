import { serialize } from 'next-mdx-remote/serialize'
import path from 'path'

// Remark packages
import remarkGfm from 'remark-gfm'
import remarkFootnotes from 'remark-footnotes'
import remarkMath from 'remark-math'
// we can also include "local" plugins!!
import remarkFrontmatter from 'remark-frontmatter'
import remarkExtractFrontmatter from '@/lib/remark-extract-frontmatter'
import remarkCodeTitles from '@/lib/remark-code-title'
import remarkImgToJsx from '@/lib/remark-img-to-jsx'
import remarkUrlEmbeds from '@/lib/remark-url-embeds'

// Rehype packages
import rehypeSlug from 'rehype-slug'
import rehypeAutolinkHeadings from 'rehype-autolink-headings'
import rehypeKatex from 'rehype-katex'
import rehypeCitation from 'rehype-citation'
import rehypePrismPlus from 'rehype-prism-plus'
import rehypePresetMinify from 'rehype-preset-minify'
// import rehypeHighlight from 'rehype-highlight'

export async function renderPostContent(content): Promise<any> {
  const mdxSource = await serialize(content, {
    mdxOptions: {
      // TODO - add remark snippets plugin and other configs here
      remarkPlugins: [
        remarkFrontmatter,
        remarkExtractFrontmatter,
        remarkGfm,
        remarkUrlEmbeds,
        remarkCodeTitles,
        [remarkFootnotes, { inlineNotes: true }],
        remarkMath,
        remarkImgToJsx,
        // remarkSnippetConfig,
      ],
      rehypePlugins: [
        rehypeSlug,
        [
          rehypeAutolinkHeadings,
          {
            properties: { className: ['anchor'] },
          },
          { behaviour: 'wrap' },
        ],
        rehypeKatex,
        // rehypeHighlight,
        [rehypeCitation, { path: path.join(process.cwd(), 'data') }],
        [rehypePrismPlus, { ignoreMissing: false }],
        rehypePresetMinify,
      ],
    },
  })
  return mdxSource
}
