package app

import (
	"fmt"
	"io/ioutil"

	"github.com/yankghjh/selfhosted_store/cli/project"
)

func init() {
	project.RegisterLoader("app", Loader)
}

// LoaderPlugin plugin for app loader
type LoaderPlugin func(*Context, *project.Application) error

// Loader load app from app path
func Loader(o *project.Operator) error {
	path := o.Config.GetString("path")
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return fmt.Errorf("read source path error: %s", err.Error())
	}

	ctxs := []*Context{}

	for _, f := range files {
		if f.IsDir() {
			ctx := NewContext(o, f.Name(), path+"/"+f.Name())
			ctxs = append(ctxs, ctx)
		}
	}

	for _, ctx := range ctxs {
		a := project.NewApplication()
		a.Name = ctx.Name

		plugins := []LoaderPlugin{LoadDockerCompose, LoadApp, LoadIcon}
		for _, f := range plugins {
			if err := f(ctx, a); err != nil {
				return fmt.Errorf("load for %s error: %s", a.Name, err.Error())
			}
		}

		o.Project.Apps = append(o.Project.Apps, a)
	}

	return nil
}
