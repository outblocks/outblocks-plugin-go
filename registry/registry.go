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

type Options struct {
	Read bool
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

		if ft.Anonymous || !fv.CanSet() {
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

		if fh, ok := f.Value.Interface().(fields.FieldDependencyHolder); ok {
			for _, dep := range fh.FieldDependencies() {
				if erw, ok := r.fieldMap[dep]; ok {
					rw.Dependencies[erw] = struct{}{}
					erw.DependedBy[rw] = struct{}{}
				}
			}
		}

		r.fieldMap[f.Value.Interface()] = rw
	}

	// Check if object has additional FieldDependencies.
	if fh, ok := o.(fields.FieldDependencyHolder); ok {
		for _, dep := range fh.FieldDependencies() {
			if erw, ok := r.fieldMap[dep]; ok {
				rw.Dependencies[erw] = struct{}{}
				erw.DependedBy[rw] = struct{}{}
			}
		}
	}

	r.resources = append(r.resources, rw)

	return nil
}

func (r *Registry) Load(ctx context.Context, state []byte, meta interface{}, opts *Options) error {
	if opts == nil {
		opts = &Options{}
	}

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

		rw.Resource.SetState(ResourceStateExisting)

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

	// Init where needed.
	err = r.init(ctx, meta, opts)
	if err != nil {
		return err
	}

	// Read where needed.
	err = r.read(ctx, meta)
	if err != nil {
		return err
	}

	return nil
}

func (r *Registry) init(ctx context.Context, meta interface{}, opts *Options) error {
	return r.processInOrder(ctx, func(res *ResourceWrapper) error {
		if rr, ok := res.Resource.(ResourceIniter); ok {
			return rr.Init(ctx, meta, opts)
		}

		return nil
	})
}

func (r *Registry) read(ctx context.Context, meta interface{}) error {
	return r.processInOrder(ctx, func(res *ResourceWrapper) error {
		if rr, ok := res.Resource.(ResourceReader); ok {
			return rr.Read(ctx, meta)
		}

		return nil
	})
}

