package types

import "fmt"

const (
	DNSObject = "dns"
)

type AppPlan struct {
	IsDeploy bool      `json:"is_deploy"`
	IsDNS    bool      `json:"is_dns"`
	App      *App      `json:"app"`
	State    *AppState `json:"state`
}

func (a *AppPlan) String() string {
	return fmt.Sprintf("AppPlan<App=%s,IsDeploy=%t,IsDNS=%t>", a.App, a.IsDeploy, a.IsDNS)
}

type DependencyPlan struct {
	Dependency *Dependency      `json:"dependency"`
	State      *DependencyState `json:"state`
}

func (d *DependencyPlan) String() string {
	return fmt.Sprintf("DepPlan<Dep=%s>", d.Dependency)
}

type Plan struct {
	Plugin       []*PluginPlanActions     `json:"plugin,omitempty"`
	Apps         []*AppPlanActions        `json:"apps,omitempty"`
	Dependencies []*DependencyPlanActions `json:"dependencies,omitempty"`
}

type PlanOperation int

const (
	PlanAdd PlanOperation = iota + 1
	PlanUpdate
	PlanDelete
)

type PlanActionOperation struct {
	Operation PlanOperation `json:"op"`
	Data      []byte        `json:"data"`
}

type PlanAction struct {
	Description string                 `json:"description"`
	Operations  []*PlanActionOperation `json:"operations"`
}

type PluginPlanActions struct {
	Actions map[string]*PlanAction `json:"actions"`
}

type AppPlanActions struct {
	App     *App                   `json:"app"`
	Actions map[string]*PlanAction `json:"actions"`
}

func NewAppPlanActions(app *App) *AppPlanActions {
	return &AppPlanActions{
		App:     app,
		Actions: make(map[string]*PlanAction),
	}
}

type DependencyPlanActions struct {
	Dependency *Dependency            `json:"dependency"`
	Actions    map[string]*PlanAction `json:"actions"`
}
