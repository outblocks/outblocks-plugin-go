package types

import (
	"context"
	"reflect"
	"strings"
	"sync"
)

type idInfo struct {
	idx int
	obj interface{}
}

func getIDsThroughReflection(lst interface{}) map[string]idInfo {
	ids := make(map[string]idInfo)

	k := reflect.TypeOf(lst).Kind()
	if k != reflect.Slice && k != reflect.Array {
		panic("lst is not a slice/array")
	}

	v := reflect.ValueOf(lst)

	for i := 0; i < v.Len(); i++ {
		rv := v.Index(i).Interface()
		v := idInfo{
			idx: i,
			obj: rv,
		}

		ids[rv.(Indexed).Key()] = v
	}

	return ids
}

type PlanActions struct {
	Actions []*PlanAction `json:"actions"`
}

func (p *PlanActions) PlanObject(ctx context.Context, o interface{}, fieldName string, dest interface{}, verify bool) error {
	structVal := reflect.ValueOf(o).Elem()
	fieldVal := structVal.FieldByName(fieldName)
	key := findPlanKey(o, fieldName)

	action, err := fieldVal.Interface().(Planner).Plan(ctx, key, dest, verify)
	if err != nil {
		return err
	}

	if action != nil {
		p.Actions = append(p.Actions, action)
	}

	return nil
}

func getStructFieldPlanKey(f *reflect.StructField) string {
	for _, tag := range []string{"plan", "json", "mapstructure"} {
		v := f.Tag.Get(tag)
		v = strings.Split(v, ",")[0]

		if v == "" {
			continue
		}

		return strings.ReplaceAll(v, "_", " ")
	}

	return ""
}

func findPlanKey(o interface{}, fieldName string) string {
	structType := reflect.TypeOf(o).Elem()
	field, ok := structType.FieldByName(fieldName)

	if !ok {
		return ""
	}

	return getStructFieldPlanKey(&field)
}

func (p *PlanActions) PlanObjectList(ctx context.Context, o interface{}, fieldName string, dest interface{}, verify bool) error {
	structVal := reflect.ValueOf(o).Elem()
	fieldVal := structVal.FieldByName(fieldName)
	lst := fieldVal.Interface()
	key := findPlanKey(o, fieldName)

	curIDs := getIDsThroughReflection(lst)
	destIDs := getIDsThroughReflection(dest)

	elemType := reflect.TypeOf(lst).Elem()
	if elemType.Kind() == reflect.Ptr {
		elemType = elemType.Elem()
	}

	val := fieldVal

	var (
		action *PlanAction
		err    error
	)

	for k, v := range destIDs {
		cur, ok := curIDs[k]
		idx := -1

		if !ok {
			// Creations.
			n := reflect.New(elemType)
			action, err = n.Interface().(Planner).Plan(ctx, key, v.obj, verify)

			if action == nil || action.Type != PlanCreate {
				val = reflect.Append(val, n)
			}
		} else {
			// Updates.
			idx = cur.idx
			action, err = cur.obj.(Planner).Plan(ctx, key, v.obj, verify)

			delete(curIDs, k)
		}

		if err != nil {
			return err
		}

		if action != nil {
			action.Index = idx
			p.Actions = append(p.Actions, action)
		}
	}

	// Deletions.
	for _, v := range curIDs {
		action, err = v.obj.(Planner).Plan(ctx, key, nil, verify)

		if err != nil {
			return err
		}

		if action != nil {
			action.Index = v.idx
			p.Actions = append(p.Actions, action)
		}
	}

	fieldVal.Set(val)

	return nil
}

