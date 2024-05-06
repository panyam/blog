package components

type Header struct {
	BaseView
	HeaderNavLinks  []HeaderNavLink
	ThemeSwitchView ThemeSwitch
	MobileNavView   MobileNav
}

func (v *Header) InitContext(c *Context, pv View) {
	v.BaseView.AddChildViews(&v.ThemeSwitchView, &v.MobileNavView)
	v.BaseView.InitContext(c, pv)
	v.MobileNavView.HeaderNavLinks = v.HeaderNavLinks
}

func (v *Header) RenderResponse() error {
	return v.Ctx.Render(v, "Header")
}
