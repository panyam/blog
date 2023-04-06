import jsx from 'acorn-jsx'
import { Parser } from 'acorn'
import { Parent, Root } from 'mdast'
import { MdxJsxFlowElement, MdxFlowExpression } from 'mdast-util-mdx'
import { visit } from 'unist-util-visit'
import { Plugin, Transformer } from 'unified'
import util from 'util'
const axios = require('axios')

const inspect = (result: any) => util.inspect(result, false, null, true)
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
  return async (tree: Parent & { lang?: string }) => {
    const allPromises = [] as Promise<string>[]

    visit(
      tree,
      'mdxJsxFlowElement',
      (node: MdxJsxFlowElement, index: number | null, parent: Parent | null) => {
        if (index == null || parent == null || node.name != 'div') {
          return
        }
        console.log('Div Node: ', inspect(node))
      }
    )

    visit(
      tree,
      'mdxJsxFlowElement',
      (node: MdxJsxFlowElement, index: number | null, parent: Parent) => {
        if (index == null || parent == null || node.name != OUR_NODE_NAME) {
          return
        }
        const url = getAttrib(node, 'url') || ''
        const lang = 'ts'
        allPromises.push(loadUrl(url, parent, index, lang, {}))
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
    const code = `<code><pre>Unable to laod url ${url}</pre></code>`
    parent.children[index] = parseMarkup(code)
  }
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
      {
        name: 'style',
        type: 'mdxJsxAttributeValueExpression',
        // value: { overflow: 'scroll', 'max-height': maxHeight },
        value: `max-height: ${maxHeight}; overflow: 'scroll';`,
        // value: { overflow: 'scroll', 'max-height': maxHeight },
        // value: `max-height: 300px; overflow: 'scroll';`,
      },
    ],
    children: [codeNode],
  }

  parent.children[index] = divNode
  return content
}

function getAttrib(node: any, attribName: string): any {
  const attrib = node.attributes.filter(
    (attr: any) => attr.type == 'mdxJsxAttribute' && attr.name == attribName
  )
  if (attrib.length == 0 || !attrib[0].value || attrib[0].value == null) {
    return null
  }
  if (typeof attrib[0].value === 'string') {
    return attrib[0].value
  } else {
    return attrib[0].value.value
  }
}

const parser = Parser.extend(jsx())
function parseMarkup(value: string): any {
  const estree = parser.parse(value, { ecmaVersion: 'latest' })
  return {
    type: 'mdxFlowExpression',
    value,
    data: { estree },
  }
}