func (p *PlanActions) Apply(ctx context.Context, o interface{}, callback ApplyActionCallbackFunc) error {
	structType := reflect.TypeOf(o).Elem()
	structVal := reflect.ValueOf(o).Elem()

	keyToField := make(map[string]reflect.Value, structType.NumField())

	for i := 0; i < structType.NumField(); i++ {
		f := structType.Field(i)
		keyToField[getStructFieldPlanKey(&f)] = structVal.Field(i)
	}

	var (
		first  []*PlanAction // deletions in reverse
		second []*PlanAction // rest
	)

	for _, act := range p.Actions {
		var (
			opFirst  []*PlanActionOperation
			opSecond []*PlanActionOperation
		)

		for _, op := range act.Operations {
			if op.Operation == PlanOpDelete {
				opFirst = append(opFirst, op)
			} else {
				opSecond = append(opSecond, op)
			}
		}

		if len(opFirst) != 0 {
			actCpy := *act
			actCpy.Operations = opFirst
			first = append(first, &actCpy)
		}

		if len(opSecond) != 0 {
			actCpy := *act
			actCpy.Operations = opSecond
			second = append(second, &actCpy)
		}
	}

	// Reverse order of deletions.
	for i, j := 0, len(first)-1; i < j; i, j = i+1, j-1 {
		first[i], first[j] = first[j], first[i]
	}

	err := applyActions(ctx, callback, keyToField, first)
	if err != nil {
		return err
	}

	return applyActions(ctx, callback, keyToField, second)
}

func applyActions(ctx context.Context, callback ApplyActionCallbackFunc, keyToField map[string]reflect.Value, actions []*PlanAction) error {
	for _, act := range actions {
		field := keyToField[act.Key]
		kind := field.Type().Kind()

		if kind == reflect.Slice || kind == reflect.Array {
			newV, err := applyObjectList(ctx, field.Interface(), act, act.Key, callback)
			if err != nil {
				return err
			}

			field.Set(reflect.ValueOf(newV))
		} else {
			err := applyObject(ctx, field.Interface().(Planner), act, act.Key, callback)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func applyObject(ctx context.Context, o Planner, act *PlanAction, obj string, callback ApplyActionCallbackFunc) error {
	progress := 0
	total := act.TotalSteps()

	if total == 0 {
		return nil
	}

	var mu sync.Mutex

	callback(obj, "start", act.Index, 0, total)

	cb := func(desc string) {
		mu.Lock()
		progress++
		callback(obj, desc, act.Index, progress, total)
		mu.Unlock()
	}

	err := o.Apply(ctx, act.Operations, cb)

	return err
}

func applyObjectList(ctx context.Context, lst interface{}, act *PlanAction, obj string, callback ApplyActionCallbackFunc) (interface{}, error) {
	typ := reflect.TypeOf(lst)
	kind := typ.Kind()

	if kind != reflect.Slice && kind != reflect.Array {
		panic("lst is not a slice/array")
	}

	val := reflect.ValueOf(lst)
	elemType := typ.Elem()

	if elemType.Kind() == reflect.Ptr {
		elemType = elemType.Elem()
	}

	if len(act.Operations) == 0 {
		return lst, nil
	}

	todel := make(map[int]struct{})

	progress := 0
	total := act.TotalSteps()

	if total == 0 {
		return lst, nil
	}

	var mu sync.Mutex

	callback(obj, "start", act.Index, 0, total)

	cb := func(desc string) {
		mu.Lock()
		progress++
		callback(obj, desc, act.Index, progress, total)
		mu.Unlock()
	}

	var err error

	if act.Index >= 0 {
		// Updates or deletions.
		err = val.Index(act.Index).Interface().(Planner).Apply(ctx, act.Operations, cb)

		if act.Type == PlanDelete {
			todel[act.Index] = struct{}{}
		}
	} else {
		// Creations.
		n := reflect.New(elemType)
		err = n.Interface().(Planner).Apply(ctx, act.Operations, cb)

		val = reflect.Append(val, n)
	}

	if err != nil {
		return nil, err
	}

	// Process deletions.
	if len(todel) != 0 {
		newLst := reflect.MakeSlice(reflect.SliceOf(typ.Elem()), 0, val.Len()-len(todel))

		for i := 0; i < val.Len(); i++ {
			if _, ok := todel[i]; !ok {
				newLst = reflect.Append(newLst, val.Index(i))
			}
		}

		return newLst.Interface(), nil
	}

	return val.Interface(), nil
}
