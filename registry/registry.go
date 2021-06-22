package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/outblocks/outblocks-plugin-go/registry/fields"
	"github.com/outblocks/outblocks-plugin-go/types"
	"github.com/outblocks/outblocks-plugin-go/util/errgroup"
)

const defaultConcurrency = 5

type ResourceTypeInfo struct {
	Type   reflect.Type
	Fields map[string]*FieldTypeInfo
}

type Registry struct {
	types    map[string]*ResourceTypeInfo
	fieldMap map[interface{}]*ResourceWrapper

	resources   []*ResourceWrapper
	resourceMap map[ResourceID]*ResourceWrapper
	missing     []*ResourceWrapper
}

func NewRegistry() *Registry {
	return &Registry{
		fieldMap:    make(map[interface{}]*ResourceWrapper),
		resourceMap: make(map[ResourceID]*ResourceWrapper),
		types:       make(map[string]*ResourceTypeInfo),
	}
}

func (r *Registry) RegisterType(o Resource) error {
	t := reflect.TypeOf(o)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if _, ok := r.types[t.Name()]; ok {
		return nil
	}

	fieldsMap := make(map[string]*FieldTypeInfo)

	for i := 0; i < t.NumField(); i++ {
		ft := t.Field(i)

		if ft.Anonymous {
			continue
		}

		def, defSet := ft.Tag.Lookup("default")

		fieldsMap[ft.Name] = &FieldTypeInfo{
			ReflectType: ft,
			Properties:  parseFieldPropertiesTag(ft.Tag.Get("state")),
			Default:     def,
			DefaultSet:  defSet,
		}
	}

	r.types[t.Name()] = &ResourceTypeInfo{
		Type:   t,
		Fields: fieldsMap,
	}

	return nil
}

func generateResourceFields(o Resource, rti *ResourceTypeInfo) map[string]*FieldInfo {
	v := reflect.ValueOf(o)
	t := rti.Type

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	fieldsMap := make(map[string]*FieldInfo)

	for i := 0; i < t.NumField(); i++ {
		ft := t.Field(i)

		if ft.Anonymous {
			continue
		}

		fieldsMap[ft.Name] = &FieldInfo{
			Type:  rti.Fields[ft.Name],
			Value: v.Field(i),
		}
	}

	return fieldsMap
}

func (r *Registry) Register(o Resource, namespace, id string) error {
	t := reflect.TypeOf(o)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	tinfo, ok := r.types[t.Name()]

	if !ok {
		err := r.RegisterType(o)
		if err != nil {
			return err
		}

		tinfo = r.types[t.Name()]
	}

	resourceID := ResourceID{
		ID:        id,
		Namespace: namespace,
		Type:      t.Name(),
	}

	if erw, ok := r.resourceMap[resourceID]; ok {
		reflect.ValueOf(o).Elem().Set(reflect.ValueOf(erw.Resource).Elem())

		return nil
	}

	o.SetNew(true)

	rw := &ResourceWrapper{
		ResourceID: resourceID,
		Fields:     generateResourceFields(o, tinfo),
		Resource:   o,
	}
	r.resourceMap[resourceID] = rw

	err := setFieldDefaults(rw)
	if err != nil {
		return err
	}

	for _, f := range rw.Fields {
		if f.Type.Properties.Ignored {
			continue
		}

		if erw, ok := r.fieldMap[f.Value.Interface()]; ok {
			rw.Dependencies = append(rw.Dependencies, erw)
			erw.DependedBy = append(erw.DependedBy, rw)

			continue
		}

		if fh, ok := f.Value.Interface().(fields.FieldHolder); ok {
			for _, dep := range fh.FieldDependencies() {
				if erw, ok := r.fieldMap[dep]; ok {
					rw.Dependencies = append(rw.Dependencies, erw)
					erw.DependedBy = append(erw.DependedBy, rw)
				}
			}
		}

		r.fieldMap[f.Value.Interface()] = rw
	}

	r.resources = append(r.resources, rw)

	return nil
}

