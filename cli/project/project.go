package project

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type (
	// Loader application
	Loader func(o *Operator) error
	// Generater site/template
	Generater func(o *Operator) error
)

var (
	loaders    map[string]Loader
	generaters map[string]Generater
)

func init() {
	loaders = make(map[string]Loader)
	generaters = make(map[string]Generater)
}

// RegisterLoader register loader
func RegisterLoader(name string, f Loader) {
	loaders[name] = f
}

// RegisterGenerater register generater
func RegisterGenerater(name string, f Generater) {
	generaters[name] = f
}

// Project is the universal data struct for selfhosted template
type Project struct {
	Dist       string
	Loaders    []*Operator
	Generaters []*Operator
	Config     *viper.Viper
	Apps       []*Application
}

// Operator loader or generater
type Operator struct {
	Name    string
	Type    string
	Config  *viper.Viper
	Project *Project
}

// NewProject create new project
func NewProject(cfg *viper.Viper) *Project {
	return &Project{
		Config:     cfg,
		Dist:       cfg.GetString("dist"),
		Apps:       []*Application{},
		Loaders:    []*Operator{},
		Generaters: []*Operator{},
	}
}

// Run the project
func (p *Project) Run() error {
	for _, o := range p.Loaders {
		err := loaders[o.Type](o)
		if err != nil {
			return fmt.Errorf("run loader %s(%s) error: %s", o.Name, o.Type, err)
		}
	}

	for _, o := range p.Generaters {
		err := generaters[o.Type](o)
		if err != nil {
			return fmt.Errorf("run generater %s(%s) error: %s", o.Name, o.Type, err)
		}
	}

	return nil
}

// AddLoader add loader
func (p *Project) AddLoader(name string) error {
	cfg := p.Config.Sub("loaders." + name)

	t := cfg.GetString("type")
	if _, ok := loaders[t]; !ok {
		return fmt.Errorf("no such loader type [%s]", t)
	}

	p.Loaders = append(p.Loaders, &Operator{
		Name:    name,
		Type:    t,
		Config:  cfg,
		Project: p,
	})

	return nil
}

// AddGenerater add generater
func (p *Project) AddGenerater(name string) error {
	cfg := p.Config.Sub("generaters." + name)

	t := cfg.GetString("type")
	if _, ok := generaters[t]; !ok {
		return fmt.Errorf("no such generater type [%s]", t)
	}

	p.Generaters = append(p.Generaters, &Operator{
		Name:    name,
		Type:    t,
		Config:  cfg,
		Project: p,
	})

	return nil
}

// GetDistPath of project
func (p *Project) GetDistPath(paths ...string) string {
	return p.Dist + "/" + strings.Join(paths, "/")
}
