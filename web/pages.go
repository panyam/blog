package web

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/panyam/blog/web/components"
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

var siteMetadata = &components.SiteMetadata{
	HeaderTitle: "Buildmage",
	Title:       "Buildmage",
	Description: "My personal blog highlighting adventures in building random things on the side",
	Author:      "Sriram Panyam",
}

// This should be mirroring how we are setting up our app.yaml
func (web *BlogWeb) setupPages(router *mux.Router) {
	site.NewPageViewFunc = web.NewPageView
	site.SiteMetadata = siteMetadata
	site.Init().Load().StartWatching()

	// Here we want to point just to the root of our blog and let it get served
	// For now we will serve via a router but then take the same router to
	// publish them for static serving too
	router.PathPrefix(site.PathPrefix).Handler(http.StripPrefix(site.PathPrefix, &site))
}

func (web *BlogWeb) NewPageView(name string) (out s3.PageView) {
	if name == "BasePage" || name == "" {
		out = &BasePage{BasePageView: s3.BasePageView{BaseView: s3.BaseView{Template: "BasePage.html"}}}
	}
	if name == "PostPage" || name == "" {
		out = &PostPage{}
	}
	return
}

type BasePage struct {
	s3.BasePageView
	PageSEO    SEO
	HeaderView Header
	BodyView   s3.View
	FooterView Footer
}

func (v *BasePage) InitView(s *s3.Site, parentView s3.View) {
	if v.Self == nil {
		v.Self = v
	}
	v.BasePageView.AddChildViews(&v.PageSEO, &v.HeaderView, v.BodyView, &v.FooterView)
	v.BasePageView.InitView(s, parentView)
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
	Post        *s3.Page
	PrevPost    *s3.Page
	NextPost    *s3.Page
}

func (v *PostSimple) InitView(s *s3.Site, parentView s3.View) {
	if v.Self == nil {
		v.Self = v
	}
	v.BaseView.AddChildViews(v.ContentView)
	v.BaseView.InitView(s, parentView)
}