func (r *Registry) Load(ctx context.Context, state []byte) error {
	if len(state) == 0 {
		return nil
	}

	var existing []*ResourceData

	err := json.Unmarshal(state, &existing)
	if err != nil {
		return err
	}

	existingMap := make(map[ResourceID]*ResourceData)
	resourceMap := make(map[ResourceID]*ResourceWrapper)

	for _, e := range existing {
		existingMap[e.ResourceID] = e
	}

	for _, res := range r.resources {
		resourceMap[res.ResourceID] = res

		e, ok := existingMap[res.ResourceID]
		if !ok {
			continue
		}

		delete(existingMap, res.ResourceID)

		err := res.SetFieldValues(e.Properties)
		if err != nil {
			return err
		}

		res.Resource.SetNew(false)
	}

	// Process missing resources.
	var missing []*ResourceWrapper

	for _, v := range existingMap {
		rti, ok := r.types[v.Type]
		if !ok {
			return fmt.Errorf("unknown resource type found: %s", v.Type)
		}

		obj := reflect.New(rti.Type)
		res := obj.Interface().(Resource)

		rw := &ResourceWrapper{
			ResourceID: v.ResourceID,
			Fields:     generateResourceFields(res, rti),
			Resource:   obj.Interface().(Resource),
		}

		err := setFieldDefaults(rw)
		if err != nil {
			return err
		}

		err = rw.SetFieldValues(v.Properties)
		if err != nil {
			return err
		}

		resourceMap[v.ResourceID] = rw
		missing = append(missing, rw)
	}

	// Fill dependencies now that resourceMap is filled.
	for _, v := range existingMap {
		rw := resourceMap[v.ResourceID]

		for _, d := range v.DependedBy {
			dep, ok := resourceMap[d]
			if !ok {
				return fmt.Errorf("dependency missing: %s", d)
			}

			rw.DependedBy = append(rw.DependedBy, dep)
		}

		for _, d := range v.Dependencies {
			dep, ok := resourceMap[d]
			if !ok {
				return fmt.Errorf("dependency missing: %s", d)
			}

			rw.Dependencies = append(rw.Dependencies, dep)
		}
	}

	r.missing = missing

	return nil
}

func (r *Registry) Read(ctx context.Context) error {
	g, _ := errgroup.WithConcurrency(ctx, defaultConcurrency)

	for _, o := range r.resources {
		o := o

		g.Go(func() error {
			return o.Resource.Read(ctx)
		})
	}

	return g.Wait()
}

var (
	stringInputType  = reflect.TypeOf((*fields.StringInputField)(nil)).Elem()
	stringOutputType = reflect.TypeOf((*fields.StringOutputField)(nil)).Elem()
	boolInputType    = reflect.TypeOf((*fields.BoolInputField)(nil)).Elem()
	boolOutputType   = reflect.TypeOf((*fields.BoolOutputField)(nil)).Elem()
	intInputType     = reflect.TypeOf((*fields.IntInputField)(nil)).Elem()
	intOutputType    = reflect.TypeOf((*fields.IntOutputField)(nil)).Elem()
	mapInputType     = reflect.TypeOf((*fields.MapInputField)(nil)).Elem()
	mapOutputType    = reflect.TypeOf((*fields.MapOutputField)(nil)).Elem()
)

func setFieldDefaults(r *ResourceWrapper) error {
	for _, f := range r.Fields {
		if f.Type.Properties.Ignored || !f.Value.IsNil() {
			continue
		}

		defaultTag := f.Type.Default
		ok := f.Type.DefaultSet

		var val interface{}

		switch f.Type.ReflectType.Type {
		// String.
		case stringInputType:
			if ok {
				val = fields.String(defaultTag)
			} else {
				val = fields.StringUnset()
			}
		case stringOutputType:
			val = fields.StringUnsetOutput()

			// Bool.
		case boolInputType:
			if ok {
				val = fields.Bool(defaultTag == "1" || defaultTag == "true")
			} else {
				val = fields.BoolUnset()
			}
		case boolOutputType:
			val = fields.StringUnsetOutput()

			// Int.
		case intInputType:
			if ok {
				v, _ := strconv.Atoi(defaultTag)
				val = fields.Int(v)
			} else {
				val = fields.IntUnset()
			}
		case intOutputType:
			val = fields.IntUnsetOutput()

			// Map.
		case mapInputType:
			val = fields.MapUnset()
		case mapOutputType:
			val = fields.MapUnsetOutput()

		default:
			return fmt.Errorf("unknown field type %s", f.Type.ReflectType.Type)
		}

		f.Value.Set(reflect.ValueOf(val))
	}

	return nil
}

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

