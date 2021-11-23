package registry

import (
	"context"

	apiv1 "github.com/outblocks/outblocks-plugin-go/gen/api/v1"
	"github.com/outblocks/outblocks-plugin-go/types"
)

type DiffType int

const (
	DiffTypeNone DiffType = iota + 1
	DiffTypeCreate
	DiffTypeUpdate
	DiffTypeRecreate
	DiffTypeDelete
	DiffTypeProcess
)

type Diff struct {
	Object   *ResourceWrapper
	Type     DiffType
	Fields   []string
	Critical bool

	applied []chan struct{}
}

func NewDiff(o *ResourceWrapper, t DiffType, fieldList []string) *Diff {
	var applied []chan struct{}

	applied = append(applied, make(chan struct{}))

	if t == DiffTypeRecreate {
		applied = append(applied, make(chan struct{}))
	}

	var critical bool

	if rcc, ok := o.Resource.(ResourceCriticalChecker); ok {
		critical = rcc.IsCritical(t, fieldList)
	}

	return &Diff{
		Object:   o,
		Type:     t,
		Fields:   fieldList,
		Critical: critical,

		applied: applied,
	}
}

func (d *Diff) Wait(steps int) {
	for i := 0; (i < steps || steps <= 0) && i < len(d.applied); i++ {
		<-d.applied[i]
	}
}

func (d *Diff) WaitContext(ctx context.Context, steps int) error {
	for i := 0; (i < steps || steps <= 0) && i < len(d.applied); i++ {
		select {
		case <-d.applied[i]:
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return nil
}

func (d *Diff) IsApplied(steps int) bool {
	for i := 0; (i < steps || steps <= 0) && i < len(d.applied); i++ {
		select {
		case <-d.applied[i]:
			return true
		default:
			return false
		}
	}

	return false
}

func (d *Diff) MarkStepAsApplied() {
	d.SetApplied(d.AppliedSteps())
}

func (d *Diff) SetApplied(step int) {
	close(d.applied[step])
}

func (d *Diff) AppliedSteps() int {
	for i := 0; i < len(d.applied); i++ {
		select {
		case <-d.applied[i]:
		default:
			return i
		}
	}

	return 0
}

func (d *Diff) RequiredSteps() int {
	return len(d.applied)
}

func (d *Diff) Applied() bool {
	return d.AppliedSteps() == d.RequiredSteps()
}

func (d *Diff) ObjectType() string {
	typ := d.Object.Type

	v, ok := d.Object.Resource.(ResourceTypeVerbose)
	if ok {
		typ = v.GetType()
	}

	return typ
}

func (d *Diff) ToPlanAction() *apiv1.PlanAction {
	switch d.Type {
	case DiffTypeCreate:
		return types.NewPlanActionCreate(d.Object.Source, d.Object.Namespace, d.Object.ID, d.ObjectType(), d.Object.Resource.GetName(), d.Critical)
	case DiffTypeUpdate:
		return types.NewPlanActionUpdate(d.Object.Source, d.Object.Namespace, d.Object.ID, d.ObjectType(), d.Object.Resource.GetName(), d.Critical)
	case DiffTypeRecreate:
		return types.NewPlanActionRecreate(d.Object.Source, d.Object.Namespace, d.Object.ID, d.ObjectType(), d.Object.Resource.GetName(), d.Critical)
	case DiffTypeDelete:
		return types.NewPlanActionDelete(d.Object.Source, d.Object.Namespace, d.Object.ID, d.ObjectType(), d.Object.Resource.GetName(), d.Critical)
	case DiffTypeProcess:
		return types.NewPlanActionProcess(d.Object.Source, d.Object.Namespace, d.Object.ID, d.ObjectType(), d.Object.Resource.GetName(), d.Critical)
	case DiffTypeNone:
		panic("unexpected diff type")
	default:
		panic("unknown diff type")
	}
}

func (d *Diff) ToApplyAction(step, total int) *apiv1.ApplyAction {
	var typ apiv1.PlanType

	switch d.Type {
	case DiffTypeCreate:
		typ = apiv1.PlanType_PLAN_TYPE_CREATE
	case DiffTypeUpdate:
		typ = apiv1.PlanType_PLAN_TYPE_UPDATE
	case DiffTypeDelete:
		typ = apiv1.PlanType_PLAN_TYPE_DELETE
	case DiffTypeRecreate:
		typ = apiv1.PlanType_PLAN_TYPE_RECREATE
	case DiffTypeProcess:
		typ = apiv1.PlanType_PLAN_TYPE_PROCESS
	case DiffTypeNone:
		panic("unexpected diff type")
	default:
		panic("unknown diff type")
	}

	return &apiv1.ApplyAction{
		Type:       typ,
		Source:     d.Object.Source,
		Namespace:  d.Object.Namespace,
		ObjectId:   d.Object.ID,
		ObjectType: d.ObjectType(),
		ObjectName: d.Object.Resource.GetName(),
		Progress:   int32(step),
		Total:      int32(total),
	}
}

func deleteObjectTree(r *ResourceWrapper, diffMap map[*ResourceWrapper]*Diff, onlyUnregistered bool) {
	for d := range r.DependedBy {
		deleteObjectTree(d, diffMap, onlyUnregistered)
	}

	if !r.Resource.IsExisting() || r.Resource.SkipState() || (onlyUnregistered && r.IsRegistered) {
		return
	}

	diff := NewDiff(r, DiffTypeDelete, r.FieldList())
	diffMap[r] = diff

	r.Resource.SetDiff(diff)
}
