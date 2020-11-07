package app

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/viper"
	compose "github.com/yankghjh/selfhosted_store/cli/modules/docker-compose"
	"github.com/yankghjh/selfhosted_store/cli/project"
)

// LoadApp from app.yml
func LoadApp(ctx *Context, a *project.Application) error {
	path := ctx.GetPath("app.yml")
	data := viper.New()
	data.SetConfigType("yml")

	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open file %s error: %s", path, err.Error())
	}

	err = data.ReadConfig(file)
	if err != nil {
		return fmt.Errorf("read file %s error: %s", path, err.Error())
	}

	a.Name = data.GetString("name")
	a.Description = data.GetString("description")
	a.Overview = data.GetString("overview")

	return nil
}

// LoadDockerCompose from docker-compose.yml
func LoadDockerCompose(ctx *Context, a *project.Application) error {
	path := ctx.GetPath("docker-compose.yml")
	payload, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read file %s error: %s", path, err.Error())
	}

	compose.LoadApplication(a, payload)
	if err != nil {
		return fmt.Errorf("load application form %s error: %s", path, err.Error())
	}

	return nil
}

// LoadIcon from icon.png
func LoadIcon(ctx *Context, a *project.Application) error {
	path := ctx.GetPath("icon.png")

	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return fmt.Errorf("stat icon %s error: %s", path, err.Error())
		}
		return nil
	}

	iconDistpath := ctx.Config.GetString("icon.distpath")
	if iconDistpath == "" {
		iconDistpath = "assets/icon"
	}
	folderPath := ctx.Operator.Project.Dist + "/" + iconDistpath
	os.MkdirAll(folderPath, os.ModePerm)
	distpath := folderPath + "/" + ctx.Name + ".png"

	err = copyFile(path, distpath)
	if err != nil {
		return fmt.Errorf("copy icon %s error: %s", path, err.Error())
	}

	a.Icon = ctx.Config.GetString("icon.basepath") + ctx.Name + ".png"

	return nil
}

func copyFile(src, dst string) error {
	input, err := ioutil.ReadFile(src)
	if err != nil {
		return fmt.Errorf("read file %s error: %s", src, err.Error())
	}

	err = ioutil.WriteFile(dst, input, 0644)
	if err != nil {
		return fmt.Errorf("write file %s error: %s", dst, err.Error())
	}

	return nil
}
