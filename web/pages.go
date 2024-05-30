package web

import (
	"html/template"
	"log"
	"net/http"
	"sort"

	"github.com/gorilla/mux"
	gfn "github.com/panyam/goutils/fn"
	s3 "github.com/panyam/s3gen/core"
)

var site = s3.Site{
	ContentRoot: "./data",
	OutputDir:   "/Users/sri/personal/golang/blog/build",
	PathPrefix:  "/published",
	LazyLoad:    true,
	HtmlTemplates: []string{
		"templates/*.html",
	},
	StaticFolders: []string{
		"/static/", "public/static",
	},
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

// This should be mirroring how we are setting up our app.yaml
func (web *BlogWeb) setupPages(router *mux.Router) {
	site.NewViewFunc = web.NewView
	site.CommonFuncMap = template.FuncMap{
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

	site.Init().Load().StartWatching()

	// Here we want to point just to the root of our blog and let it get served
	// For now we will serve via a router but then take the same router to
	// publish them for static serving too
	router.PathPrefix(site.PathPrefix).Handler(http.StripPrefix(site.PathPrefix, &site))
}

func (web *BlogWeb) NewView(name string) (out s3.View) {
	if name == "BasePage" || name == "" {
		out = &BasePage{BaseView: s3.BaseView{Template: "BasePage.html"}}
	}
	if name == "PostPage" || name == "" {
		out = &PostPage{}
	}
	if name == "HomePage" || name == "" {
		out = &HomePage{}
	}
	if name == "AuthorPage" || name == "" {
		out = &AuthorPage{}
	}
	return
}

type BasePage struct {
	s3.BaseView
	PageSEO    SEO
	HeaderView Header
	BodyView   s3.View
	FooterView Footer
}

func (v *BasePage) InitView(s *s3.Site, parentView s3.View) {
	if v.Self == nil {
		v.Self = v
	}
	v.BaseView.AddChildViews(&v.PageSEO, &v.HeaderView, v.BodyView, &v.FooterView)
	v.BaseView.InitView(s, parentView)
}

type HomePageBodyView struct {
	s3.BaseView
	ContentView       s3.View
	MaxPostsToDisplay int
}

func (v *HomePageBodyView) InitView(s *s3.Site, parentView s3.View) {
	if v.Self == nil {
		v.Self = v
	}
	v.BaseView.AddChildViews(v.ContentView)
	v.BaseView.InitView(s, parentView)
}

type HomePage struct {
	BasePage
	BodyView HomePageBodyView
}

func (v *HomePage) InitView(s *s3.Site, parentView s3.View) {
	if v.Self == nil {
		v.Self = v
	}
	if v.Template == "" {
		v.Template = "BasePage.html"
	}
	v.BasePage.AddChildViews(&v.BodyView)
	v.BasePage.InitView(s, parentView)
}

type PostListView struct {
	s3.BaseView
}

type AuthorPage struct {
	BasePage
	BodyView AuthorLayout
}

type AuthorLayout struct {
	s3.BaseView
	ContentView s3.View
}

func (v *AuthorLayout) InitView(s *s3.Site, parentView s3.View) {
	if v.Self == nil {
		v.Self = v
	}
	v.BaseView.AddChildViews(v.ContentView)
	v.BaseView.InitView(s, parentView)
}

func (v *AuthorPage) InitView(s *s3.Site, parentView s3.View) {
	if v.Self == nil {
		v.Self = v
	}
	if v.Template == "" {
		v.Template = "BasePage.html"
	}
	v.BasePage.AddChildViews(&v.BodyView)
	v.BasePage.InitView(s, parentView)
}

type PostPage struct {
	BasePage
	BodyView PostSimple
}

func (v *PostPage) InitView(s *s3.Site, parentView s3.View) {
	if v.Self == nil {
		v.Self = v
	}
	if v.Template == "" {
		v.Template = "BasePage.html"
	}
	v.BasePage.AddChildViews(&v.BodyView)
	v.BasePage.InitView(s, parentView)
}

type Header struct {
	s3.BaseView
	ThemeSwitchView ThemeSwitch
	MobileNavView   MobileNav
}

func (v *Header) InitView(s *s3.Site, pv s3.View) {
	if v.Self == nil {
		v.Self = v
	}
	v.BaseView.AddChildViews(&v.ThemeSwitchView, &v.MobileNavView)
	v.BaseView.InitView(s, pv)
}

type MobileNav struct {
	s3.BaseView
	ShowNav bool
}

func (v *MobileNav) InitView(s *s3.Site, pv s3.View) {
	if v.Self == nil {
		v.Self = v
	}
	v.BaseView.InitView(s, pv)
}

type ThemeSwitch struct {
	s3.BaseView
	ThemeName   string
	IsDarkTheme bool
}

func (v *ThemeSwitch) InitView(s *s3.Site, pv s3.View) {
	if v.Self == nil {
		v.Self = v
	}
	v.BaseView.InitView(s, pv)
}

type Footer struct {
	s3.BaseView
	ThemeName   string
	IsDarkTheme bool
}

func (v *Footer) InitView(s *s3.Site, pv s3.View) {
	if v.Self == nil {
		v.Self = v
	}
	v.BaseView.InitView(s, pv)
}

type SEO struct {
	s3.BaseView
	OgType   string
	OgImages []string
	TwImage  string
}

func (v *SEO) InitView(s *s3.Site, pv s3.View) {
	if v.Self == nil {
		v.Self = v
	}
	v.BaseView.InitView(s, pv)
}

/*
func (v *s3.PageSEO) InitView(s *s3.Site, pv s3.View) {
	smd := v.Site.SiteMetadata
	SiteUrl := GetProp(v.Site.SiteMetadata, "SiteUrl").(string)
	SocialBanner:= GetProp(v.Site.SiteMetadata, "SocialBanner").(string)
	v.CommonSEO.OgImages = []string{
		SiteUrl + SocialBanner,
	}
	v.CommonSEO.TwImage = SiteUrl + SocialBanner
	v.CommonSEO.InitView(s, pv)
}
*/

type PostSimple struct {
	s3.BaseView
	ContentView s3.View
}

func (v *PostSimple) InitView(s *s3.Site, parentView s3.View) {
	if v.Self == nil {
		v.Self = v
	}
	v.BaseView.AddChildViews(v.ContentView)
	v.BaseView.InitView(s, parentView)
}
