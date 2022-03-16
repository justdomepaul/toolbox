package restful

import (
	"github.com/cockroachdb/errors"
	"github.com/gin-gonic/gin/render"
	"html/template"
)

// NewRender returns a fresh instance of Render
func NewRender() *Render {
	return &Render{
		templates: make(map[string]*template.Template),
		Files:     make(map[string][]string),
	}
}

// Render implements gin's HTMLRender and provides some sugar on top of it
type Render struct {
	templates    map[string]*template.Template
	Files        map[string][]string
	TemplatesDir string
	Layout       string
	Ext          string
	Debug        bool
}

// Add assigns the name to the template
func (r *Render) Add(name string, tmpl *template.Template) error {
	if tmpl == nil {
		return errors.New("template can not be nil")
	}
	if len(name) == 0 {
		return errors.New("template name cannot be empty")
	}
	r.templates[name] = tmpl
	return nil
}

// Instance implements gin's HTML render interface
func (r *Render) Instance(name string, data interface{}) render.Render {
	tpl := r.templates[name]

	return &render.HTML{
		Template: tpl,
		Data:     data,
	}
}
