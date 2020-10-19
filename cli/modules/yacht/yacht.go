package yacht

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/yankghjh/selfhosted_store/cli/modules/icon"

	"github.com/yankghjh/selfhosted_store/cli/modules/app"
	compose "github.com/yankghjh/selfhosted_store/cli/modules/docker-compose"
	"github.com/yankghjh/selfhosted_store/cli/pipe"
)

func init() {
	pipe.RegisterInitFunc("yacht", InitFunc)
	pipe.RegisterFinishFunc("yacht", FinishFunc)
	pipe.RegisterSourceHandler("yacht", Handler)
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
	Categories  []string `json:"categories"`
	Platform    string   `json:"platform"`
	Note        string   `json:"note,omitempty"`
	Logo        string   `json:"logo,omitempty"`

	Name          string              `json:"name"`
	Image         string              `json:"image"`
	RestartPolicy string              `json:"restart_policy,omitempty"`
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

// FinishFunc write templetes to file
func FinishFunc(pipe *pipe.Pipe) error {
	dataset := pipe.Get("yacht").(*Dataset)

	path := pipe.GetDistPath("yacht")
	os.MkdirAll(path, os.ModePerm)

	res, err := json.MarshalIndent(dataset.Templates, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal templates error: %s", err.Error())
	}

	filename := path + "/yacht.json"
	ioutil.WriteFile(filename, res, os.ModePerm)
	if err != nil {
		return fmt.Errorf("write template file [%s] error: %s", filename, err.Error())
	}

	return nil
}

// Handler generate yacht app definition
func Handler(pipe *pipe.Pipe, ctx *pipe.Context) error {
	t := new(Template)
	a := ctx.Get("app").(*app.App)
	c := ctx.Get("compose").(*compose.Compose)
	if len(c.Config.Services) < 1 {
		return fmt.Errorf("no service found in compose")
	}
	service := c.Config.Services[0]

	t.Type = 1
	t.Title = a.Name
	t.Description = a.Data.GetString("description")
	t.Categories = a.Data.GetStringSlice("categories")
	t.Platform = "linux"
	t.Note = a.Data.GetString("note")
	if v := ctx.Get("icon"); v != nil {
		t.Logo = v.(*icon.Icon).URL
	}

	t.Name = service.ContainerName
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

	dataset := pipe.Get("yacht").(*Dataset)
	dataset.Templates = append(dataset.Templates, t)

	return nil
}
