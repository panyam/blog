package components

import (
	"net/http"
)

type HomePage struct {
	BaseView
}

func (v *HomePage) InitContext(c *Context, pv View) {
	v.BaseView.InitContext(c, pv)
}

func (v *HomePage) ValidateRequest(w http.ResponseWriter, r *http.Request) (err error) {
	v.BaseView.ValidateRequest(w, r)
	return
}

func (h *HomePage) RenderResponse() error {
	return h.Ctx.Template.ExecuteTemplate(h.Ctx.Writer, "index.html", h)
}
