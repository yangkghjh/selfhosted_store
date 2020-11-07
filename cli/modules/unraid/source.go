package unraid

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/docker/cli/cli/compose/types"
	"github.com/spf13/cast"

	"github.com/yankghjh/selfhosted_store/cli/project"
)

func init() {
	project.RegisterLoader("unraid", Loader)
}

// Application struct for unraid community application feed
type Application struct {
	Plugin      bool
	Name        string
	Description string
	Overview    string
	Category    string
	Icon        string
	Repository  string
	Environment interface{}
	Networking  interface{}
	Data        interface{}
	Config      interface{}

	network     map[string]*types.ServicePortConfig
	volumn      map[string]*types.ServiceVolumeConfig
	environment map[string]*string
	networkMode string
}

// FeedFile struct for unraid community application feed
type FeedFile struct {
	Apps    int            `json:"apps"`
	AppList []*Application `json:"applist"`
}

var defaultRestartPolicy string = "unless-stopped"

// Loader for unraid community applications
func Loader(o *project.Operator) error {
	path := o.Config.GetString("application_feed_file")

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
		o.Project.Apps = append(o.Project.Apps, a.ToProjectApplication())
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

// ToProjectApplication convert to project application
func (a *Application) ToProjectApplication() *project.Application {
	app := project.NewApplication()

	app.Name = a.Name
	app.Description = a.Description
	app.Overview = a.Overview
	app.Icon = a.Icon
	app.Category = strings.Split(a.Category, " ")

	app.Services = append(app.Services, a.GetServiceConfig())

	return app
}

// GetServiceConfig from application
func (a *Application) GetServiceConfig() *types.ServiceConfig {
	service := &types.ServiceConfig{}
	service.ContainerName = a.Name
	service.Environment = a.environment
	service.Image = a.Repository
	service.Restart = defaultRestartPolicy
	service.NetworkMode = a.networkMode

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
	if n.Protocol != "tcp" && n.Protocol != "udp" {
		n.Protocol = "tcp"
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

	network := cast.ToStringMap(a.Networking)
	if m, ok := network["Mode"]; ok {
		mode := cast.ToString(m)
		if mode == "bridge" || mode == "host" || mode == "none" {
			a.networkMode = mode
		}
	}

	publish, ok := network["Publish"]
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
