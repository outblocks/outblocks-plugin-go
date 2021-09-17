package types

import (
	"fmt"
)

type AppRun struct {
	App        *App                   `json:"app"`
	Command    string                 `json:"command"`
	Env        map[string]string      `json:"env"`
	URL        string                 `json:"url"`
	IP         string                 `json:"ip"`
	Port       int                    `json:"port"`
	Properties map[string]interface{} `json:"properties"`
}

func (a *AppRun) String() string {
	return fmt.Sprintf("AppRun<App=%s,IP=%s,Port=%d>", a.App, a.IP, a.Port)
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
