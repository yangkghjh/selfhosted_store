package generator

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/yankghjh/store/cli/app"
)

// YachtTemplate struct for Yacht template
type YachtTemplate struct {
	Type          int `json:"type"`
	app.Metadata  `json:",inline"`
	Image         string                   `json:"image"`
	RestartPolicy string                   `json:"restart_policy,omitempty"`
	Ports         []map[string]string      `json:"ports,omitempty"`
	Volumes       []YachtVolumeConfig      `json:"volumes,omitempty"`
	Environment   []YachtEnvironmentConfig `json:"env,omitempty"`
}

// YachtVolumeConfig for Yacht template volumn bind
type YachtVolumeConfig struct {
	Container string `json:"container"`
	Bind      string `json:"bind"`
}

// YachtEnvironmentConfig for Yacht template environment
type YachtEnvironmentConfig struct {
	Name    string `json:"name"`
	Label   string `json:"label"`
	Default string `json:"default"`
}

// GenerateYachtTemplate from App and DockerCompose struct
func GenerateYachtTemplate(apps []*app.App) ([]byte, []error, error) {
	template := []*YachtTemplate{}
	errs := []error{}

	for _, app := range apps {
		t, err := app2YachtTemplate(app)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		template = append(template, t)
	}

	res, err := json.MarshalIndent(template, "", "  ")

	return res, errs, err
}

func app2YachtTemplate(app *app.App) (*YachtTemplate, error) {
	if app == nil {
		return nil, errors.New("nil app")
	}
	t := new(YachtTemplate)

	t.Metadata = app.Metadata
	t.Type = 1

	service := app.Compose.Services[0]

	t.Image = service.Image
	if service.Restart == "always" {
		t.RestartPolicy = "unless-stopped"
	}

	if len(service.Ports) > 0 {
		t.Ports = []map[string]string{}
		ports := map[string]string{}

		for _, port := range service.Ports {
			published := strconv.Itoa(int(port.Published))
			target := strconv.Itoa(int(port.Target))
			ports[target] = published + ":" + target + "/" + port.Protocol
		}

		t.Ports = []map[string]string{ports}
	}

	if len(service.Volumes) > 0 {
		t.Volumes = []YachtVolumeConfig{}

		for _, volumn := range service.Volumes {
			t.Volumes = append(t.Volumes, YachtVolumeConfig{
				Container: volumn.Target,
				Bind:      volumn.Source,
			})
		}
	}

	if len(service.Environment) > 0 {
		t.Environment = []YachtEnvironmentConfig{}

		for name, point := range service.Environment {
			value := ""
			if point != nil {
				value = *point
			}
			t.Environment = append(t.Environment, YachtEnvironmentConfig{
				Name:    name,
				Label:   name,
				Default: value,
			})
		}
	}

	return t, nil
}
