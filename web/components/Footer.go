package components

type Footer struct {
	BaseView
	ThemeName   string
	IsDarkTheme bool
}

func (v *Footer) RenderResponse() error {
	return v.Ctx.Render(v, "Footer")
}
