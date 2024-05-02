package components

import (
	"html/template"
	"io"
	"net/http"
)

type Context struct {
	Writer   io.Writer
	Template *template.Template
}

type View interface {
	InitContext(*Context, View)
	ValidateRequest(w http.ResponseWriter, r *http.Request) error
	RenderResponse() error
}

type BaseView struct {
	Parent View
	Ctx    *Context
}

func (v *BaseView) InitContext(c *Context, parent View) {
	v.Ctx = c
	v.Parent = parent
}

func (v *BaseView) ValidateRequest(w http.ResponseWriter, r *http.Request) (err error) {
	return
}

func (v *BaseView) RenderResponse() (err error) {
	_, err = v.Ctx.Writer.Write([]byte("TemplateName not provided"))
	return
}
