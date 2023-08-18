import jsx from 'acorn-jsx'
import { Parser } from 'acorn'

export function getAttrib(node: any, attribName: string): any {
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
export function parseMarkup(value: string): any {
  try {
    const estree = parser.parse(value, { ecmaVersion: 'latest' })
    return {
      type: 'mdxFlowExpression',
      value,
      data: { estree },
    }
  } catch (e) {
    console.log('Markup Parse Error: ', e)
    throw e
  }
}
