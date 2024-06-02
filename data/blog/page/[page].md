---
title: BuildMage Blog
date: 2024-03-26T00:00:00AM
draft: false
images: []
authors: ['Sri Panyam']
page: BasePage
location: "BodyView"
---

{{ . }}

{{ $siteMetadata := (json "SiteMetadata.json" "") }}
{{ $pageSize := Int $siteMetadata.MaxPostsPerPage }}

{{ if eq .Res.CurrentParamName "" }}
  {{ $posts := LeafPages $siteMetadata.ShowDrafts "" 0 -1 }}
  {{ $numPages := len $posts }}
  {{ range (IntList 1 $numPages 1) }}
    {{ $.Res.AddParam (String .) }}
  {{ end }}
{{ else }}
  {{/* Do what is needed to show data for this page - list layout etc */}}
  Current Page {{.Res.CurrentParamName}}
{{ end }}

