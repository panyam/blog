package components

import (
	"net/http"
)

type Footer struct {
	BaseView
	ThemeName   string
	IsDarkTheme bool
}

func (v *Footer) InitContext(c *Context, pv View) {
	v.BaseView.InitContext(c, pv)
}

func (v *Footer) ValidateRequest(w http.ResponseWriter, r *http.Request) (err error) {
	v.BaseView.ValidateRequest(w, r)
	return
}

func (h *Footer) RenderResponse() error {
	return h.Ctx.Template.ExecuteTemplate(h.Ctx.Writer, "Footer", h)
}
