package web

import (
	"html/template"
	"net/http"
	"sort"
	"strings"

	"github.com/gorilla/mux"
	gfn "github.com/panyam/goutils/fn"
	s3 "github.com/panyam/s3gen/core"
)

var site = s3.Site{
	ContentRoot: "./data",
	OutputDir:   "/Users/sri/personal/golang/blog/published",
	PathPrefix:  "/published",
	LazyLoad:    true,
	HtmlTemplates: []string{
		"templates/*.html",
	},
	StaticFolders: []string{
		"/static/", "public/static",
	},
}

// This should be mirroring how we are setting up our app.yaml
func (web *BlogWeb) setupPages(router *mux.Router) {
	site.NewViewFunc = web.NewView
	site.CommonFuncMap = template.FuncMap{
		"AllPages": func() []*s3.Page {
			resources := site.ListResources(
				func(res *s3.Resource) bool {
					return strings.HasSuffix(res.FullPath, ".md") ||
						strings.HasSuffix(res.FullPath, ".mdx")
				},
				// sort by reverse date order
				/*sort=*/
				nil, -1, -1)
			pages := gfn.Map(resources, func(r *s3.Resource) *s3.Page {
				p, _ := site.GetPage(r)
				return p
			})
			sort.Slice(pages, func(idx1, idx2 int) bool {
				page1 := pages[idx1]
				page2 := pages[idx2]
				return page1.CreatedAt.Sub(page2.CreatedAt) > 0
			})
			return pages
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
	if name == "BlogsPage" || name == "" {
		out = &BlogsPage{}
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
	if v.MaxPostsToDisplay == 0 {
		v.MaxPostsToDisplay = 20
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

type BlogsPage struct {
	BasePage
}
