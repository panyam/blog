package web

import (
	"log"
	"net/http"
	"reflect"

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

func (web *BlogWeb) NewPageView(name string) (out s3.PageView) {
	if name == "BasePage" || name == "" {
		out = &BasePage{BaseView: s3.BaseView{Template: "BasePage.html"}}
	}
	if name == "PostPage" || name == "" {
		out = &PostPage{}
	}
	return
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

	/*
		site.HandlePage("/:slug", func(w http.ResponseWriter, r *http.Request) {
			// need a function to go slug -> s3.View
			view := components.BasePage{
				HeaderNavLinks: headerNavLinks,
				BodyView:       &components.HomePage{},
			}
		})

		router.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
			aboutpage := &components.BasePage{
				HeaderNavLinks: headerNavLinks,
				BodyView:       &components.AboutPage{},
			}
			web.RenderView(aboutpage, w, r, "AboutPage")
		}).Methods("GET")

		// router.SetFuncMap(funcmap)
		// router.LoadHTMLGlob("./web/templates/*.*")
		router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			homepage := &components.BasePage{
				s3.BaseView:       components.BaseView{},
				HeaderNavLinks: headerNavLinks,
				BodyView:       &components.HomePage{},
			}
			web.RenderView(homepage, w, r, "")
		}).Methods("GET")
	*/
}

/*
func (web *BlogWeb) RenderView(v components.View, w http.ResponseWriter, r *http.Request, templateName string) {
	ctx := &components.Context{
		Writer:       w,
		Template:     web.Template,
		SiteMetadata: siteMetadata,
	}
	// w.WriteHeader(http.StatusOK)
	v.InitContext(ctx, nil)
	err := v.ValidateRequest(w, r)
	if err == nil {
		err = ctx.Render(v, templateName)
	}
	if err != nil {
		slog.Error("Render Error: ", "err", err)
		http.Error(w, err.Error(), 500)
		// c.Abort()
	}
}

func CustomFuncMap() template.FuncMap {
	return template.FuncMap{
		"RenderView": func(v components.View) template.URL {
			log.Println("V: ", v)
			v.RenderResponse()
			return template.URL(" ")
		},
	}
}
*/

type BasePage struct {
	s3.BaseView
	Page       *s3.Page
	PageSEO    SEO
	HeaderView Header
	BodyView   s3.View
	FooterView Footer
}

func (v *BasePage) InitContext(s *s3.Site, parentView s3.View) {
	if v.Self == nil {
		v.Self = v
	}
	v.BaseView.AddChildViews(&v.PageSEO, &v.HeaderView, v.BodyView, &v.FooterView)
	v.BaseView.InitContext(s, parentView)
}

func (v *BasePage) GetPage() *s3.Page {
	return v.Page
}

func (v *BasePage) SetPage(p *s3.Page) {
	v.Page = p
}

type PostPage struct {
	BasePage
	BodyView PostSimple
}

func (v *PostPage) InitContext(s *s3.Site, parentView s3.View) {
	if v.Self == nil {
		v.Self = v
	}
	if v.Template == "" {
		v.Template = "BasePage.html"
	}
	v.BasePage.AddChildViews(&v.BodyView)
	v.BasePage.InitContext(s, parentView)
	log.Println("PP After: ", reflect.TypeOf(v), reflect.TypeOf(v.BodyView.Parent))
}

type Header struct {
	s3.BaseView
	ThemeSwitchView ThemeSwitch
	MobileNavView   MobileNav
}

func (v *Header) InitContext(s *s3.Site, pv s3.View) {
	if v.Self == nil {
		v.Self = v
	}
	v.BaseView.AddChildViews(&v.ThemeSwitchView, &v.MobileNavView)
	v.BaseView.InitContext(s, pv)
}

type MobileNav struct {
	s3.BaseView
	ShowNav bool
}

func (v *MobileNav) InitContext(s *s3.Site, pv s3.View) {
	if v.Self == nil {
		v.Self = v
	}
	v.BaseView.InitContext(s, pv)
}

type ThemeSwitch struct {
	s3.BaseView
	ThemeName   string
	IsDarkTheme bool
}

func (v *ThemeSwitch) InitContext(s *s3.Site, pv s3.View) {
	if v.Self == nil {
		v.Self = v
	}
	v.BaseView.InitContext(s, pv)
}

type Footer struct {
	s3.BaseView
	ThemeName   string
	IsDarkTheme bool
}

func (v *Footer) InitContext(s *s3.Site, pv s3.View) {
	if v.Self == nil {
		v.Self = v
	}
	v.BaseView.InitContext(s, pv)
}

type SEO struct {
	s3.BaseView
	OgType   string
	OgImages []string
	TwImage  string
}

func (v *SEO) InitContext(s *s3.Site, pv s3.View) {
	if v.Self == nil {
		v.Self = v
	}
	v.BaseView.InitContext(s, pv)
}

/*
func (v *s3.PageSEO) InitContext(s *s3.Site, pv s3.View) {
	smd := v.Site.SiteMetadata
	SiteUrl := GetProp(v.Site.SiteMetadata, "SiteUrl").(string)
	SocialBanner:= GetProp(v.Site.SiteMetadata, "SocialBanner").(string)
	v.CommonSEO.OgImages = []string{
		SiteUrl + SocialBanner,
	}
	v.CommonSEO.TwImage = SiteUrl + SocialBanner
	v.CommonSEO.InitContext(s, pv)
}
*/

type PostSimple struct {
	s3.BaseView
	ContentView s3.View
	Post        *s3.Page
	PrevPost    *s3.Page
	NextPost    *s3.Page
}

func (v *PostSimple) InitContext(s *s3.Site, parentView s3.View) {
	if v.Self == nil {
		v.Self = v
	}
	v.BaseView.AddChildViews(v.ContentView)
	v.BaseView.InitContext(s, parentView)
}
