package types

import (
	"fmt"
	"strings"

	"github.com/creasty/defaults"
	"github.com/mitchellh/mapstructure"
	"github.com/outblocks/outblocks-plugin-go/util"
)

type App struct {
	ID           string                 `json:"id"`
	DeployPlugin string                 `json:"deploy_plugin"`
	Dir          string                 `json:"dir"`
	Name         string                 `json:"name"`
	Type         string                 `json:"type"`
	URL          string                 `json:"url"`
	PathRedirect string                 `json:"path_redirect"`
	Env          map[string]string      `json:"env"`
	Needs        map[string]*AppNeed    `json:"needs"`
	Properties   map[string]interface{} `json:"properties"`
}

func (a *App) EnvPrefix() string {
	return fmt.Sprintf("APP_%s_%s_", strings.ToUpper(a.Type), util.SanitizeEnvVar(strings.ToUpper(a.Name)))
}

func (a *App) String() string {
	return fmt.Sprintf("App<Name=%s,Type=%s>", a.Name, a.Type)
}

type AppNeed struct {
	Dependency string                 `json:"dependency"`
	Properties map[string]interface{} `json:"properties"`
}

func (a *AppNeed) String() string {
	return fmt.Sprintf("AppNeed<Dep=%s>", a.Dependency)
}

func mapstructureJSONDecode(in, out interface{}) error {
	cfg := &mapstructure.DecoderConfig{
		Metadata: nil,
		Result:   out,
		TagName:  "json",
	}

	decoder, err := mapstructure.NewDecoder(cfg)
	if err != nil {
		return err
	}

	return decoder.Decode(in)
}

// Service app properties.

type ServiceAppBuild struct {
	Dockerfile      string            `json:"dockerfile"`
	DockerContext   string            `json:"context"`
	DockerBuildArgs map[string]string `json:"build_args"`
}

type ServiceAppContainer struct {
	Port int `json:"port" default:"80"`
}

type ServiceAppCDN struct {
	Enabled bool `json:"enabled"`
}

type ServiceAppProperties struct {
	Public bool `json:"public" default:"true"`

	Build     *ServiceAppBuild     `json:"build"`
	Container *ServiceAppContainer `json:"container"`
	CDN       *ServiceAppCDN       `json:"cdn"`

	LocalDockerImage string `json:"local_docker_image"`
	LocalDockerHash  string `json:"local_docker_hash"`
}

func NewServiceAppProperties(in interface{}) (*ServiceAppProperties, error) {
	o := &ServiceAppProperties{
		Build:     &ServiceAppBuild{},
		Container: &ServiceAppContainer{},
		CDN:       &ServiceAppCDN{},
	}

	err := mapstructureJSONDecode(in, o)
	if err != nil {
		return nil, err
	}

	return o, defaults.Set(o)
}

func (p *ServiceAppProperties) Encode() (map[string]interface{}, error) {
	out := make(map[string]interface{})
	err := mapstructureJSONDecode(p, &out)

	return out, err
}

// Static app properties.

type StaticAppBuild struct {
	Env     map[string]string `json:"env"`
	Command string            `json:"command"`
	Dir     string            `json:"dir"`
}

type StaticAppCDN struct {
	Enabled bool `json:"enabled"`
}

type StaticAppProperties struct {
	Build *StaticAppBuild `json:"build"`
	CDN   *StaticAppCDN   `json:"cdn"`

	Routing string `json:"routing"`
}

func NewStaticAppProperties(in interface{}) (*StaticAppProperties, error) {
	o := &StaticAppProperties{
		Build: &StaticAppBuild{},
		CDN:   &StaticAppCDN{},
	}

	return o, mapstructureJSONDecode(in, o)
}

func (p *StaticAppProperties) Encode() (map[string]interface{}, error) {
	out := make(map[string]interface{})
	err := mapstructureJSONDecode(p, &out)

	return out, err
}
