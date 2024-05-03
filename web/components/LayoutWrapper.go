package components

type LayoutWrapper struct {
	BaseView
	SectionContainer SectionContainer
}

func (v *LayoutWrapper) RenderResponse() error {
	return v.Ctx.Template.ExecuteTemplate(v.Ctx.Writer, "LayoutWrapper.html", v)
}
