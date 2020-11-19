package compose

import (
	"github.com/docker/cli/cli/compose/types"
	"github.com/yankghjh/selfhosted_store/cli/project"
	"gopkg.in/yaml.v2"
)

func init() {
	project.RegisterDecoder("docker-compose", Decoder)
	project.RegisterEncoder("docker-compose", Encoder)
}

// Decoder for docker-compose.yml
func Decoder(payload []byte) (*project.Application, error) {
	a := project.NewApplication()

	err := LoadApplication(a, payload)
	if err != nil {
		return nil, err
	}

	return a, nil
}

// Encoder for docker-compose.yml
func Encoder(a *project.Application) ([]byte, error) {
	var services types.Services = make([]types.ServiceConfig, len(a.Services))

	for i, service := range a.Services {
		services[i] = *service
	}

	cfg := &types.Config{
		Version:  "3.0",
		Services: services,
	}

	out, err := yaml.Marshal(cfg)
	if err != nil {
		return nil, err
	}

	return out, nil
}
