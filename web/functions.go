package web

import (
	"html/template"
	"log"
	"sort"

	gfn "github.com/panyam/goutils/fn"
	s3 "github.com/panyam/s3gen/core"
)

func (web *BlogWeb) TemplateFunctions() template.FuncMap {
	return template.FuncMap{
		"PagesByDate":   web.GetPagesByDate,
		"AllTags":       web.GetAllTags,
		"KeysForTagMap": web.KeysForTagMap,
		"AllRes": func() []*s3.Resource {
			resources := site.ListResources(
				func(res *s3.Resource) bool {
					return !res.IsParametric
				},
				// sort by reverse date order
				/*sort=*/
				nil, -1, -1)
			sort.Slice(resources, func(idx1, idx2 int) bool {
				res1 := resources[idx1]
				res2 := resources[idx2]
				return res1.CreatedAt.Sub(res2.CreatedAt) > 0
			})
			return resources
		},
	}
}

func (web *BlogWeb) GetPagesByDate(desc bool, offset, count int) (out []*s3.Resource) {
	return site.ListResources(
		func(res *s3.Resource) bool {
			return !res.IsParametric && (res.NeedsIndex || res.IsIndex)
			// && (strings.HasSuffix(res.FullPath, ".md") || strings.HasSuffix(res.FullPath, ".mdx"))
		},
		func(res1, res2 *s3.Resource) bool {
			d1 := res1.DestPage
			d2 := res2.DestPage
			if d1 == nil || d2 == nil {
				log.Println("D1: ", res1.FullPath)
				log.Println("D2: ", res2.FullPath)
				return false
			}
			sub := res1.DestPage.CreatedAt.Sub(res2.DestPage.CreatedAt)
			if desc {
				return sub > 0
			} else {
				return sub < 0
			}
		}, offset, count+1)
}

func (web *BlogWeb) KeysForTagMap(tagmap map[string]int, orderby string) []string {
	out := gfn.MapKeys(tagmap)
	sort.Slice(out, func(i1, i2 int) bool {
		c1 := tagmap[out[i1]]
		c2 := tagmap[out[i2]]
		if c1 == c2 {
			return out[i1] < out[i2]
		}
		return c1 > c2
	})
	return out
}

func (web *BlogWeb) GetAllTags(resources []*s3.Resource) (tagCount map[string]int) {
	tagCount = make(map[string]int)
	for _, res := range resources {
		if res.FrontMatter().Data != nil {
			if t, ok := res.FrontMatter().Data["tags"]; ok && t != nil {
				if tags, ok := t.([]any); ok && tags != nil {
					for _, tag := range tags {
						tagCount[tag.(string)] += 1
					}
				}
			}
		}
	}
	return
}
