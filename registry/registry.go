package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
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
	existing    []*ResourceWrapper
}

func NewRegistry() *Registry {
	return &Registry{
		fieldMap:    make(map[interface{}]*ResourceWrapper),
		resourceMap: make(map[ResourceID]*ResourceWrapper),
		types:       make(map[string]*ResourceTypeInfo),
	}
}

func mapFieldTypeInfo(fieldsMap map[string]*FieldTypeInfo, t reflect.Type, prefix string) {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	for i := 0; i < t.NumField(); i++ {
		ft := t.Field(i)

		if ft.Anonymous {
			continue
		}

		if ft.Type.Kind() == reflect.Struct {
			mapFieldTypeInfo(fieldsMap, ft.Type, ft.Name+".")

			continue
		}

		def, defSet := ft.Tag.Lookup("default")

		fieldsMap[prefix+ft.Name] = &FieldTypeInfo{
			ReflectType: ft,
			Properties:  parseFieldPropertiesTag(ft.Tag.Get("state")),
			Default:     def,
			DefaultSet:  defSet,
		}
	}
}

func (r *Registry) RegisterType(o Resource) error {
	t := reflect.TypeOf(o)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// Check if type wasn't registered.
	if _, ok := r.types[t.Name()]; ok {
		return nil
	}

	fieldsMap := make(map[string]*FieldTypeInfo)

	mapFieldTypeInfo(fieldsMap, t, "")

	r.types[t.Name()] = &ResourceTypeInfo{
		Type:   t,
		Fields: fieldsMap,
	}

	return nil
}

func mapFieldInfo(rti *ResourceTypeInfo, fieldsMap map[string]*FieldInfo, t reflect.Type, v reflect.Value, prefix string) {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	for i := 0; i < t.NumField(); i++ {
		ft := t.Field(i)
		fv := v.Field(i)

		if ft.Anonymous {
			continue
		}

		if ft.Type.Kind() == reflect.Struct {
			mapFieldInfo(rti, fieldsMap, ft.Type, fv, ft.Name+".")

			continue
		}

		fName := prefix + ft.Name
		fieldsMap[fName] = &FieldInfo{
			Type:  rti.Fields[fName],
			Value: fv,
		}
	}
}

func generateResourceFields(o Resource, rti *ResourceTypeInfo) map[string]*FieldInfo {
	v := reflect.ValueOf(o)
	t := rti.Type

	fieldsMap := make(map[string]*FieldInfo)

	mapFieldInfo(rti, fieldsMap, t, v, "")

	return fieldsMap
}

func (r *Registry) Register(o Resource, namespace, id string) error {
	t := reflect.TypeOf(o)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	id = fmt.Sprintf("%s:%s", t.Name(), id)
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

	o.SetState(ResourceStateNew)

	rw := &ResourceWrapper{
		ResourceID:   resourceID,
		Fields:       generateResourceFields(o, tinfo),
		Resource:     o,
		DependedBy:   make(map[*ResourceWrapper]struct{}),
		Dependencies: make(map[*ResourceWrapper]struct{}),
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
			rw.Dependencies[erw] = struct{}{}
			erw.DependedBy[rw] = struct{}{}

			f.Value.Set(reflect.ValueOf(fields.MakeProxyField(f.Value.Interface())))

			continue
		}

		if fh, ok := f.Value.Interface().(fields.FieldHolder); ok {
			for _, dep := range fh.FieldDependencies() {
				if erw, ok := r.fieldMap[dep]; ok {
					rw.Dependencies[erw] = struct{}{}
					erw.DependedBy[rw] = struct{}{}
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

	var existing []*ResourceSerialized

	err := json.Unmarshal(state, &existing)
	if err != nil {
		return err
	}

	existingMap := make(map[ResourceID]*ResourceSerialized)
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

		res.Resource.SetState(ResourceStateExisting)
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
			ResourceID:   v.ResourceID,
			Fields:       generateResourceFields(res, rti),
			Resource:     obj.Interface().(Resource),
			DependedBy:   make(map[*ResourceWrapper]struct{}),
			Dependencies: make(map[*ResourceWrapper]struct{}),
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

			rw.DependedBy[dep] = struct{}{}
		}

		for _, d := range v.Dependencies {
			dep, ok := resourceMap[d]
			if !ok {
				return fmt.Errorf("dependency missing: %s", d)
			}

			rw.Dependencies[dep] = struct{}{}
		}
	}

	r.existing = missing

	return nil
}

func (r *Registry) Read(ctx context.Context, meta interface{}) error {
	pool, ctx := errgroup.WithConcurrency(ctx, defaultConcurrency)
	g, _ := errgroup.WithContext(ctx)

	resMap := make(map[*ResourceWrapper]*sync.WaitGroup, len(r.resources))

	for _, res := range r.resources {
		var s sync.WaitGroup

		s.Add(len(res.Dependencies))

		resMap[res] = &s
	}

	for res, wg := range resMap {
		res := res
		wg := wg

		g.Go(func() error {
			// Wait for all dependencies to finish.
			err := waitContext(ctx, wg)
			if err != nil {
				return err
			}

			pool.Go(func() error {
				err := res.Resource.Read(ctx, meta)
				if err != nil {
					return err
				}

				for dep := range res.DependedBy {
					resMap[dep].Done()
				}

				return nil
			})

			return nil
		})
	}

	err := g.Wait()
	if err != nil && err != context.Canceled {
		return err
	}

	return pool.Wait()
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
	arrayInputType   = reflect.TypeOf((*fields.ArrayInputField)(nil)).Elem()
	arrayOutputType  = reflect.TypeOf((*fields.ArrayOutputField)(nil)).Elem()
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

			// Array.
		case arrayInputType:
			val = fields.ArrayUnset()
		case arrayOutputType:
			val = fields.ArrayUnsetOutput()

		default:
			return fmt.Errorf("unknown field type %s", f.Type.ReflectType.Type)
		}

		f.Value.Set(reflect.ValueOf(val))
	}

	return nil
}

