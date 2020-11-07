package portainer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/yankghjh/selfhosted_store/cli/project"
)

func init() {
	project.RegisterGenerater("portainer", Generater)
}

// Dataset portainer template format
type Dataset struct {
	Version   string      `json:"version"`
	Templates []*Template `json:"templates"`
}

// Template struct for portainer application
type Template struct {
	Type        int      `json:"type"`
	Title       string   `json:"title"`
	Description string   `json:"description,omitempty"`
	Categories  []string `json:"categories,omitempty"`
	Platform    string   `json:"platform"`
	Note        string   `json:"note,omitempty"`
	Logo        string   `json:"logo,omitempty"`

	Name          string              `json:"name"`
	Command       string              `json:"command,omitempty"`
	Image         string              `json:"image"`
	RestartPolicy string              `json:"restart_policy,omitempty"`
	NetworkMode   string              `json:"network_mode,omitempty"`
	Ports         []string            `json:"ports,omitempty"`
	Volumes       []VolumeConfig      `json:"volumes,omitempty"`
	Environment   []EnvironmentConfig `json:"env,omitempty"`
}

// VolumeConfig for portainer template volumn bind
type VolumeConfig struct {
	Container string `json:"container"`
	Bind      string `json:"bind"`
}

// EnvironmentConfig for portainer template environment
// TODO label and description
type EnvironmentConfig struct {
	Name        string `json:"name"`
	Label       string `json:"label,omitempty"`
	Default     string `json:"default,omitempty"`
	Description string `json:"description,omitempty"`
}

// Generater portainer template
func Generater(o *project.Operator) error {
	res, err := Convert(o.Project.Apps)
	if err != nil {
		return err
	}

	path := o.Project.GetDistPath("templates", "portainer")
	os.MkdirAll(path, os.ModePerm)
	filename := path + "/template.json"
	err = ioutil.WriteFile(filename, res, os.ModePerm)
	if err != nil {
		return fmt.Errorf("write template file [%s] error: %s", filename, err.Error())
	}

	return nil
}

// Convert applications to portainer template
func Convert(apps []*project.Application) ([]byte, error) {
	dataset := &Dataset{
		Version:   "2",
		Templates: []*Template{},
	}

	for _, app := range apps {
		t, err := ConvertApplication(app)
		if err != nil {
			return nil, fmt.Errorf("convert application to yacht template error: %s", err.Error())
		}
		dataset.Templates = append(dataset.Templates, t)
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
		t.Ports = []string{}

		for _, port := range service.Ports {
			published := strconv.Itoa(int(port.Published))
			target := strconv.Itoa(int(port.Target))
			t.Ports = append(t.Ports, published+":"+target+"/"+port.Protocol)
		}

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
