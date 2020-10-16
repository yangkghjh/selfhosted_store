package app

import (
	"fmt"
	"os"

	"github.com/spf13/viper"

	"github.com/yankghjh/selfhosted_store/cli/pipe"
)

func init() {
	pipe.RegisterSourceLoader("app", Loader)
}

// App read from app.yml
type App struct {
	Name string
	Path string
	Data *viper.Viper
}

// Loader for app.yml
func Loader(pipe *pipe.Pipe, ctx *pipe.Context) error {
	path := ctx.GetPath("app.yml")
	data := viper.New()

	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open file %s error: %s", path, err.Error())
	}

	err = data.ReadConfig(file)
	if err != nil {
		return fmt.Errorf("read file %s error: %s", path, err.Error())
	}

	app := &App{
		Name: data.GetString("name"),
		Data: data,
		Path: path,
	}

	ctx.Set("app", app)

	return nil
}