func (r *Registry) Diff(ctx context.Context, destroy bool) ([]*Diff, error) {
	diffMap := make(map[*ResourceWrapper]*Diff)

	// Add all missing resources as deletions.
	for _, o := range r.existing {
		deleteObjectTree(o, diffMap)
	}

	// Process other ops.
	for _, o := range r.resources {
		if destroy {
			deleteObjectTree(o, diffMap)
			continue
		}

		existing := diffMap[o]
		if existing != nil && existing.Type == DiffTypeRecreate {
			continue
		}

		d := calculateDiff(o, false)
		if d != nil {
			diffMap[o] = d

			if d.Type == DiffTypeRecreate {
				recreateObjectTree(d.Object, diffMap)
			}
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
	for _, res := range append(r.resources, r.existing...) {
		if !res.Resource.IsExisting() {
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

func calculateDependencyWaitGroup(diffActionMap map[diffActionType]*diffAction, deps map[*ResourceWrapper]struct{}, del bool) int {
	add := 0

	for d := range deps {
		if _, ok := diffActionMap[diffActionType{rw: d, delete: del}]; ok {
			add++
		}
	}

	return add
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
		case DiffTypeDelete:
			typ.delete = true

		case DiffTypeRecreate:
			// Add one more action before create.
			diffActionMap[diffActionType{rw: d.Object, delete: true}] = &diffAction{
				diff: d,
			}

			action.wg.Add(1)
		}

		diffActionMap[typ] = action
	}

	for t, a := range diffActionMap {
		if t.delete {
			a.wg.Add(calculateDependencyWaitGroup(diffActionMap, a.diff.Object.DependedBy, t.delete))
		} else {
			a.wg.Add(calculateDependencyWaitGroup(diffActionMap, a.diff.Object.Dependencies, t.delete))
		}
	}

	return diffActionMap
}

func handleDiffAction(ctx context.Context, meta interface{}, d *Diff, del bool, callback func(*types.ApplyAction)) error {
	if callback == nil {
		callback = func(*types.ApplyAction) {}
	}

	switch d.Type {
	case DiffTypeCreate:
		callback(d.ToApplyAction(0, 1))

		err := d.Object.Resource.Create(ctx, meta)
		if err != nil {
			return err
		}

		d.Object.Resource.SetState(ResourceStateExisting)
		d.Object.MarkAllWantedAsCurrent()
		callback(d.ToApplyAction(1, 1))

	case DiffTypeUpdate:
		callback(d.ToApplyAction(0, 1))

		err := d.Object.Resource.Update(ctx, meta)
		if err != nil {
			return err
		}

		d.Object.MarkAllWantedAsCurrent()
		callback(d.ToApplyAction(1, 1))

	case DiffTypeDelete:
		callback(d.ToApplyAction(0, 1))

		err := d.Object.Resource.Delete(ctx, meta)
		if err != nil {
			return err
		}

		d.Object.Resource.SetState(ResourceStateDeleted)
		callback(d.ToApplyAction(1, 1))

	case DiffTypeRecreate:
		if del {
			callback(d.ToApplyAction(0, 2))

			err := d.Object.Resource.Delete(ctx, meta)
			if err != nil {
				return err
			}

			d.Object.Resource.SetState(ResourceStateDeleted)
			callback(d.ToApplyAction(1, 2))
		} else {
			err := d.Object.Resource.Create(ctx, meta)
			if err != nil {
				return err
			}

			d.Object.Resource.SetState(ResourceStateExisting)
			d.Object.MarkAllWantedAsCurrent()
			callback(d.ToApplyAction(2, 2))
		}
	}

	return nil
}

func (r *Registry) Apply(ctx context.Context, meta interface{}, diff []*Diff, callback func(*types.ApplyAction)) error {
	pool, ctx := errgroup.WithConcurrency(ctx, defaultConcurrency)
	g, _ := errgroup.WithContext(ctx)

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

				err = handleDiffAction(ctx, meta, d, t.delete, callback)
				if err != nil {
					return fmt.Errorf("applying changes to %s '%s' error: %w", d.ObjectType(), d.Object.Resource.GetName(), err)
				}

				if action.diff.Type == DiffTypeCreate || action.diff.Type == DiffTypeUpdate || (action.diff.Type == DiffTypeRecreate && !t.delete) {
					// Tell all objects that depend on me that I have been created.
					for dep := range action.diff.Object.DependedBy {
						a, ok := diffActionMap[diffActionType{rw: dep, delete: false}]
						if ok {
							a.wg.Done()
						}
					}

					return nil
				}

				// When dealing with delete, tell all my dependencies that I have been deleted so they can be deleted safely as well.
				for dep := range action.diff.Object.Dependencies {
					a, ok := diffActionMap[diffActionType{rw: dep, delete: true}]
					if ok {
						a.wg.Done()
					}
				}

				if action.diff.Type == DiffTypeRecreate {
					// Tell myself that I have been deleted as well if we are dealing with recreate.
					a, ok := diffActionMap[diffActionType{rw: action.diff.Object, delete: false}]
					if ok {
						a.wg.Done()
					}
				}

				return nil
			})

			return nil
		})
	}

	err := g.Wait()
	if err != nil && err != context.Canceled {
		return err
	}

	return pool.Wait()
}
