package pipe

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/spf13/viper"
)

/**
0. Call init functions
1. Load apps
2. For the each of apps:
	1. Execute source loaders
	2. Execute source handlers
3. Call finish fucntions
**/

var initFuncs map[string]InitFunc
var finishFuncs map[string]FinishFunc
var sourceLoaders map[string]SourceLoader
var sourceHandlers map[string]SourceHandler

func init() {
	initFuncs = map[string]InitFunc{}
	finishFuncs = map[string]FinishFunc{}
	sourceLoaders = map[string]SourceLoader{}
	sourceHandlers = map[string]SourceHandler{}
}

// Pipe is the pipeline for the templetes generator
type Pipe struct {
	Opition
	Count int
	Apps  []*Context

	data map[string]interface{}
}

// Opition for the pipe
type Opition struct {
	SourcePath     string
	DistPath       string
	Sources        []string
	Handlers       []string
	Config         *viper.Viper
	SkipSourceFile bool
}

// InitFunc called before loading apps
type InitFunc func(pipe *Pipe) error

// FinishFunc called after handling apps
type FinishFunc func(pipe *Pipe) error

// SourceLoader load source and save to context
type SourceLoader func(pipe *Pipe, ctx *Context) error

// SourceHandler process the loaded source
type SourceHandler func(pipe *Pipe, ctx *Context) error

// NewPipe with the opition
func NewPipe(opt Opition) (*Pipe, error) {
	p := &Pipe{
		Opition: opt,
		Apps:    []*Context{},

		data: map[string]interface{}{},
	}

	for _, source := range p.Sources {
		if _, ok := sourceLoaders[source]; !ok {
			return nil, fmt.Errorf("source %s not exist", source)
		}
	}

	for _, handler := range p.Handlers {
		if _, ok := sourceHandlers[handler]; !ok {
			return nil, fmt.Errorf("handler %s not exist", handler)
		}
	}

	return p, nil
}

// Get value stored in pipe
func (p *Pipe) Get(key string) interface{} {
	value, ok := p.data[key]
	if !ok {
		return nil
	}

	return value
}

// Set vaule to pipe
func (p *Pipe) Set(key string, value interface{}) {
	p.data[key] = value
}

// GetDistPath get path of dist
func (p *Pipe) GetDistPath(paths ...string) string {
	return p.DistPath + "/" + strings.Join(paths, "/")
}

// Run the pipe
func (p *Pipe) Run() error {
	for name, initFunc := range initFuncs {
		if err := initFunc(p); err != nil {
			return fmt.Errorf("init %s error: %s", name, err.Error())
		}
	}

	if !p.SkipSourceFile {
		files, err := ioutil.ReadDir(p.SourcePath)
		if err != nil {
			return fmt.Errorf("read source path error: %s", err.Error())
		}

		for _, f := range files {
			if f.IsDir() {
				p.Apps = append(p.Apps, NewContext(f.Name(), p.SourcePath+"/"+f.Name()))
			}
		}
	}

	p.Count = len(p.Apps)

	for _, ctx := range p.Apps {
		// load source
		for _, source := range p.Sources {
			f := sourceLoaders[source]

			err := f(p, ctx)
			if err != nil {
				return fmt.Errorf("load source [%s] of path [%s] error: %s", source, ctx.Path, err.Error())
			}
		}

		// handle
		for _, handler := range p.Handlers {
			f := sourceHandlers[handler]

			err := f(p, ctx)
			if err != nil {
				return fmt.Errorf("run handler [%s] of path [%s] error: %s", handler, ctx.Path, err.Error())
			}
		}
	}

	for name, finishFunc := range finishFuncs {
		if err := finishFunc(p); err != nil {
			return fmt.Errorf("finish %s error: %s", name, err.Error())
		}
	}

	return nil
}

// RegisterInitFunc register init function
func RegisterInitFunc(name string, f InitFunc) {
	initFuncs[name] = f
}

// RegisterFinishFunc register finish function
func RegisterFinishFunc(name string, f FinishFunc) {
	finishFuncs[name] = f
}

// RegisterSourceLoader register source loader
func RegisterSourceLoader(name string, f SourceLoader) {
	sourceLoaders[name] = f
}

// RegisterSourceHandler register source handler
func RegisterSourceHandler(name string, f SourceHandler) {
	sourceHandlers[name] = f
}
