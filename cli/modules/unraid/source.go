package unraid

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/docker/cli/cli/compose/types"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	compose "github.com/yankghjh/selfhosted_store/cli/modules/docker-compose"

	"github.com/yankghjh/selfhosted_store/cli/modules/app"
	"github.com/yankghjh/selfhosted_store/cli/modules/icon"
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
	Config      interface{}

	network     map[string]*types.ServicePortConfig
	volumn      map[string]*types.ServiceVolumeConfig
	environment map[string]*string
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

	appList, err := LoadApplications(payload)
	if err != nil {
		return err
	}

	for _, a := range appList {
		if a.Plugin || a.Repository == "" {
			continue
		}

		a.Parse()

		ctx := pipe.NewContext(a.Name, "")

		ctx.Set("app", a.GetApp())
		ctx.Set("icon", a.GetIcon())

		// docker compose
		ctx.Set("compose", &compose.Compose{
			Config: &types.Config{
				Services: []types.ServiceConfig{*a.GetServiceConfig()},
			},
		})

		p.Apps = append(p.Apps, ctx)
	}

	return nil
}

// LoadApplications form application feed file data
func LoadApplications(payload []byte) ([]*Application, error) {
	feed := &FeedFile{
		AppList: []*Application{},
	}

	err := json.Unmarshal(payload, feed)

	if err != nil {
		return nil, fmt.Errorf("unmarshal unraid application_feed_file error: %s", err.Error())
	}

	return feed.AppList, nil
}

// Parse unraid application after unmarshaled form json
func (a *Application) Parse() error {
	a.network = map[string]*types.ServicePortConfig{}
	a.volumn = map[string]*types.ServiceVolumeConfig{}
	a.environment = map[string]*string{}
	// config
	if a.Config != nil {
		cfgs := []map[string]interface{}{}
		cos := cast.ToSlice(a.Config)
		if len(cos) > 0 {
			for _, c := range cos {
				if cfg := cast.ToStringMap(c); cfg != nil {
					cfgs = append(cfgs, cfg)
				}
			}
		} else if cfg := cast.ToStringMap(a.Config); cfg != nil {
			cfgs = append(cfgs, cfg)
		}

		for _, item := range cfgs {
			a.parseConfigItem(item)
		}
	}

	a.parseEnvironment()
	a.parsePorts()
	a.parseVolumns()

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

// GetServiceConfig from application
func (a *Application) GetServiceConfig() *types.ServiceConfig {
	service := &types.ServiceConfig{}
	service.ContainerName = a.Name
	service.Environment = a.environment
	service.Image = a.Repository

	ports := []types.ServicePortConfig{}
	for _, p := range a.network {
		ports = append(ports, *p)
	}
	service.Ports = ports

	volumns := []types.ServiceVolumeConfig{}
	for _, v := range a.volumn {
		volumns = append(volumns, *v)
	}
	service.Volumes = volumns

	return service
}

func (a *Application) parseConfigItem(v map[string]interface{}) {
	attributes := cast.ToStringMapString(v["@attributes"])
	value := cast.ToString(v["value"])
	if value == "" {
		value = attributes["Default"]
	}

	switch attributes["Type"] {
	case "Port":
		a.addNetwork(&types.ServicePortConfig{
			Published: cast.ToUint32(value),
			Target:    cast.ToUint32(attributes["Target"]),
			Protocol:  attributes["Mode"],
		})
	case "Path":
		a.addVolumn(&types.ServiceVolumeConfig{
			Target: attributes["Target"],
			Source: value,
		})
	case "Variable":
		a.addEnvironment(attributes["Target"], value)
	}
}

func (a *Application) addNetwork(n *types.ServicePortConfig) {
	if n.Mode != "tcp" && n.Mode != "udp" {
		n.Mode = "tcp"
	}

	name := strconv.Itoa(int(n.Target)) + "/" + n.Mode

	if _, isExisted := a.network[name]; isExisted {
		return
	}

	a.network[name] = n
}

func (a *Application) addVolumn(v *types.ServiceVolumeConfig) {
	if _, isExisted := a.volumn[v.Target]; isExisted {
		return
	}

	a.volumn[v.Target] = v
}

func (a *Application) addEnvironment(key, value string) {
	a.environment[key] = &value
}

// parse environment of application
func (a *Application) parseEnvironment() {
	if a.Environment == nil {
		return
	}

	v, ok := cast.ToStringMap(a.Environment)["Variable"]
	if !ok {
		return
	}

	vs := cast.ToSlice(v)
	if len(vs) == 0 {
		vs = []interface{}{v}
	}

	for _, pair := range vs {
		kv := cast.ToStringMapString(pair)
		key := kv["Name"]
		if key != "" {
			a.addEnvironment(key, kv["Value"])
		}
	}
}

// parse ports of application
func (a *Application) parsePorts() {
	if a.Networking == nil {
		return
	}

	publish, ok := cast.ToStringMap(a.Networking)["Publish"]
	if !ok {
		return
	}

	pts, ok := cast.ToStringMap(publish)["Port"]
	if !ok {
		return
	}

	ps := cast.ToSlice(pts)
	if len(ps) == 0 {
		ps = []interface{}{pts}
	}

	for _, port := range ps {
		kv := cast.ToStringMapString(port)
		if kv["ContainerPort"] != "" {
			protocol := kv["Protocol"]
			a.addNetwork(&types.ServicePortConfig{
				Published: cast.ToUint32(kv["HostPort"]),
				Target:    cast.ToUint32(kv["ContainerPort"]),
				Protocol:  protocol,
			})
		}
	}
}

// parse volumns of application
func (a *Application) parseVolumns() {
	if a.Data == nil {
		return
	}

	vs, ok := cast.ToStringMap(a.Data)["Volume"]
	if !ok {
		return
	}

	vos := cast.ToSlice(vs)
	if len(vos) == 0 {
		vos = []interface{}{vs}
	}

	for _, v := range vos {
		kv := cast.ToStringMapString(v)
		if kv["ContainerDir"] != "" {
			a.addVolumn(&types.ServiceVolumeConfig{
				Target: kv["ContainerDir"],
				Source: kv["HostDir"],
			})
		}
	}
}
