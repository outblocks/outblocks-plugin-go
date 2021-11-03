package types

import (
	"fmt"
)

const (
	DNSObject = "dns"
)

const (
	SourceApp        = "app"
	SourceDependency = "dependency"
	SourcePlugin     = "plugin"
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
	Actions []*PlanAction `json:"actions"`
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
	Source     string   `json:"source"`
	Namespace  string   `json:"namespace"`
	ObjectID   string   `json:"object_id"`
	ObjectType string   `json:"object_type"`
	ObjectName string   `json:"object_name"`
}

func NewPlanAction(typ PlanType, source, namespace, objectID, objectType, objectName string) *PlanAction {
	return &PlanAction{
		Type:       typ,
		Source:     source,
		Namespace:  namespace,
		ObjectID:   objectID,
		ObjectType: objectType,
		ObjectName: objectName,
	}
}

func NewPlanActionCreate(source, namespace, objectID, objectType, objectName string) *PlanAction {
	return NewPlanAction(PlanCreate, source, namespace, objectID, objectType, objectName)
}

func NewPlanActionRecreate(source, namespace, objectID, objectType, objectName string) *PlanAction {
	return NewPlanAction(PlanRecreate, source, namespace, objectID, objectType, objectName)
}

func NewPlanActionUpdate(source, namespace, objectID, objectType, objectName string) *PlanAction {
	return NewPlanAction(PlanUpdate, source, namespace, objectID, objectType, objectName)
}

func NewPlanActionDelete(source, namespace, objectID, objectType, objectName string) *PlanAction {
	return NewPlanAction(PlanDelete, source, namespace, objectID, objectType, objectName)
}

func NewPlanActionProcess(source, namespace, objectID, objectType, objectName string) *PlanAction {
	return NewPlanAction(PlanProcess, source, namespace, objectID, objectType, objectName)
}
