package components

type ThemeSwitch struct {
	BaseView
	ThemeName   string
	IsDarkTheme bool
}

func (h *ThemeSwitch) RenderResponse() error {
	return h.Ctx.Template.ExecuteTemplate(h.Ctx.Writer, "ThemeSwitch", h)
}