func (r *Registry) processInOrder(ctx context.Context, f func(res *ResourceWrapper) error) error {
	pool, ctx := errgroup.WithConcurrency(ctx, defaultConcurrency)
	g, _ := errgroup.WithContext(ctx)

	resMap := make(map[*ResourceWrapper]*sync.WaitGroup, len(r.resources))

	for _, res := range r.resources {
		var wg sync.WaitGroup

		wg.Add(len(res.Dependencies))

		resMap[res] = &wg
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
				err := f(res)
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

func mapFieldDefaultValue(typ *FieldTypeInfo) interface{} {
	defaultTag := typ.Default
	ok := typ.DefaultSet

	var val interface{}

	switch typ.ReflectType.Type {
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
	}

	return val
}

func setFieldDefaults(r *ResourceWrapper) error {
	for _, f := range r.Fields {
		if f.Type.Properties.Ignored || !f.Value.IsNil() || !f.Value.CanSet() {
			continue
		}

		if f.Type.Properties.Computed && !f.Value.IsNil() {
			return fmt.Errorf("%s.%s: computed field set to non-nil", r.Type, f.Type.ReflectType.Name)
		}

		val := mapFieldDefaultValue(f.Type)
		if val == nil {
			return fmt.Errorf("%s.%s: unknown field type %s", r.Type, f.Type.ReflectType.Name, f.Type.ReflectType.Type)
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

	for k, v := range diffMap {
		if v == nil {
			continue
		}

		k.Resource.SetDiff(v)

		diff = append(diff, v)
	}

	return diff, nil
}

func (r *Registry) Dump() ([]byte, error) {
	var resources []*ResourceWrapper
	for _, res := range append(r.resources, r.existing...) {
		if !res.Resource.IsExisting() || res.Resource.SkipState() {
			continue
		}

		resources = append(resources, res)
	}

	return json.Marshal(resources)
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

func handleDiffAction(ctx context.Context, meta interface{}, d *Diff, callback func(*types.ApplyAction)) error {
	if callback == nil {
		callback = func(*types.ApplyAction) {}
	}

	switch d.Type {
	case DiffTypeCreate:
		callback(d.ToApplyAction(0, 1))

		err := d.Object.Resource.(ResourceCUD).Create(ctx, meta)
		if err != nil {
			return err
		}

		d.Object.Resource.SetState(ResourceStateExisting)
		d.Object.MarkAllWantedAsCurrent()
		callback(d.ToApplyAction(1, 1))

	case DiffTypeUpdate:
		callback(d.ToApplyAction(0, 1))

		err := d.Object.Resource.(ResourceCUD).Update(ctx, meta)
		if err != nil {
			return err
		}

		d.Object.MarkAllWantedAsCurrent()
		callback(d.ToApplyAction(1, 1))

	case DiffTypeProcess:
		callback(d.ToApplyAction(0, 1))

		err := d.Object.Resource.(ResourceProcessor).Process(ctx, meta)
		if err != nil {
			return err
		}

		d.Object.MarkAllWantedAsCurrent()
		callback(d.ToApplyAction(1, 1))

	case DiffTypeDelete:
		callback(d.ToApplyAction(0, 1))

		err := d.Object.Resource.(ResourceCUD).Delete(ctx, meta)
		if err != nil {
			return err
		}

		d.Object.Resource.SetState(ResourceStateDeleted)
		callback(d.ToApplyAction(1, 1))

	case DiffTypeRecreate:
		if d.AppliedSteps() == 0 {
			callback(d.ToApplyAction(0, 2))

			err := d.Object.Resource.(ResourceCUD).Delete(ctx, meta)
			if err != nil {
				return err
			}

			d.Object.Resource.SetState(ResourceStateDeleted)
			callback(d.ToApplyAction(1, 2))
		} else {
			err := d.Object.Resource.(ResourceCUD).Create(ctx, meta)
			if err != nil {
				return err
			}

			d.Object.Resource.SetState(ResourceStateExisting)
			d.Object.MarkAllWantedAsCurrent()
			callback(d.ToApplyAction(2, 2))
		}

	case DiffTypeNone:
		panic("unexpected diff type")
	default:
		panic("unknown diff type")
	}

	return nil
}

func waitForDiffDeps(ctx context.Context, d *Diff, step int) error {
	// Wait for all dependencies to finish on create/update/process.
	if d.Type == DiffTypeCreate || d.Type == DiffTypeUpdate || d.Type == DiffTypeProcess || (d.Type == DiffTypeRecreate && step == 1) {
		for dep := range d.Object.Dependencies {
			resDiff := dep.Resource.Diff()

			if resDiff != nil {
				if err := resDiff.WaitContext(ctx, -1); err != nil {
					return err
				}
			}
		}

		return nil
	}

	// Wait for objects that depend on me to be deleted first.
	if d.Type == DiffTypeDelete || (d.Type == DiffTypeRecreate && step == 0) {
		for dep := range d.Object.DependedBy {
			resDiff := dep.Resource.Diff()

			// Only wait for deletions/first-phase recreations or (if it's a simple deletion) updates.
			if resDiff != nil && (resDiff.Type == DiffTypeDelete || resDiff.Type == DiffTypeRecreate || (d.Type == DiffTypeDelete && resDiff.Type == DiffTypeUpdate)) {
				if err := resDiff.WaitContext(ctx, 1); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (r *Registry) Apply(ctx context.Context, meta interface{}, diff []*Diff, callback func(*types.ApplyAction)) error {
	pool, ctx := errgroup.WithConcurrency(ctx, defaultConcurrency)

	// Add another cancel for errgroup context so that we can handle error from pool at all times.
	ctxErrgroup, cancelErrgroup := context.WithCancel(ctx)

	defer cancelErrgroup()

	g, _ := errgroup.WithContext(ctxErrgroup)

	for _, d := range diff {
		d := d

		for step := 0; step < d.RequiredSteps(); step++ {
			step := step

			g.Go(func() error {
				if step > 0 {
					if err := d.WaitContext(ctxErrgroup, step); err != nil {
						return err
					}
				}

				err := waitForDiffDeps(ctxErrgroup, d, step)
				if err != nil {
					return err
				}

				pool.Go(func() error {
					err := handleDiffAction(ctx, meta, d, callback)
					if err != nil {
						cancelErrgroup()

						return fmt.Errorf("applying changes to %s '%s' error: %w", d.ObjectType(), d.Object.Resource.GetName(), err)
					}

					d.MarkStepAsApplied()

					return nil
				})

				return nil
			})
		}
	}

	err := g.Wait()
	if err != nil && err != context.Canceled {
		return err
	}

	return pool.Wait()
}
