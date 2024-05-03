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
	Parent       View
	Ctx          *Context
	TemplateName string
	Children     []View
}

func (v *BaseView) InitContext(c *Context, parent View) {
	v.Ctx = c
	v.Parent = parent
	for _, child := range v.Children {
		child.InitContext(c, v)
	}
}

func (v *BaseView) ValidateRequest(w http.ResponseWriter, r *http.Request) (err error) {
	for _, child := range v.Children {
		err = child.ValidateRequest(w, r)
		if err != nil {
			return
		}
	}
	return
}

func (v *BaseView) RenderResponse() (err error) {
	if v.TemplateName == "" {
		_, err = v.Ctx.Writer.Write([]byte("TemplateName not provided"))
	} else {
		return v.Ctx.Template.ExecuteTemplate(v.Ctx.Writer, v.TemplateName, v)
	}
	return
}

func (v *BaseView) AddChildViews(views ...View) {
	for _, child := range views {
		v.Children = append(v.Children, child)
	}
}
