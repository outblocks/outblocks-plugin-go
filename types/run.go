package types

import (
	"fmt"
	"strings"

	"github.com/outblocks/outblocks-plugin-go/util"
)

type AppRun struct {
	App        *App                   `json:"app"`
	Command    string                 `json:"command"`
	Path       string                 `json:"path"`
	Env        map[string]string      `json:"env"`
	URL        string                 `json:"url"`
	IP         string                 `json:"ip"`
	Port       int                    `json:"port"`
	Properties map[string]interface{} `json:"properties"`
}

func (a *AppRun) String() string {
	return fmt.Sprintf("AppRun<App=%s,IP=%s,Port=%d>", a.App, a.IP, a.Port)
}

func (a *AppRun) EnvPrefix() string {
	return fmt.Sprintf("APP_%s_%s", strings.ToUpper(a.App.Type), util.SanitizeEnvVar(strings.ToUpper(a.App.Name)))
}

type DependencyRun struct {
	Dependency *Dependency            `json:"dependency"`
	Env        map[string]string      `json:"env"`
	IP         string                 `json:"ip"`
	Port       int                    `json:"port"`
	Properties map[string]interface{} `json:"properties"`
}

func (d *DependencyRun) String() string {
	return fmt.Sprintf("DepRun<Dep=%s>", d.Dependency)
}

func (d *DependencyRun) EnvPrefix() string {
	return fmt.Sprintf("DEP_%s_%s", strings.ToUpper(d.Dependency.Type), util.SanitizeEnvVar((strings.ToUpper(d.Dependency.Name))))
}
