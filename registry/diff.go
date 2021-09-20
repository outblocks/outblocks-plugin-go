package registry

import (
	"context"

	"github.com/outblocks/outblocks-plugin-go/registry/fields"
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
	Object *ResourceWrapper
	Type   DiffType
	Fields []string

	applied []chan struct{}
}

func NewDiff(o *ResourceWrapper, t DiffType, fieldList []string) *Diff {
	var applied []chan struct{}

	applied = append(applied, make(chan struct{}))

	if t == DiffTypeRecreate {
		applied = append(applied, make(chan struct{}))
	}

	return &Diff{
		Object: o,
		Type:   t,
		Fields: fieldList,

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

func (d *Diff) ToPlanAction() *types.PlanAction {
	switch d.Type {
	case DiffTypeCreate:
		return types.NewPlanActionCreate(d.Object.ID, d.ObjectType(), d.Object.Resource.GetName())
	case DiffTypeUpdate:
		return types.NewPlanActionUpdate(d.Object.ID, d.ObjectType(), d.Object.Resource.GetName())
	case DiffTypeRecreate:
		return types.NewPlanActionRecreate(d.Object.ID, d.ObjectType(), d.Object.Resource.GetName())
	case DiffTypeDelete:
		return types.NewPlanActionDelete(d.Object.ID, d.ObjectType(), d.Object.Resource.GetName())
	case DiffTypeProcess:
		return types.NewPlanActionProcess(d.Object.ID, d.ObjectType(), d.Object.Resource.GetName())
	case DiffTypeNone:
		panic("unexpected diff type")
	default:
		panic("unknown diff type")
	}
}

func (d *Diff) ToApplyAction(step, total int) *types.ApplyAction {
	var typ types.PlanType

	switch d.Type {
	case DiffTypeCreate:
		typ = types.PlanCreate
	case DiffTypeUpdate:
		typ = types.PlanUpdate
	case DiffTypeDelete:
		typ = types.PlanDelete
	case DiffTypeRecreate:
		typ = types.PlanRecreate
	case DiffTypeProcess:
		typ = types.PlanProcess
	case DiffTypeNone:
		panic("unexpected diff type")
	default:
		panic("unknown diff type")
	}

	return &types.ApplyAction{
		Type:       typ,
		Namespace:  d.Object.Namespace,
		ObjectID:   d.Object.ID,
		ObjectType: d.ObjectType(),
		ObjectName: d.Object.Resource.GetName(),
		Progress:   step,
		Total:      total,
	}
}

func calculateDiff(r *ResourceWrapper, recreate bool) *Diff {
	if rdc, ok := r.Resource.(ResourceDiffCalculator); ok {
		typ := rdc.CalculateDiff()
		if typ != DiffTypeNone {
			return NewDiff(r, typ, r.FieldList())
		}

		return nil
	}

	if rbdh, ok := r.Resource.(ResourceBeforeDiffHook); ok {
		rbdh.BeforeDiff()
	}

	if r.Resource.IsNew() || recreate {
		typ := DiffTypeCreate

		if recreate && r.Resource.IsExisting() {
			typ = DiffTypeRecreate
		}

		return NewDiff(r, typ, r.FieldList())
	}

	forceNew := false

	var fieldsList []string

	for name, f := range r.Fields {
		if f.Type.Properties.Ignored {
			continue
		}

		v := f.Value.Interface().(fields.ValueTracker)
		if !v.IsChanged() {
			continue
		}

		fieldsList = append(fieldsList, name)

		if f.Type.Properties.ForceNew {
			forceNew = true
		}
	}

	if len(fieldsList) == 0 {
		return nil
	}

	typ := DiffTypeUpdate
	if forceNew {
		typ = DiffTypeRecreate
	}

	return NewDiff(r, typ, fieldsList)
}

func recreateObjectTree(r *ResourceWrapper, diffMap map[*ResourceWrapper]*Diff) {
	for d := range r.DependedBy {
		if c, ok := diffMap[d]; ok && c.Type == DiffTypeRecreate {
			continue
		}

		diffMap[d] = calculateDiff(d, true)

		if len(d.DependedBy) > 0 {
			recreateObjectTree(d, diffMap)
		}
	}
}

func deleteObjectTree(r *ResourceWrapper, diffMap map[*ResourceWrapper]*Diff) {
	for d := range r.DependedBy {
		if _, ok := d.Dependencies[r]; !ok {
			// If it's a hanging dependency (non-mutual), skip it.
			continue
		}

		deleteObjectTree(d, diffMap)
	}

	if !r.Resource.IsExisting() || r.Resource.SkipState() {
		return
	}

	diffMap[r] = NewDiff(r, DiffTypeDelete, r.FieldList())
}
