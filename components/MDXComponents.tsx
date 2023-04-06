/* eslint-disable react/display-name */
import React from 'react'
import Image from './Image'
import CustomLink from './Link'
import TOCInline from './TOCInline'
import Pre from './Pre'
import CodeRenderer from './CodeRenderer'
import { BlogNewsletterForm } from './NewsletterForm'

interface MDXLayout {
  layout: string
  content: any
  [key: string]: unknown
}

interface Wrapper {
  layout: string
  [key: string]: unknown
}

export const MDXComponents: any = {
  Image,
  TOCInline,
  a: CustomLink,
  pre: Pre,
  CodeEmbed: CodeRenderer,
  BlogNewsletterForm,
}

/*
export const MDXLayoutRenderer = ({ layout, content, ...rest }: MDXLayout) => {
  const MDXLayout = useMDXComponent(content.body.code)
  const mainContent = content // coreContent(content);

  return <MDXLayout layout={layout} content={mainContent} components={MDXComponents} {...rest} />
}
*/

// Using the MDXProvider version - with @next/mdx
/*
import { MDXProvider } from '@mdx-js/react'
export function OurMDXProvider({ layout, content, ...rest }) {
  console.log('Content: ', content)
  console.log('Layout: ', layout)
  return (
    <MDXProvider components={MDXComponents}>
      <Wrapper layout={layout} content={content}>
        {rest}
      </Wrapper>
    </MDXProvider>
  )
}
*/

// Using the mdx-next-remote plugin
import { MDXRemote } from 'next-mdx-remote'
export function OurMDXRemote({ source, layout, content, ...rest }) {
  const Layout = require(`../layouts/${layout}`).default

  // Instead of a fixed components dict, we are making it dynamic
  // so that we can pass in the layout as a prop
  const components = { ...MDXComponents }
  components.wrapper = (props) => {
    const Layout = require(`../layouts/${layout}`).default
    return (
      <Layout content={content} {...rest}>
        {props.children}
      </Layout>
    )
  }
  return (
    <MDXRemote {...source} components={components}>
      <Layout content={content} {...rest} />
    </MDXRemote>
  )
}
