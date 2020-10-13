package app

import (
	"os"

	"github.com/docker/cli/cli/compose/loader"
	"github.com/docker/cli/cli/compose/types"
	"gopkg.in/yaml.v3"
)

// NewApp form yml format string
func NewApp(payload []byte) (*App, error) {
	app := new(App)
	err := yaml.Unmarshal([]byte(payload), app)

	if err != nil {
		return nil, err
	}

	return app, nil
}

// LoadDockerCompose from payload, the payload is the content of docker-compose.yml
func (app *App) LoadDockerCompose(payload []byte) error {
	source, err := loader.ParseYAML(payload)

	if err != nil {
		return err
	}

	workingDir, err := os.Getwd()
	if err != nil {
		return err
	}

	compose, err := loader.Load(types.ConfigDetails{
		WorkingDir: workingDir,
		ConfigFiles: []types.ConfigFile{
			{Filename: "filename.yml", Config: source},
		},
		Environment: map[string]string{},
	})

	if err != nil {
		return err
	}

	app.Compose = compose

	return nil
}
