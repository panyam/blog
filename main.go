package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/felixge/httpsnoop"
	"github.com/gorilla/mux"
	gfn "github.com/panyam/goutils/fn"
	s3 "github.com/panyam/s3gen"
)

var (
	addr = flag.String("addr", DefaultGatewayAddress(), "Address where the http grpc gateway endpoint is running")
)

var site = s3.Site{
	ContentRoot: "./content",
	OutputDir:   "./output",
	HtmlTemplates: []string{
		"templates/*.html",
	},
	StaticFolders: []string{
		"/static/", "static",
	},
}

func main() {
	flag.Parse()
	router := mux.NewRouter()

	if os.Getenv("APP_ENV") != "production" {
		site.CommonFuncMap = TemplateFunctions()
		site.NewViewFunc = NewView
		site.Watch()
	}

	// Attach our site to be at /`PathPrefix`
	// The site will also take care of serving static files from /`PathPrefix`/static paths
	router.PathPrefix(site.PathPrefix).Handler(http.StripPrefix(site.PathPrefix, &site))

	srv := &http.Server{
		Handler: withLogger(router),
		Addr:    *addr,
		// Good practice: enforce timeouts for servers you create!
		// WriteTimeout: 15 * time.Second,
		// ReadTimeout:  15 * time.Second,
	}
	log.Printf("Serving Gateway endpoint on %s:", *addr)
	log.Fatal(srv.ListenAndServe())
}

func DefaultGatewayAddress() string {
	gateway_addr := os.Getenv("BLOG_WEB_PORT")
	if gateway_addr != "" {
		return gateway_addr
	}
	return ":8080"
}

func withLogger(handler http.Handler) http.Handler {
	// the create a handler
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		// pass the handler to httpsnoop to get http status and latency
		m := httpsnoop.CaptureMetrics(handler, writer, request)
		// printing exracted data
		log.Printf("http[%d]-- %s -- %s\n", m.Code, m.Duration, request.URL.Path)
	})
}

// /////////// Page View related items
func NewView(name string) (out s3.View) {
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

// //////////// Functions for our site
func TemplateFunctions() template.FuncMap {
	return template.FuncMap{
		"LeafPages":     LeafPages,
		"PagesByDate":   GetPagesByDate,
		"AllTags":       GetAllTags,
		"KeysForTagMap": KeysForTagMap,
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

func LeafPages(hideDrafts bool, orderby string, offset, count any) (out []*s3.Resource) {
	var sortFunc s3.ResourceSortFunc = nil
	if orderby != "" {
		desc := orderby[0] == '-'
		if desc {
			orderby = orderby[1:]
		}
		sortFunc = func(res1, res2 *s3.Resource) bool {
			d1 := res1.DestPage
			d2 := res2.DestPage
			if d1 == nil || d2 == nil {
				log.Println("D1: ", res1.FullPath)
				log.Println("D2: ", res2.FullPath)
				return false
			}
			sub := 0
			if orderby == "date" {
				sub = int(res1.DestPage.CreatedAt.Sub(res2.DestPage.CreatedAt))
			} else if orderby == "title" {
				sub = strings.Compare(d1.Title, d2.Title)
			}
			if desc {
				return sub > 0
			} else {
				return sub < 0
			}
		}
	}
	return site.ListResources(
		func(res *s3.Resource) bool {
			// Leaf pages only - not indexes
			if res.IsParametric || !res.NeedsIndex || res.IsIndex {
				return false
			}

			if hideDrafts {
				draft := res.FrontMatter().Data["draft"]
				if draft == true {
					return false
				}
			}
			return true
			// && (strings.HasSuffix(res.FullPath, ".md") || strings.HasSuffix(res.FullPath, ".mdx"))
		},
		sortFunc,
		s3.ToInt(offset), s3.ToInt(count))
}

func GetPagesByDate(hideDrafts bool, desc bool, offset, count any) (out []*s3.Resource) {
	return site.ListResources(
		func(res *s3.Resource) bool {
			if res.IsParametric || !(res.NeedsIndex || res.IsIndex) {
				return false
			}

			if hideDrafts {
				draft := res.FrontMatter().Data["draft"]
				if draft == true {
					return false
				}
			}
			return true
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
		},
		s3.ToInt(offset), s3.ToInt(count))
}

func KeysForTagMap(tagmap map[string]int, orderby string) []string {
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

func GetAllTags(resources []*s3.Resource) (tagCount map[string]int) {
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
