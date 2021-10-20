package types

import (
	"fmt"
)

const (
	DNSObject = "dns"
)

type AppPlan struct {
	IsDeploy   bool                   `json:"is_deploy"`
	IsDNS      bool                   `json:"is_dns"`
	App        *App                   `json:"app"`
	Env        map[string]string      `json:"env"`
	Properties map[string]interface{} `json:"properties"`
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

type PlanType int

const (
	PlanCreate PlanType = iota + 1
	PlanRecreate
	PlanUpdate
	PlanDelete
	PlanProcess
)

type PlanAction struct {
	Type       PlanType `json:"type"`
	ObjectID   string   `json:"object_id"`
	ObjectType string   `json:"object_type"`
	ObjectName string   `json:"object_name"`
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

func NewPlanAction(typ PlanType, objectID, objectType, objectName string) *PlanAction {
	return &PlanAction{
		Type:       typ,
		ObjectID:   objectID,
		ObjectType: objectType,
		ObjectName: objectName,
	}
}

func NewPlanActionCreate(objectID, objectType, objectName string) *PlanAction {
	return NewPlanAction(PlanCreate, objectID, objectType, objectName)
}

func NewPlanActionRecreate(objectID, objectType, objectName string) *PlanAction {
	return NewPlanAction(PlanRecreate, objectID, objectType, objectName)
}

func NewPlanActionUpdate(objectID, objectType, objectName string) *PlanAction {
	return NewPlanAction(PlanUpdate, objectID, objectType, objectName)
}

func NewPlanActionDelete(objectID, objectType, objectName string) *PlanAction {
	return NewPlanAction(PlanDelete, objectID, objectType, objectName)
}

func NewPlanActionProcess(objectID, objectType, objectName string) *PlanAction {
	return NewPlanAction(PlanProcess, objectID, objectType, objectName)
}

type PlanActions struct {
	Actions []*PlanAction `json:"actions"`
}
