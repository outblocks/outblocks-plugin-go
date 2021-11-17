package types

import (
	"fmt"
)

type AppRun struct {
	App     *App   `json:"app"`
	Command string `json:"command"`
	URL     string `json:"url"`
	IP      string `json:"ip"`
	Port    int    `json:"port"`
}

func (a *AppRun) String() string {
	return fmt.Sprintf("AppRun<App=%s,IP=%s,Port=%d>", a.App, a.IP, a.Port)
}

type DependencyRun struct {
	Dependency *Dependency `json:"dependency"`
	IP         string      `json:"ip"`
	Port       int         `json:"port"`
}

func (d *DependencyRun) String() string {
	return fmt.Sprintf("DepRun<Dep=%s>", d.Dependency)
}
