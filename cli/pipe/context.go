package pipe

import "strings"

// Context of the pipe
type Context struct {
	Path string
	data map[string]interface{}
}

// NewContext for the pipeline of one app
func NewContext(path string) *Context {
	return &Context{
		Path: path,
		data: map[string]interface{}{},
	}
}

// Get value store in context
func (c *Context) Get(key string) interface{} {
	value, ok := c.data[key]
	if !ok {
		return nil
	}

	return value
}

// Set vaule to context
func (c *Context) Set(key string, value interface{}) {
	c.data[key] = value
}

// GetPath of source file in app path
func (c *Context) GetPath(paths ...string) string {
	return c.Path + "/" + strings.Join(paths, "/")
}
