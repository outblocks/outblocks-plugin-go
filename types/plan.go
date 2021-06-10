package types

import (
	"fmt"
)

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
	Key         string                 `json:"key"`
	Index       int                    `json:"idx"`
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
	*PlanActions
	Object string `json:"object"`
}

func NewPluginPlanActions(obj string) *PluginPlanActions {
	return &PluginPlanActions{
		PlanActions: &PlanActions{},
		Object:      obj,
	}
}

type AppPlanActions struct {
	*PlanActions
	App *App `json:"app"`
}

func NewAppPlanActions(app *App) *AppPlanActions {
	return &AppPlanActions{
		PlanActions: &PlanActions{},
		App:         app,
	}
}

type DependencyPlanActions struct {
	*PlanActions
	Dependency *Dependency `json:"dependency"`
}

func NewDependencyPlanActions(dep *Dependency) *DependencyPlanActions {
	return &DependencyPlanActions{
		PlanActions: &PlanActions{},
		Dependency:  dep,
	}
}

func NewPlanAction(typ PlanType, key, desc string, op []*PlanActionOperation) *PlanAction {
	return &PlanAction{
		Index:       -1,
		Key:         key,
		Type:        typ,
		Description: desc,
		Operations:  op,
	}
}

func NewPlanActionCreate(key, desc string, op []*PlanActionOperation) *PlanAction {
	return NewPlanAction(PlanCreate, key, desc, op)
}

func NewPlanActionRecreate(key, desc string, op []*PlanActionOperation) *PlanAction {
	return NewPlanAction(PlanRecreate, key, desc, op)
}

func NewPlanActionUpdate(key, desc string, op []*PlanActionOperation) *PlanAction {
	return NewPlanAction(PlanUpdate, key, desc, op)
}

func NewPlanActionDelete(key, desc string, op []*PlanActionOperation) *PlanAction {
	return NewPlanAction(PlanDelete, key, desc, op)
}
