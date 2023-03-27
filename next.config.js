import nextMDX from '@next/mdx'

const withBundleAnalyzer = require('@next/bundle-analyzer')({
  enabled: process.env.ANALYZE === 'true',
})

// You might need to insert additional domains in script-src
// if you are using external services
const ContentSecurityPolicy = `
  default-src 'self';
  script-src 'self' 'unsafe-eval' 'unsafe-inline' giscus.app;
  style-src 'self' 'unsafe-inline';
  img-src * blob: data:;
  media-src 'none';
  connect-src *;
  font-src 'self';
  frame-src giscus.app
`

const securityHeaders = [
  // https://developer.mozilla.org/en-US/docs/Web/HTTP/CSP
  {
    key: 'Content-Security-Policy',
    value: ContentSecurityPolicy.replace(/\n/g, ''),
  },
  // https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Referrer-Policy
  {
    key: 'Referrer-Policy',
    value: 'strict-origin-when-cross-origin',
  },
  // https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/X-Frame-Options
  {
    key: 'X-Frame-Options',
    value: 'DENY',
  },
  // https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/X-Content-Type-Options
  {
    key: 'X-Content-Type-Options',
    value: 'nosniff',
  },
  // https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/X-DNS-Prefetch-Control
  {
    key: 'X-DNS-Prefetch-Control',
    value: 'on',
  },
  // https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Strict-Transport-Security
  {
    key: 'Strict-Transport-Security',
    value: 'max-age=31536000; includeSubDomains',
  },
  // https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Feature-Policy
  {
    key: 'Permissions-Policy',
    value: 'camera=(), microphone=(), geolocation=()',
  },
]

import remarkFrontmatter from 'remark-frontmatter'
// const rehypeHighlight = import("rehype-highlight");
import rehypeHighlight from 'rehype-highlight'
import remarkSnippets from 'remark-snippets' // const remarkSnippets = import("remark-snippets");
// import remarkFrontmatter from "remark-frontmatter";
// import rehypeHighlight from "rehype-highlight";
const withMDX = nextMDX({
  options: {
    // If you use remark-gfm, you'll need to use next.config.mjs
    // as the package is ESM only
    // https://github.com/remarkjs/remark-gfm#install

    // remarkPlugins: [remarkFrontmatter],
    // rehypePlugins: [rehypeHighlight],
    remarkPlugins: [
      remarkFrontmatter,
      [
        remarkSnippets,
        {
          /**
           * Address of the snippet service to be invokved to run our snippets.
           */
          snippetsvc: {
            addr: 'localhost:7000', // default
          },
          /**
           * Different environments that can be used so they dont have to
           * be defined in the mdx files.  These environments can be overridden
           * in specific mdx files and new ones can also be created.
           */
          envinfo: {
            default: 'default',
            envs: [
              {
                name: 'default',
                packages: [
                  {
                    tlex: '*',
                  },
                ],
              },
            ],
          },
          foo: 'bar',
        },
      ],
    ],
    rehypePlugins: [rehypeHighlight],

    // If you use `MDXProvider`, uncomment the following line.
    // providerImportSource: "@mdx-js/react",
  },
})

/**
 * @type {import('next/dist/next-server/server/config').NextConfig}
 **/
module.exports = withMDX(
  withBundleAnalyzer({
    basePath: '/',
    reactStrictMode: true,
    trailingSlash: true,
    productionBrowserSourceMaps: true,
    distDir: 'build',
    pageExtensions: ['ts', 'tsx', 'js', 'jsx', 'md', 'mdx'],
    eslint: {
      dirs: ['pages', 'components', 'lib', 'layouts', 'scripts'],
    },
    webpack: (config, { dev, isServer }) => {
      config.module.rules.push({
        test: /\.svg$/,
        use: ['@svgr/webpack'],
      })

      if (!dev && !isServer) {
        // Replace React with Preact only in client production build
        Object.assign(config.resolve.alias, {
          'react/jsx-runtime.js': 'preact/compat/jsx-runtime',
          react: 'preact/compat',
          'react-dom/test-utils': 'preact/test-utils',
          'react-dom': 'preact/compat',
        })
      }

      return config
    },
    // https://github.com/vercel/next.js/issues/21079
    // Remove this workaround whenever the issue is fixed
    // images: { unoptimized: true, },
    images: {
      loader: 'imgix',
      path: '/',
      remotePatterns: [
        {
          protocol: 'https',
          hostname: 'raw.githubusercontent.com',
          port: '',
          pathname: '/grpc-ecosystem/**',
        },
      ],
    },
    async headers() {
      return [
        {
          source: '/(.*)',
          headers: securityHeaders,
        },
      ]
    },
  })
)
