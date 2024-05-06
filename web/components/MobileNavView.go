package components

type MobileNav struct {
	BaseView
	HeaderNavLinks []HeaderNavLink
	ShowNav        bool
}

func (h *MobileNav) RenderResponse() error {
	return h.Ctx.Render(h, "MobileNav")
}
