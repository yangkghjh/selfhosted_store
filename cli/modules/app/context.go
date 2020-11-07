package app

import (
	"strings"

	"github.com/yankghjh/selfhosted_store/cli/project"
)

// Context of the app load pipe
type Context struct {
	*project.Operator
	Path string
	Name string
}

// NewContext for the pipeline of one app
func NewContext(o *project.Operator, name, path string) *Context {
	return &Context{
		Operator: o,
		Name:     name,
		Path:     path,
	}
}

// GetPath of source file in app path
func (c *Context) GetPath(paths ...string) string {
	return c.Path + "/" + strings.Join(paths, "/")
}
