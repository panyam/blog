package components

type Header struct {
	BaseView
	HeaderNavLinks  []HeaderNavLink
	ThemeSwitchView ThemeSwitch
	MobileNavView   MobileNav
	SiteMetadata    SiteMetadata
}

func (v *Header) InitContext(c *Context, pv View) {
	v.BaseView.AddChildViews(&v.ThemeSwitchView, &v.MobileNavView)
	v.BaseView.InitContext(c, pv)

	v.MobileNavView.HeaderNavLinks = v.HeaderNavLinks
}

func (h *Header) RenderResponse() error {
	return h.Ctx.Template.ExecuteTemplate(h.Ctx.Writer, "Header", h)
}
