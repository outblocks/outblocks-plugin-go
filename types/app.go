package types

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/creasty/defaults"
	apiv1 "github.com/outblocks/outblocks-plugin-go/gen/api/v1"
	"github.com/outblocks/outblocks-plugin-go/util"
	"github.com/outblocks/outblocks-plugin-go/util/command"
)

func AppEnvPrefix(a *apiv1.App) string {
	return fmt.Sprintf("APP_%s_%s_", strings.ToUpper(a.Type), util.SanitizeEnvVar(strings.ToUpper(a.Name)))
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

// Common properties.
type AppCDN struct {
	Enabled bool `json:"enabled"`
}

type AppScheduler struct {
	Cron    string            `json:"cron"`
	Name    string            `json:"name,omitempty"`
	Method  string            `json:"method,omitempty"`
	Path    string            `json:"path,omitempty"`
	Headers map[string]string `json:"headers,omitempty"`
}

// Service app properties.

type ServiceAppBuild struct {
	DockerImage     string            `json:"image"`
	SkipBuild       bool              `json:"skip_build"`
	SkipPull        bool              `json:"skip_pull"`
	Dockerfile      string            `json:"dockerfile"`
	DockerContext   string            `json:"context"`
	DockerBuildArgs map[string]string `json:"build_args"`
}

type ServiceAppContainer struct {
	Entrypoint *command.StringCommand `json:"entrypoint,omitempty"`
	Command    *command.StringCommand `json:"command,omitempty"`
	Port       int                    `json:"port" default:"8080"`
}

type ServiceAppProperties struct {
	Private bool `json:"private"`

	Build     *ServiceAppBuild     `json:"build,omitempty"`
	Container *ServiceAppContainer `json:"container,omitempty"`
	CDN       *AppCDN              `json:"cdn,omitempty"`
	Scheduler []*AppScheduler      `json:"scheduler,omitempty"`
}

func NewServiceAppProperties(in map[string]interface{}) (*ServiceAppProperties, error) {
	o := &ServiceAppProperties{
		Build:     &ServiceAppBuild{},
		Container: &ServiceAppContainer{},
		CDN:       &AppCDN{},
	}

	err := util.MapstructureJSONDecode(in, o)
	if err != nil {
		return nil, fmt.Errorf("error decoding service app properties: %w", err)
	}

	return o, defaults.Set(o)
}

func (p *ServiceAppProperties) Encode() (map[string]interface{}, error) {
	return encodeToMap(p)
}

type ServiceAppDeployOptions struct {
	CPULimit    float64 `json:"cpu_limit,omitempty"`
	MemoryLimit int     `json:"memory_limit,omitempty"`
	MinScale    int     `json:"min_scale,omitempty"`
	MaxScale    int     `json:"max_scale,omitempty"`
	Timeout     int     `json:"timeout,omitempty"`
}

func NewServiceAppDeployOptions(in map[string]interface{}) (*ServiceAppDeployOptions, error) {
	o := &ServiceAppDeployOptions{}

	return o, util.MapstructureJSONDecode(in, o)
}

// Static app properties.

type StaticAppBuild struct {
	Env     map[string]string      `json:"env"`
	Command *command.StringCommand `json:"command"`
	Dir     string                 `json:"dir"`
}

type StaticAppBasicAuth struct {
	Realm string            `json:"realm,omitempty"`
	Users map[string]string `json:"users,omitempty"`
}

type StaticAppProperties struct {
	Build     *StaticAppBuild     `json:"build,omitempty"`
	CDN       *AppCDN             `json:"cdn,omitempty"`
	BasicAuth *StaticAppBasicAuth `json:"basic_auth,omitempty"`

	Routing             string `json:"routing,omitempty"`
	RemoveTrailingSlash *bool  `json:"remove_trailing_slash,omitempty"`
}

func NewStaticAppProperties(in map[string]interface{}) (*StaticAppProperties, error) {
	o := &StaticAppProperties{
		Build: &StaticAppBuild{},
		CDN:   &AppCDN{},
	}

	return o, util.MapstructureJSONDecode(in, o)
}

func (p *StaticAppProperties) Encode() (map[string]interface{}, error) {
	return encodeToMap(p)
}

type StaticAppDeployOptions struct {
	MinScale int      `json:"min_scale,omitempty"`
	MaxScale int      `json:"max_scale,omitempty"`
	Timeout  int      `json:"timeout,omitempty"`
	Patterns []string `json:"patterns,omitempty"`
}

func NewStaticAppDeployOptions(in map[string]interface{}) (*StaticAppDeployOptions, error) {
	o := &StaticAppDeployOptions{}

	return o, util.MapstructureJSONDecode(in, o)
}

// Function app properties.
type FunctionAppBuild struct {
	Env     map[string]string      `json:"env"`
	Command *command.StringCommand `json:"command"`
	Dir     string                 `json:"dir"`
}

type FunctionAppProperties struct {
	Private    bool   `json:"private,omitempty"`
	Entrypoint string `json:"entrypoint,omitempty"`
	Runtime    string `json:"runtime,omitempty"`

	Build     *FunctionAppBuild `json:"build,omitempty"`
	CDN       *AppCDN           `json:"cdn,omitempty"`
	Scheduler []*AppScheduler   `json:"scheduler,omitempty"`
}

func NewFunctionAppProperties(in map[string]interface{}) (*FunctionAppProperties, error) {
	o := &FunctionAppProperties{
		Build: &FunctionAppBuild{},
		CDN:   &AppCDN{},
	}

	return o, util.MapstructureJSONDecode(in, o)
}

type FunctionAppDeployOptions struct {
	MemoryLimit int `json:"memory_limit,omitempty"`
	MinScale    int `json:"min_scale,omitempty"`
	MaxScale    int `json:"max_scale,omitempty"`
	Timeout     int `json:"timeout,omitempty"`
}

func NewFunctionAppDeployOptions(in map[string]interface{}) (*FunctionAppDeployOptions, error) {
	o := &FunctionAppDeployOptions{}

	return o, util.MapstructureJSONDecode(in, o)
}

func (p *FunctionAppProperties) Encode() (map[string]interface{}, error) {
	return encodeToMap(p)
}

type AppVars map[string]interface{}

func VarsFromAppType(app *apiv1.App) map[string]interface{} {
	return map[string]interface{}{
		"url":         app.Url,
		"cloud_url":   app.Url,
		"private_url": app.Url,
	}
}

func VarsFromAppRunType(app *apiv1.AppRun) map[string]interface{} {
	return map[string]interface{}{
		"url":         app.Url,
		"cloud_url":   app.Url,
		"private_url": app.Url,
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

func AppVarsFromAppStates(apps []*apiv1.AppState) AppVars {
	appVars := make(map[string]interface{}) // type->name->value

	for _, state := range apps {
		vars := VarsFromAppType(state.App)

		if _, ok := appVars[state.App.Type]; !ok {
			appVars[state.App.Type] = map[string]interface{}{
				state.App.Name: vars,
			}
		} else {
			appVars[state.App.Type].(map[string]interface{})[state.App.Name] = vars
		}
	}

	for _, state := range apps {
		if state.Dns == nil {
			continue
		}

		vars := appVars[state.App.Type].(map[string]interface{})[state.App.Name].(map[string]interface{})

		if state.Dns.CloudUrl != "" {
			vars["cloud_url"] = state.Dns.CloudUrl
		}

		if state.Dns.InternalUrl != "" {
			vars["private_url"] = state.Dns.InternalUrl
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

func MergeAppVars(a1, a2 AppVars) AppVars {
	return util.MergeMaps(a1, a2)
}
