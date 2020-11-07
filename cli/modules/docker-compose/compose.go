package compose

import (
	"fmt"
	"io/ioutil"

	"github.com/docker/cli/cli/compose/loader"
	"github.com/docker/cli/cli/compose/types"
	"github.com/yankghjh/selfhosted_store/cli/pipe"
	"github.com/yankghjh/selfhosted_store/cli/project"
)

func init() {
	pipe.RegisterSourceLoader("docker-compose", Loader)
}

// Compose struct for docker-compose
type Compose struct {
	Config *types.Config
}

// Loader for docker-compose.yml
func Loader(pipe *pipe.Pipe, ctx *pipe.Context) error {
	path := ctx.GetPath("docker-compose.yml")

	payload, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read file %s error: %s", path, err.Error())
	}

	source, err := loader.ParseYAML(payload)
	if err != nil {
		return fmt.Errorf("parse yaml %s error: %s", path, err.Error())
	}

	config, err := loader.Load(types.ConfigDetails{
		ConfigFiles: []types.ConfigFile{
			{Filename: "docker-compose.yml", Config: source},
		},
		Environment: map[string]string{},
	})
	if err != nil {
		return fmt.Errorf("load docker compose conifg error: %s", err.Error())
	}

	ctx.Set("compose", &Compose{
		Config: config,
	})

	return nil
}

// LoadApplication from docker-compose.yml
func LoadApplication(a *project.Application, payload []byte) error {
	source, err := loader.ParseYAML(payload)
	if err != nil {
		return fmt.Errorf("parse docker compose yaml error: %s", err.Error())
	}

	config, err := loader.Load(types.ConfigDetails{
		ConfigFiles: []types.ConfigFile{
			{Filename: "docker-compose.yml", Config: source},
		},
		Environment: map[string]string{},
	})
	if err != nil {
		return fmt.Errorf("load docker compose conifg error: %s", err.Error())
	}

	if len(config.Services) == 0 {
		return fmt.Errorf("load docker compose services error: no service found")
	}

	for _, service := range config.Services {
		a.Services = append(a.Services, &service)
	}

	return nil
}
