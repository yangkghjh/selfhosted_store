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
		if !a.Plugin {
			ctx := pipe.NewContext(a.Name, "")

			data := viper.New()
			data.Set("description", a.Description)
			appyml := &app.App{
				Name: a.Name,
				Data: data,
			}
			ctx.Set("app", appyml)

			iconURL := icon.Icon{
				URL: a.Icon,
			}
			ctx.Set("icon", &iconURL)

			// docker compose
			env := map[string]*string{}
			if a.Environment != nil {
				if v, ok := cast.ToStringMap(a.Environment)["Variable"]; ok {
					for _, pair := range cast.ToSlice(v) {
						kv := cast.ToStringMapString(pair)
						key := kv["Name"]
						if key != "" {
							value := kv["Value"]
							env[key] = &value
						}
					}
				}
			}
			ports := []types.ServicePortConfig{}
			if a.Networking != nil {
				if publish, ok := cast.ToStringMap(a.Networking)["Publish"]; ok {
					if pts, ok := cast.ToStringMap(publish)["Port"]; ok {
						for _, port := range cast.ToSlice(pts) {
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
					}
				}
			}
			volumes := []types.ServiceVolumeConfig{}
			if a.Networking != nil {
				if vs, ok := cast.ToStringMap(a.Data)["Volume"]; ok {
					for _, v := range cast.ToSlice(vs) {
						kv := cast.ToStringMapString(v)
						if kv["ContainerDir"] != "" {
							volumes = append(volumes, types.ServiceVolumeConfig{
								Target: kv["ContainerDir"],
								Source: kv["HostDir"],
							})
						}
					}
				}
			}
			dc := &types.Config{
				Services: []types.ServiceConfig{
					{
						ContainerName: a.Name,
						Environment:   env,
						Volumes:       volumes,
						Ports:         ports,
						Image:         a.Repository,
					},
				},
			}
			ctx.Set("compose", &compose.Compose{
				Config: dc,
			})

			p.Apps = append(p.Apps, ctx)
		}
	}

	return nil
}
