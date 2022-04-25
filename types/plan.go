package types

import apiv1 "github.com/outblocks/outblocks-plugin-go/gen/api/v1"

func NewPlanAction(typ apiv1.PlanType, source, namespace, objectID, objectType, objectName string, fields []string, critical bool) *apiv1.PlanAction {
	return &apiv1.PlanAction{
		Type:       typ,
		Source:     source,
		Namespace:  namespace,
		ObjectId:   objectID,
		ObjectType: objectType,
		ObjectName: objectName,
		Critical:   critical,
	}
}

func NewPlanActionCreate(source, namespace, objectID, objectType, objectName string, fields []string, critical bool) *apiv1.PlanAction {
	return NewPlanAction(apiv1.PlanType_PLAN_TYPE_CREATE, source, namespace, objectID, objectType, objectName, fields, critical)
}

func NewPlanActionRecreate(source, namespace, objectID, objectType, objectName string, fields []string, critical bool) *apiv1.PlanAction {
	return NewPlanAction(apiv1.PlanType_PLAN_TYPE_RECREATE, source, namespace, objectID, objectType, objectName, fields, critical)
}

func NewPlanActionUpdate(source, namespace, objectID, objectType, objectName string, fields []string, critical bool) *apiv1.PlanAction {
	return NewPlanAction(apiv1.PlanType_PLAN_TYPE_UPDATE, source, namespace, objectID, objectType, objectName, fields, critical)
}

func NewPlanActionDelete(source, namespace, objectID, objectType, objectName string, fields []string, critical bool) *apiv1.PlanAction {
	return NewPlanAction(apiv1.PlanType_PLAN_TYPE_DELETE, source, namespace, objectID, objectType, objectName, fields, critical)
}

func NewPlanActionProcess(source, namespace, objectID, objectType, objectName string, fields []string, critical bool) *apiv1.PlanAction {
	return NewPlanAction(apiv1.PlanType_PLAN_TYPE_PROCESS, source, namespace, objectID, objectType, objectName, fields, critical)
}