type FieldProperties struct {
	Ignored     bool
	ForceWanted bool
}

func parseFieldPropertiesTag(tag string) *FieldProperties {
	ret := &FieldProperties{}
	taginfo := strings.Split(tag, ",")

	for _, t := range taginfo {
		switch t {
		case "-":
			ret.Ignored = true

		case "force_wanted":
			ret.ForceWanted = true
		}
	}

	return ret
}

func calculateDiff(r *ResourceWrapper, recreate bool) *Diff {
	if r.Resource.IsNew() || recreate {
		typ := DiffTypeCreate

		if recreate {
			typ = DiffTypeRecreate
		}

		return &Diff{
			Object: r,
			Type:   typ,
			Fields: r.FieldList(),
		}
	}

	forceWanted := false

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

		if f.Type.Properties.ForceWanted {
			forceWanted = true
		}
	}

	if len(fieldsList) == 0 {
		return nil
	}

	typ := DiffTypeUpdate
	if forceWanted {
		typ = DiffTypeRecreate
	}

	return &Diff{
		Type:   typ,
		Object: r,
		Fields: fieldsList,
	}
}

func recreateObjectTree(r *ResourceWrapper, diffMap map[*ResourceWrapper]*Diff) {
	for _, d := range r.DependedBy {
		if _, ok := diffMap[d]; ok {
			continue
		}

		diffMap[d] = calculateDiff(r, true)

		recreateObjectTree(d, diffMap)
	}
}

func deleteObjectTree(r *ResourceWrapper, diffMap map[*ResourceWrapper]*Diff) {
	for _, d := range r.DependedBy {
		deleteObjectTree(d, diffMap)
	}

	d := &Diff{
		Type:   DiffTypeDelete,
		Object: r,
		Fields: r.FieldList(),
	}

	diffMap[r] = d
}

func isTreeMarkedForDeletion(r *ResourceWrapper, diffMap map[*ResourceWrapper]*Diff) bool {
	for _, d := range r.DependedBy {
		v, ok := diffMap[d]
		if !ok {
			return false
		}

		if v.Type != DiffTypeDelete {
			return false
		}
	}

	return true
}

func (r *Registry) Diff(ctx context.Context) ([]*Diff, error) {
	diffMap := make(map[*ResourceWrapper]*Diff)

	for _, o := range r.resources {
		d := calculateDiff(o, false)
		if d != nil {
			diffMap[o] = d

			if d.Type == DiffTypeRecreate {
				recreateObjectTree(d.Object, diffMap)
			}
		}
	}

	// Add all missing resources as deletions.
	for _, o := range r.missing {
		deleteObjectTree(o, diffMap)
	}

	// Verify if all that depends on it are also deleted.
	for _, d := range diffMap {
		if d.Type == DiffTypeDelete && !isTreeMarkedForDeletion(d.Object, diffMap) {
			return nil, fmt.Errorf("object: '%s' marked for deletion but there are still objects depending on it", d.Object.ResourceID)
		}
	}

	var diff []*Diff
	for _, v := range diffMap {
		diff = append(diff, v)
	}

	return diff, nil
}

func (r *Registry) Dump() ([]byte, error) {
	var resources []*ResourceWrapper
	for _, res := range r.resources {
		if res.Resource.IsNew() {
			continue
		}

		resources = append(resources, res)
	}

	return json.Marshal(resources)
}

type diffActionType struct {
	rw     *ResourceWrapper
	delete bool
}

type diffAction struct {
	diff *Diff
	wg   sync.WaitGroup
}

