package components

type SectionContainer struct {
	BaseView
}

func (v *SectionContainer) RenderResponse() error {
	return v.Ctx.Template.ExecuteTemplate(v.Ctx.Writer, "SectionContainer.html", v)
}
