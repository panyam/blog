import { Parent, Root } from 'mdast'
import { MdxJsxFlowElement } from 'mdast-util-mdx'
// import { MdxFlowExpression } from 'mdast-util-mdx'
import { visit } from 'unist-util-visit'
import { Plugin, Transformer } from 'unified'
import { getAttrib, parseMarkup } from './utils/common'
// import util from 'util'
const axios = require('axios')

// const inspect = (result: any) => util.inspect(result, false, null, true)
const OUR_NODE_NAME = 'CodeEmbed'

/**
 * Replaces code snippets referring to github embeds with actual file contents
 * The following options are allowed:
 *
 * url: URL of file to embed - Needs to be public and viewable
 * version: tag:<tagname> | branch:<branchname> | commit:<commitid>
 * start_line: Which line to start showing from
 * num_lines: Max number of lines to show (<= 50)
 */

function createTransformer(options: any): Transformer<Root> {
  options = options || {}
  return async (tree: Parent & { lang?: string }) => {
    const allPromises = [] as Promise<string>[]

    visit(
      tree,
      'mdxJsxFlowElement',
      (node: MdxJsxFlowElement, index: number | null, parent: Parent | null) => {
        if (index == null || parent == null || node.name != OUR_NODE_NAME) {
          return
        }
        const url = getAttrib(node, 'url') || ''
        const lang = getAttrib(node, 'language') || 'ts'
        const title = getAttrib(node, 'title') || ''
        const height = getAttrib(node, 'height') || options.defaultCodeHeight || '300px'
        // console.log('Div Node: ', inspect(node))
        allPromises.push(
          loadUrl(url, parent, index, lang, {
            title: title,
            height: height,
          })
        )
      }
    )

    await Promise.all(allPromises)
  }
}

const remarkGithubEmbeds: Plugin<[any], Root> = createTransformer
export default remarkGithubEmbeds

async function loadUrl(
  url: string,
  parent,
  index: number,
  lang: string,
  configs: any
): Promise<string> {
  configs = configs || {}
  const maxHeight = configs.maxHeight || '400px'
  if (url == null) {
    // replace with a warning
    const code = `<pre><code>Unable to laod url ${url}</code></pre>`
    parent.children[index] = parseMarkup(code)
  } else {
    const response = await axios.get(url)
    const content = response.data
    const codeNode = {
      type: 'code',
      value: content,
      meta: 'showLineNumbers',
      lang: lang,
    }
    const divNode = {
      name: 'div',
      type: 'mdxJsxFlowElement',
      attributes: [
        {
          type: 'mdxJsxAttribute',
          name: 'className',
          value: 'CodeEmbedContainer',
        },
      ],
      children: [
        {
          name: 'h4',
          type: 'mdxJsxFlowElement',
          attributes: [
            {
              type: 'mdxJsxAttribute',
              name: 'className',
              value: 'CodeEmbedHeading',
            },
          ],
          children: [
            {
              name: 'a',
              type: 'mdxJsxFlowElement',
              attributes: [
                {
                  type: 'mdxJsxAttribute',
                  name: 'className',
                  value: 'CodeEmbedUrlLink',
                },
                {
                  type: 'mdxJsxAttribute',
                  name: 'href',
                  value: url,
                },
              ],
              children: [
                {
                  type: 'text',
                  value: configs.title || '',
                },
              ],
            },
          ],
        },
        {
          name: 'div',
          type: 'mdxJsxFlowElement',
          attributes: [
            {
              type: 'mdxJsxAttribute',
              name: 'className',
              value: 'CodeEmbedBlockContainer',
            },
            {
              name: 'style',
              type: 'mdxJsxAttributeValueExpression',
              value: `max-height: ${maxHeight}; overflow: scroll;`,
              // type: 'mdxJsxAttribute',
              // value: "{ height: '300px', overflow: 'scroll' }",
              // value: { overflow: 'scroll', 'max-height': maxHeight },
              // value: { overflow: 'scroll', 'max-height': maxHeight },
              // value: `max-height: 300px; overflow: 'scroll';`,
            },
          ],
          children: [codeNode],
        },
      ],
    }
    parent.children[index] = divNode
    return content
  }
  return ''
}

