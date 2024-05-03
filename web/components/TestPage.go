package components

type TestPage struct {
	BaseView
}

func (v *TestPage) RenderResponse() error {
	return v.Ctx.Template.ExecuteTemplate(v.Ctx.Writer, "TestPage.html", v)
}
