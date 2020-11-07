package yacht

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/yankghjh/selfhosted_store/cli/project"

	"github.com/yankghjh/selfhosted_store/cli/pipe"
)

func init() {
	project.RegisterGenerater("yacht", Generater)
}

// Dataset yacht app dataset
type Dataset struct {
	Templates []*Template
}

// Template struct for Yacht application
type Template struct {
	Type        int      `json:"type"`
	Title       string   `json:"title"`
	Description string   `json:"description,omitempty"`
	Categories  []string `json:"categories,omitempty"`
	Platform    string   `json:"platform"`
	Note        string   `json:"note,omitempty"`
	Logo        string   `json:"logo,omitempty"`

	Name          string              `json:"name"`
	Image         string              `json:"image"`
	RestartPolicy string              `json:"restart_policy,omitempty"`
	NetworkMode   string              `json:"network_mode,omitempty"`
	Ports         []map[string]string `json:"ports,omitempty"`
	Volumes       []VolumeConfig      `json:"volumes,omitempty"`
	Environment   []EnvironmentConfig `json:"env,omitempty"`
}

// VolumeConfig for Yacht template volumn bind
type VolumeConfig struct {
	Container string `json:"container"`
	Bind      string `json:"bind"`
}

// EnvironmentConfig for Yacht template environment
type EnvironmentConfig struct {
	Name    string `json:"name"`
	Label   string `json:"label"`
	Default string `json:"default"`
}

// InitFunc init dataset
func InitFunc(pipe *pipe.Pipe) error {
	pipe.Set("yacht", &Dataset{
		Templates: []*Template{},
	})

	return nil
}

// Generater yacht template
func Generater(o *project.Operator) error {
	res, err := Convert(o.Project.Apps)
	if err != nil {
		return err
	}

	path := o.Project.GetDistPath("templates", "yacht")
	os.MkdirAll(path, os.ModePerm)
	filename := path + "/yacht.json"
	err = ioutil.WriteFile(filename, res, os.ModePerm)
	if err != nil {
		return fmt.Errorf("write template file [%s] error: %s", filename, err.Error())
	}

	return nil
}

// Convert applications to yacht template
func Convert(apps []*project.Application) ([]byte, error) {
	dataset := []*Template{}

	for _, app := range apps {
		t, err := ConvertApplication(app)
		if err != nil {
			return nil, fmt.Errorf("convert application to yacht template error: %s", err.Error())
		}
		dataset = append(dataset, t)
	}

	res, err := json.MarshalIndent(dataset, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal templates error: %s", err.Error())
	}

	return res, nil
}

// ConvertApplication convert single application
func ConvertApplication(a *project.Application) (*Template, error) {
	t := new(Template)
	service := a.Services[0]

	t.Type = 1
	t.Title = a.Name
	t.Description = a.Overview
	t.Categories = a.Category
	t.Platform = "linux"
	t.Note = a.Description
	t.Logo = a.Icon

	t.Name = service.ContainerName
	t.Image = service.Image
	t.RestartPolicy = service.Restart
	t.NetworkMode = service.NetworkMode

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
		t.Volumes = []VolumeConfig{}

		for _, volumn := range service.Volumes {
			t.Volumes = append(t.Volumes, VolumeConfig{
				Container: volumn.Target,
				Bind:      volumn.Source,
			})
		}
	}

	if len(service.Environment) > 0 {
		t.Environment = []EnvironmentConfig{}

		for name, point := range service.Environment {
			value := ""
			if point != nil {
				value = *point
			}
			t.Environment = append(t.Environment, EnvironmentConfig{
				Name:    name,
				Label:   name,
				Default: value,
			})
		}
	}

	return t, nil
}
