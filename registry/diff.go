package registry

import (
	"fmt"

	"github.com/outblocks/outblocks-plugin-go/registry/fields"
	"github.com/outblocks/outblocks-plugin-go/types"
)

type DiffType int

const (
	DiffTypeCreate DiffType = iota + 1
	DiffTypeUpdate
	DiffTypeRecreate
	DiffTypeDelete
)

type Diff struct {
	Object *ResourceWrapper
	Type   DiffType
	Fields []string
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
	}

	panic("unknown diff type")
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
	if r.Resource.IsNew() || recreate {
		typ := DiffTypeCreate

		if recreate && r.Resource.IsExisting() {
			typ = DiffTypeRecreate
		}

		return &Diff{
			Object: r,
			Type:   typ,
			Fields: r.FieldList(),
		}
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

	return &Diff{
		Type:   typ,
		Object: r,
		Fields: fieldsList,
	}
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
	fmt.Println(diffMap[r])

	for d := range r.DependedBy {
		deleteObjectTree(d, diffMap)
	}

	if !r.Resource.IsExisting() {
		return
	}

	d := &Diff{
		Type:   DiffTypeDelete,
		Object: r,
		Fields: r.FieldList(),
	}

	diffMap[r] = d
}
