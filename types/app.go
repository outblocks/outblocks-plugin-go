package types

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/creasty/defaults"
	"github.com/mitchellh/mapstructure"
	apiv1 "github.com/outblocks/outblocks-plugin-go/gen/api/v1"
	"github.com/outblocks/outblocks-plugin-go/util"
)

func AppEnvPrefix(a *apiv1.App) string {
	return fmt.Sprintf("APP_%s_%s_", strings.ToUpper(a.Type), util.SanitizeEnvVar(strings.ToUpper(a.Name)))
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

func encodeToMap(in interface{}) (map[string]interface{}, error) {
	out := make(map[string]interface{})

	b, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, &out)

	return out, err
}

// Service app properties.

type ServiceAppBuild struct {
	Dockerfile      string            `json:"dockerfile"`
	DockerContext   string            `json:"context"`
	DockerBuildArgs map[string]string `json:"build_args"`
}

type ServiceAppContainer struct {
	Port int `json:"port" default:"8080"`
}

type ServiceAppCDN struct {
	Enabled bool `json:"enabled"`
}

type ServiceAppProperties struct {
	Private bool `json:"private"`

	Build     *ServiceAppBuild     `json:"build,omitempty"`
	Container *ServiceAppContainer `json:"container,omitempty"`
	CDN       *ServiceAppCDN       `json:"cdn,omitempty"`

	LocalDockerImage string `json:"local_docker_image"`
	LocalDockerHash  string `json:"local_docker_hash"`
}

func NewServiceAppProperties(in map[string]interface{}) (*ServiceAppProperties, error) {
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
	return encodeToMap(p)
}

// Static app properties.

type StaticAppBuild struct {
	Env     map[string]string `json:"env,omitempty"`
	Command string            `json:"command"`
	Dir     string            `json:"dir"`
}

type StaticAppCDN struct {
	Enabled bool `json:"enabled"`
}

type StaticAppProperties struct {
	Build *StaticAppBuild `json:"build,omitempty"`
	CDN   *StaticAppCDN   `json:"cdn,omitempty"`

	Routing string `json:"routing"`
}

func NewStaticAppProperties(in map[string]interface{}) (*StaticAppProperties, error) {
	o := &StaticAppProperties{
		Build: &StaticAppBuild{},
		CDN:   &StaticAppCDN{},
	}

	return o, mapstructureJSONDecode(in, o)
}

func (p *StaticAppProperties) Encode() (map[string]interface{}, error) {
	return encodeToMap(p)
}

type AppVars map[string]interface{}

func VarsFromAppType(app *apiv1.App) map[string]interface{} {
	return map[string]interface{}{
		"url": app.Url,
	}
}

func VarsFromAppRunType(app *apiv1.AppRun) map[string]interface{} {
	return map[string]interface{}{
		"url": app.Url,
	}
}

func AppVarsFromApps(apps []*apiv1.App) AppVars {
	appVars := make(map[string]interface{}) // type->name->value

	for _, app := range apps {
		vars := VarsFromAppType(app)

		if _, ok := appVars[app.Type]; !ok {
			appVars[app.Type] = map[string]interface{}{
				app.Name: vars,
			}
		} else {
			appVars[app.Type].(map[string]interface{})[app.Name] = vars
		}
	}

	return appVars
}

func AppVarsFromAppRun(apps []*apiv1.AppRun) AppVars {
	appVars := make(map[string]interface{}) // type->name->value

	for _, app := range apps {
		vars := VarsFromAppRunType(app)

		if _, ok := appVars[app.App.Type]; !ok {
			appVars[app.App.Type] = map[string]interface{}{
				app.App.Name: vars,
			}
		} else {
			appVars[app.App.Type].(map[string]interface{})[app.App.Name] = vars
		}
	}

	return appVars
}

func (v AppVars) ForApp(a *apiv1.App) map[string]interface{} {
	return v[a.Type].(map[string]interface{})[a.Name].(map[string]interface{})
}

func VarsForApp(av AppVars, a *apiv1.App, depVars interface{}) map[string]interface{} {
	return map[string]interface{}{
		"app":  map[string]interface{}(av),
		"self": av.ForApp(a),
		"dep":  depVars,
	}
}
