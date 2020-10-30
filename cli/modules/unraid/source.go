package unraid

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/spf13/cast"

	"github.com/docker/cli/cli/compose/types"
	compose "github.com/yankghjh/selfhosted_store/cli/modules/docker-compose"
	"github.com/yankghjh/selfhosted_store/cli/modules/icon"

	"github.com/spf13/viper"
	"github.com/yankghjh/selfhosted_store/cli/modules/app"

	"github.com/yankghjh/selfhosted_store/cli/pipe"
)

func init() {
	pipe.RegisterInitFunc("unraid", InitFunc)
}

// Application struct for unraid community application feed
type Application struct {
	Plugin      bool
	Name        string
	Description string
	Icon        string
	Repository  string
	Environment interface{}
	Networking  interface{}
	Data        interface{}
}

// FeedFile struct for unraid community application feed
type FeedFile struct {
	Apps    int            `json:"apps"`
	AppList []*Application `json:"applist"`
}

// InitFunc for unraid community applications
func InitFunc(p *pipe.Pipe) error {
	if !p.Config.GetBool("unraid.enable") {
		return nil
	}

	path := p.Config.GetString("unraid.application_feed_file")

	payload, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read file %s error: %s", path, err.Error())
	}

	feed := &FeedFile{
		AppList: []*Application{},
	}

	err = json.Unmarshal(payload, feed)

	if err != nil {
		return fmt.Errorf("unmarshal unraid application_feed_file error: %s", err.Error())
	}

	for _, a := range feed.AppList {
		if a.Plugin || a.Repository == "" {
			continue
		}

		ctx := pipe.NewContext(a.Name, "")

		ctx.Set("app", a.GetApp())
		ctx.Set("icon", a.GetIcon())

		// docker compose
		service := types.ServiceConfig{}
		service.ContainerName = a.Name
		service.Environment = a.GetEnvironment()
		service.Ports = a.GetPorts()
		service.Volumes = a.GetVolumns()
		service.Image = a.Repository

		ctx.Set("compose", &compose.Compose{
			Config: &types.Config{
				Services: []types.ServiceConfig{service},
			},
		})

		p.Apps = append(p.Apps, ctx)
	}

	return nil
}

// GetApp parse app metadata from application
func (a *Application) GetApp() *app.App {
	data := viper.New()
	data.Set("description", a.Description)
	appyml := &app.App{
		Name: a.Name,
		Data: data,
	}

	return appyml
}

// GetIcon parse icon data from application
func (a *Application) GetIcon() *icon.Icon {
	return &icon.Icon{
		URL: a.Icon,
	}
}

// GetEnvironment parse environment of application
func (a *Application) GetEnvironment() map[string]*string {
	env := map[string]*string{}
	if a.Environment == nil {
		return env
	}

	v, ok := cast.ToStringMap(a.Environment)["Variable"]
	if !ok {
		return env
	}

	vs := cast.ToSlice(v)
	if len(vs) == 0 {
		vs = []interface{}{v}
	}

	for _, pair := range vs {
		kv := cast.ToStringMapString(pair)
		key := kv["Name"]
		if key != "" {
			value := kv["Value"]
			env[key] = &value
		}
	}

	return env
}

// GetPorts parse ports of application
func (a *Application) GetPorts() []types.ServicePortConfig {
	ports := []types.ServicePortConfig{}
	if a.Networking == nil {
		return ports
	}

	publish, ok := cast.ToStringMap(a.Networking)["Publish"]
	if !ok {
		return ports
	}

	pts, ok := cast.ToStringMap(publish)["Port"]
	if !ok {
		return ports
	}

	ps := cast.ToSlice(pts)
	if len(ps) == 0 {
		ps = []interface{}{pts}
	}

	for _, port := range ps {
		kv := cast.ToStringMapString(port)
		if kv["ContainerPort"] != "" {
			protocol := kv["Protocol"]
			if protocol != "tcp" && protocol != "udp" {
				protocol = "tcp"
			}
			ports = append(ports, types.ServicePortConfig{
				Published: cast.ToUint32(kv["HostPort"]),
				Target:    cast.ToUint32(kv["ContainerPort"]),
				Protocol:  protocol,
			})
		}
	}

	return ports
}

// GetVolumns parse environment of application
func (a *Application) GetVolumns() []types.ServiceVolumeConfig {
	volumes := []types.ServiceVolumeConfig{}
	if a.Networking == nil {
		return volumes
	}
	vs, ok := cast.ToStringMap(a.Data)["Volume"]
	if !ok {
		return volumes
	}

	vos := cast.ToSlice(vs)
	if len(vos) == 0 {
		vos = []interface{}{vs}
	}

	for _, v := range vos {
		kv := cast.ToStringMapString(v)
		if kv["ContainerDir"] != "" {
			volumes = append(volumes, types.ServiceVolumeConfig{
				Target: kv["ContainerDir"],
				Source: kv["HostDir"],
			})
		}
	}
	return volumes
}
