package components

type CommonSEO struct {
	BaseView
	Title        string
	Description  string
	OgType       string
	OgImages     []string
	TwImage      string
	CanonicalUrl string
}

type PageSEO struct {
	CommonSEO
}

func (v *PageSEO) InitContext(c *Context, pv View) {
	v.CommonSEO.OgImages = []string{
		c.SiteMetadata.SiteUrl + c.SiteMetadata.SocialBanner,
	}
	v.CommonSEO.TwImage = c.SiteMetadata.SiteUrl + c.SiteMetadata.SocialBanner
	v.CommonSEO.InitContext(c, pv)
}

func (v *PageSEO) RenderResponse() (err error) {
	return v.Ctx.Render(v, "PageSEO")
}
