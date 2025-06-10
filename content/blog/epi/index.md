---
title: Elements of Programming Interviews
date: 2014-03-26T00:00:00AM
tags: ['elements of programming interviews', 'algorithms', 'programming', 'haskell']
draft: false
images: []
authors: ['Sri Panyam']
template: BasePageMD.html
---

This is an index to the problems from THE book being solved in haskell (and counting). Some of the trivial ones will only be attempted on demand!

## Chapter 5 – Primitive Types

{{ $siteMetadata := (json "SiteMetadata.json" "") }}
  
{{ $posts := ( PagesByDate $siteMetadata.HideDrafts false 0 -1 ) }}

{{ range $posts }}{{ if HasPrefix .Base.Title  "EPI 5." }}<p><a href="{{.Base.Link}}"><span>{{.Base.Title}}</span></a></p>{{ end }}{{ end }}

## Chapter 6 – Strings and Arrays

{{ range $posts }}{{ if HasPrefix .Base.Title  "EPI 6." }}<p><a href="{{.Base.Link}}"><span>{{.Base.Title}}</span></a></p>{{ end }}{{ end }}
