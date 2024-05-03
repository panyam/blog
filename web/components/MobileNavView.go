package components

type MobileNav struct {
	BaseView
	HeaderNavLinks []HeaderNavLink
	ShowNav        bool
}

func (h *MobileNav) RenderResponse() error {
	return h.Ctx.Template.ExecuteTemplate(h.Ctx.Writer, "MobileNav", h)
}
