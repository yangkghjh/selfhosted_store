package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/yankghjh/store/cli/generator"

	"github.com/yankghjh/store/cli/app"
)

var appspath = "./apps"

func main() {
	files, _ := ioutil.ReadDir(appspath)

	apps := []*app.App{}
	for _, f := range files {
		if f.IsDir() {
			app, err := handleApp(f)
			if err != nil {
				fmt.Printf("handle %s error: %s\n", f.Name(), err)
				os.Exit(1)
			}

			apps = append(apps, app)
		}
	}

	templete, errs, err := generator.GenerateYachtTemplate(apps)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if errs != nil {
		for _, e := range errs {
			fmt.Println(e)
		}
	}

	fmt.Println(string(templete))
}

func handleApp(f os.FileInfo) (*app.App, error) {
	appFile := appspath + "/" + f.Name() + "/" + "app.yml"
	appYaml, err := ioutil.ReadFile(appFile)
	if err != nil {
		return nil, err
	}

	composeFile := appspath + "/" + f.Name() + "/" + "docker-compose.yml"
	composeYaml, err := ioutil.ReadFile(composeFile)
	if err != nil {
		return nil, err
	}

	app, err := app.NewApp(appYaml)
	if err != nil {
		return nil, err
	}

	err = app.LoadDockerCompose(composeYaml)
	if err != nil {
		return nil, err
	}

	return app, nil
}
