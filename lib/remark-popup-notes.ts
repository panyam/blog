import { Parent, Root } from 'mdast'
const util = require('util')
import { MdxJsxFlowElement } from 'mdast-util-mdx'
// import { MdxFlowExpression } from 'mdast-util-mdx'
import { visit } from 'unist-util-visit'
import { getAttrib, parseMarkup } from './utils/common'
import { Plugin, Transformer } from 'unified'

// const inspect = (result: any) => util.inspect(result, false, null, true)
const NOTE_NODE_NAME = 'Note'

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
    const allNotesDefs = [] as any[]
    const allNotesLinks = [] as any[]

    // visit(tree, '', (node: any, index: number | null, parent: Parent | null) => { console.log("Foremost Node: ", util.inspect(node)) })
    visit(tree, 'code', (node: any, index: number | null, parent: Parent | null) => {
      if (node.lang != 'note' || parent == null || index == null) {
        return
      }
      console.log('In ZZZ: ', node, index)
      node.name = 'pre'
      node.type = 'mdxJsxFlowElement'
      node.attributes = [
        {
          type: 'mdxJsxAttribute',
          name: 'className',
          value: 'PopupNoteContentPre',
        },
        {
          type: 'mdxJsxAttribute',
          name: 'hidden',
          value: true,
        },
      ]
      node.children = [
        {
          name: 'code',
          type: 'mdxJsxFlowElement',
          children: [
            {
              name: 'text',
              type: 'text',
              value: node.value,
            },
          ],
        },
      ]
      console.log('Parent Afterwards: ', parent.children)
    })

    visit(tree, 'link', (node: any, index: number | null, parent: Parent | null) => {
      const url = node.url || ''
      if (!url.startsWith('#NOTE=') || index == null || parent == null) {
        return
      }
      const href = url.substring('#NOTE='.length)
      console.log('Here as LINK: ', node, parent, index)
      node.name = 'a'
      node.type = 'mdxJsxFlowElement'
      node.attributes = [
        {
          type: 'mdxJsxAttribute',
          name: 'className',
          value: 'PopupNoteAnchor',
        },
        {
          type: 'mdxJsxAttribute',
          name: 'noteref',
          value: href,
        },
      ]
    })

    // Now go through all links and turn them into clickable notes
    allNotesLinks.forEach((nl, index) => {
      //
    })
  }
}

const remarkPopupNotes: Plugin<[any], Root> = createTransformer
export default remarkPopupNotes
