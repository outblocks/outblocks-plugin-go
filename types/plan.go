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
	Namespace  string   `json:"namespace"`
	ObjectID   string   `json:"object_id"`
	ObjectType string   `json:"object_type"`
	ObjectName string   `json:"object_name"`
}

func NewPlanAction(typ PlanType, namespace, objectID, objectType, objectName string) *PlanAction {
	return &PlanAction{
		Type:       typ,
		Namespace:  namespace,
		ObjectID:   objectID,
		ObjectType: objectType,
		ObjectName: objectName,
	}
}

func NewPlanActionCreate(namespace, objectID, objectType, objectName string) *PlanAction {
	return NewPlanAction(PlanCreate, namespace, objectID, objectType, objectName)
}

func NewPlanActionRecreate(namespace, objectID, objectType, objectName string) *PlanAction {
	return NewPlanAction(PlanRecreate, namespace, objectID, objectType, objectName)
}

func NewPlanActionUpdate(namespace, objectID, objectType, objectName string) *PlanAction {
	return NewPlanAction(PlanUpdate, namespace, objectID, objectType, objectName)
}

func NewPlanActionDelete(namespace, objectID, objectType, objectName string) *PlanAction {
	return NewPlanAction(PlanDelete, namespace, objectID, objectType, objectName)
}

func NewPlanActionProcess(namespace, objectID, objectType, objectName string) *PlanAction {
	return NewPlanAction(PlanProcess, namespace, objectID, objectType, objectName)
}