func waitContext(ctx context.Context, wg *sync.WaitGroup) error {
	c := make(chan struct{})

	go func() {
		defer close(c)
		wg.Wait()
	}()

	select {
	case <-c:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func prepareDiffActionMap(diff []*Diff) map[diffActionType]*diffAction {
	diffActionMap := make(map[diffActionType]*diffAction)

	for _, d := range diff {
		action := &diffAction{
			diff: d,
		}
		typ := diffActionType{rw: d.Object, delete: false}

		switch d.Type {
		case DiffTypeCreate, DiffTypeUpdate:
			action.wg.Add(len(d.Object.Dependencies))

		case DiffTypeDelete:
			action.wg.Add(len(d.Object.DependedBy))

			typ.delete = true

		case DiffTypeRecreate:
			preAction := &diffAction{
				diff: d,
			}
			preAction.wg.Add(len(d.Object.DependedBy))
			diffActionMap[diffActionType{rw: d.Object, delete: true}] = preAction

			action.wg.Add(len(d.Object.Dependencies))
		}

		diffActionMap[typ] = action
	}

	return diffActionMap
}

func handleDiffAction(ctx context.Context, d *Diff, del bool, callback func(*types.ApplyAction)) error {
	if callback == nil {
		callback = func(*types.ApplyAction) {}
	}

	switch d.Type {
	case DiffTypeCreate:
		callback(d.ToApplyAction(0, 1))

		err := d.Object.Resource.Create(ctx)
		if err != nil {
			return err
		}

		d.Object.Resource.SetNew(false)
		d.Object.MarkAllWantedAsCurrent()
		callback(d.ToApplyAction(1, 1))

	case DiffTypeUpdate:
		callback(d.ToApplyAction(0, 1))

		err := d.Object.Resource.Update(ctx)
		if err != nil {
			return err
		}

		d.Object.MarkAllWantedAsCurrent()
		callback(d.ToApplyAction(1, 1))

	case DiffTypeDelete:
		callback(d.ToApplyAction(0, 1))

		err := d.Object.Resource.Delete(ctx)
		if err != nil {
			return err
		}

		callback(d.ToApplyAction(1, 1))

	case DiffTypeRecreate:
		if del {
			callback(d.ToApplyAction(0, 2))

			err := d.Object.Resource.Delete(ctx)
			if err != nil {
				return err
			}

			callback(d.ToApplyAction(1, 2))
		} else {
			err := d.Object.Resource.Create(ctx)
			if err != nil {
				return err
			}

			d.Object.Resource.SetNew(false)
			d.Object.MarkAllWantedAsCurrent()
			callback(d.ToApplyAction(2, 2))
		}
	}

	return nil
}

func (r *Registry) Apply(ctx context.Context, diff []*Diff, callback func(*types.ApplyAction)) error {
	g, _ := errgroup.WithContext(ctx)
	pool, _ := errgroup.WithConcurrency(ctx, defaultConcurrency)

	diffActionMap := prepareDiffActionMap(diff)

	for t, action := range diffActionMap {
		t := t
		action := action

		g.Go(func() error {
			// Wait for all dependencies to finish.
			err := waitContext(ctx, &action.wg)
			if err != nil {
				return err
			}

			pool.Go(func() error {
				d := action.diff

				err = handleDiffAction(ctx, d, t.delete, callback)
				if err != nil {
					return err
				}

				if action.diff.Type == DiffTypeCreate || action.diff.Type == DiffTypeUpdate || (action.diff.Type == DiffTypeRecreate && !t.delete) {
					// Tell all objects that depepnd on me that I have been created.
					for _, dep := range action.diff.Object.DependedBy {
						a, ok := diffActionMap[diffActionType{rw: dep, delete: false}]
						if ok {
							a.wg.Done()
						}
					}
				} else {
					// Tell all my dependencies that I have been deleted.
					for _, dep := range action.diff.Object.Dependencies {
						a, ok := diffActionMap[diffActionType{rw: dep, delete: true}]
						if ok {
							a.wg.Done()
						}
					}
				}

				return nil
			})

			return nil
		})
	}

	err := g.Wait()
	if err != nil {
		return err
	}

	err = pool.Wait()

	return err
}
