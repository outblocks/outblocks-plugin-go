package types

import "fmt"

const (
	DNSObject = "dns"
)

type AppInfo struct {
	IsDeploy bool      `json:"is_deploy"`
	IsDNS    bool      `json:"is_dns"`
	App      *App      `json:"app"`
	State    *AppState `json:"state`
}

func (a *AppInfo) String() string {
	return fmt.Sprintf("AppInfo<App=%s,IsDeploy=%t,IsDNS=%t>", a.App, a.IsDeploy, a.IsDNS)
}

type DependencyInfo struct {
	Dependency *Dependency      `json:"dependency"`
	State      *DependencyState `json:"state`
}

func (d *DependencyInfo) String() string {
	return fmt.Sprintf("DepInfo<Dep=%s>", d.Dependency)
}

type Plan struct {
	Plugin       []*PluginPlan     `json:"plugin,omitempty"`
	Apps         []*AppPlan        `json:"apps,omitempty"`
	Dependencies []*DependencyPlan `json:"dependencies,omitempty"`
}

type PlanAction struct {
	Object      string `json:"object"`
	Description string `json:"description"`
	Data        []byte `json:"data"`
}

func (a *PlanAction) IsDNS() bool {
	return a.Object == DNSObject
}

type PluginPlan struct {
	Add    []*PlanAction `json:"add"`
	Change []*PlanAction `json:"change"`
	Remove []*PlanAction `json:"remove"`
}

type AppPlan struct {
	App    *App          `json:"app"`
	Add    []*PlanAction `json:"add"`
	Change []*PlanAction `json:"change"`
	Remove []*PlanAction `json:"remove"`
}

type DependencyPlan struct {
	Dependency *Dependency   `json:"dependency"`
	Add        []*PlanAction `json:"add"`
	Change     []*PlanAction `json:"change"`
	Remove     []*PlanAction `json:"remove"`
}
