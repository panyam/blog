package web

import (
	s3 "github.com/panyam/s3gen/core"
)

func (web *BlogWeb) NewView(name string) (out s3.View) {
	if name == "BasePage" || name == "" {
		out = &BasePage{BaseView: s3.BaseView{Template: "BasePage.html"}}
	}
	if name == "PostPage" || name == "" {
		out = &PostPage{}
	}
	if name == "AuthorPage" || name == "" {
		out = &AuthorPage{}
	}
	return
}

type BasePage struct {
	s3.BaseView
	// PageSEO    SEO
	HeaderView Header
	BodyView   s3.View
	FooterView Footer
}

func (v *BasePage) InitView(s *s3.Site, parentView s3.View) {
	if v.Self == nil {
		v.Self = v
	}
	v.BaseView.AddChildViews( /*&v.PageSEO, */ &v.HeaderView, v.BodyView, &v.FooterView)
	v.BaseView.InitView(s, parentView)
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

/*
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
