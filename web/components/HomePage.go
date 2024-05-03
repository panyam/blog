package components

type AboutPage struct {
	BaseView
	BodyView       View
	HeaderView     Header
	FooterView     Footer
	SiteMetadata   SiteMetadata
	HeaderNavLinks []HeaderNavLink
}

func (v *AboutPage) InitContext(c *Context, pv View) {
	v.BaseView.AddChildViews(&v.HeaderView, v.BodyView, &v.FooterView)
	v.HeaderView.HeaderNavLinks = v.HeaderNavLinks
	v.HeaderView.SiteMetadata = v.SiteMetadata
	v.BaseView.InitContext(c, pv)
}

func (h *AboutPage) RenderResponse() error {
	return h.Ctx.Template.ExecuteTemplate(h.Ctx.Writer, "AboutPage.html", h)
}
