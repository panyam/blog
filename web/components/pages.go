package components

type BasePage struct {
	BaseView
	PageSEO        PageSEO
	HeaderView     Header
	BodyView       View
	FooterView     Footer
	HeaderNavLinks []HeaderNavLink
}

func (v *BasePage) InitContext(c *Context, parentView View) {
	v.BaseView.AddChildViews(&v.PageSEO, &v.HeaderView, v.BodyView, &v.FooterView)
	v.HeaderView.HeaderNavLinks = v.HeaderNavLinks
	v.BaseView.InitContext(c, parentView)
}

func (v *BasePage) RenderResponse() (err error) {
	return v.Ctx.Render(v, "BasePage.html")
}

type HomePage struct {
	BaseView
	Posts []Post
}

func (v *HomePage) RenderResponse() (err error) {
	return v.Ctx.Render(v, "HomePage")
}

type AboutPage struct {
	BaseView
}
