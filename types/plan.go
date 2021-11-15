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
	App  *AppState `json:"app"`
	Skip bool      `json:"skip"`

	IsDeploy bool `json:"is_deploy"`
	IsDNS    bool `json:"is_dns"`
}

func (a *AppPlan) String() string {
	return fmt.Sprintf("AppPlan<App=%s,IsDeploy=%t,IsDNS=%t>", a.App, a.IsDeploy, a.IsDNS)
}

type DependencyPlan struct {
	Dependency *DependencyState `json:"dependency"`
	Skip       bool             `json:"skip"`
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
	Critical   bool     `json:"critical"`
}

func NewPlanAction(typ PlanType, source, namespace, objectID, objectType, objectName string, critical bool) *PlanAction {
	return &PlanAction{
		Type:       typ,
		Source:     source,
		Namespace:  namespace,
		ObjectID:   objectID,
		ObjectType: objectType,
		ObjectName: objectName,
		Critical:   critical,
	}
}

func NewPlanActionCreate(source, namespace, objectID, objectType, objectName string, critical bool) *PlanAction {
	return NewPlanAction(PlanCreate, source, namespace, objectID, objectType, objectName, critical)
}

func NewPlanActionRecreate(source, namespace, objectID, objectType, objectName string, critical bool) *PlanAction {
	return NewPlanAction(PlanRecreate, source, namespace, objectID, objectType, objectName, critical)
}

func NewPlanActionUpdate(source, namespace, objectID, objectType, objectName string, critical bool) *PlanAction {
	return NewPlanAction(PlanUpdate, source, namespace, objectID, objectType, objectName, critical)
}

func NewPlanActionDelete(source, namespace, objectID, objectType, objectName string, critical bool) *PlanAction {
	return NewPlanAction(PlanDelete, source, namespace, objectID, objectType, objectName, critical)
}

func NewPlanActionProcess(source, namespace, objectID, objectType, objectName string, critical bool) *PlanAction {
	return NewPlanAction(PlanProcess, source, namespace, objectID, objectType, objectName, critical)
}
