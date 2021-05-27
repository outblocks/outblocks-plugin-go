package types

import "fmt"

const (
	DNSObject = "dns"
)

type AppPlan struct {
	IsDeploy bool `json:"is_deploy"`
	IsDNS    bool `json:"is_dns"`
	App      *App `json:"app"`
}

func (a *AppPlan) String() string {
	return fmt.Sprintf("AppPlan<App=%s,IsDeploy=%t,IsDNS=%t>", a.App, a.IsDeploy, a.IsDNS)
}

type DependencyPlan struct {
	Dependency *Dependency `json:"dependency"`
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
	PlanOpAdd PlanOperation = iota + 1
	PlanOpUpdate
	PlanOpDelete
)

type PlanType int

const (
	PlanCreate PlanType = iota + 1
	PlanRecreate
	PlanUpdate
	PlanDelete
)

type PlanActionOperation struct {
	Operation PlanOperation `json:"op"`
	Steps     int           `json:"steps"`
	Data      []byte        `json:"data"`
}

type PlanAction struct {
	Type        PlanType               `json:"type"`
	Description string                 `json:"description"`
	Operations  []*PlanActionOperation `json:"operations"`
}

func (a *PlanAction) TotalSteps() int {
	steps := 0

	for _, op := range a.Operations {
		steps += op.Steps
	}

	return steps
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

func NewPlanAction(typ PlanType, desc string, op []*PlanActionOperation) *PlanAction {
	return &PlanAction{
		Type:        typ,
		Description: desc,
		Operations:  op,
	}
}

func NewPlanActionCreate(desc string, op []*PlanActionOperation) *PlanAction {
	return NewPlanAction(PlanCreate, desc, op)
}

func NewPlanActionRecreate(desc string, op []*PlanActionOperation) *PlanAction {
	return NewPlanAction(PlanRecreate, desc, op)
}

func NewPlanActionUpdate(desc string, op []*PlanActionOperation) *PlanAction {
	return NewPlanAction(PlanUpdate, desc, op)
}

func NewPlanActionDelete(desc string, op []*PlanActionOperation) *PlanAction {
	return NewPlanAction(PlanDelete, desc, op)
}
