/* eslint-disable react/display-name */
import React from 'react'
import Image from './Image'
import CustomLink from './Link'
import TOCInline from './TOCInline'
import Pre from './Pre'
import BookMark from './BookMark'
import { BlogNewsletterForm } from './NewsletterForm'

export const MDXComponents: any = {
  Image,
  TOCInline,
  a: CustomLink,
  pre: Pre,
  BookMark: BookMark,
  // CodeEmbed: CodeRenderer,
  BlogNewsletterForm,
}

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
