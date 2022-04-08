package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"sync"

	apiv1 "github.com/outblocks/outblocks-plugin-go/gen/api/v1"
	"github.com/outblocks/outblocks-plugin-go/registry/fields"
	"github.com/outblocks/outblocks-plugin-go/util/errgroup"
)

const (
	defaultConcurrency = 5

	SourceApp        = "app"
	SourceDependency = "dependency"
	SourcePlugin     = "plugin"
)

type ResourceTypeInfo struct {
	Type   reflect.Type
	Fields map[string]*FieldTypeInfo
}

type Options struct {
	Read            bool
	Destroy         bool
	AllowDuplicates bool
}

type Registry struct {
	opts          *Options
	types         map[string]*ResourceTypeInfo
	fieldMap      map[interface{}]*ResourceWrapper
	skippedAppIDs map[string]bool

	resources       map[ResourceID]*ResourceWrapper
	loadedResources map[ResourceID]*ResourceSerialized
}

func NewRegistry(opts *Options) *Registry {
	if opts == nil {
		opts = &Options{}
	}

	return &Registry{
		opts:          opts,
		types:         make(map[string]*ResourceTypeInfo),
		fieldMap:      make(map[interface{}]*ResourceWrapper),
		skippedAppIDs: make(map[string]bool),
		resources:     make(map[ResourceID]*ResourceWrapper),
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

func (r *Registry) SkipAppResources(app *apiv1.App) {
	r.skippedAppIDs[app.Id] = true
}

func (r *Registry) RegisterAppResource(app *apiv1.App, id string, o Resource) (added bool, err error) {
	resID := r.createResourceID(SourceApp, app.Id, id, o)
	return r.register(resID, o)
}

func (r *Registry) RegisterDependencyResource(dep *apiv1.Dependency, id string, o Resource) (added bool, err error) {
	resID := r.createResourceID(SourceDependency, dep.Id, id, o)
	return r.register(resID, o)
}

func (r *Registry) RegisterPluginResource(scope, id string, o Resource) (added bool, err error) {
	resID := r.createResourceID(SourcePlugin, scope, id, o)
	return r.register(resID, o)
}

func (r *Registry) GetPluginResource(scope, id string, o Resource) (ok bool) {
	resID := r.createResourceID(SourcePlugin, scope, id, o)
	return r.get(resID, o)
}

func (r *Registry) GetAppResource(app *apiv1.App, id string, o Resource) (ok bool) {
	resID := r.createResourceID(SourceApp, app.Id, id, o)
	return r.get(resID, o)
}

func (r *Registry) GetDependencyResource(dep *apiv1.Dependency, id string, o Resource) (ok bool) {
	resID := r.createResourceID(SourceDependency, dep.Id, id, o)
	return r.get(resID, o)
}

func (r *Registry) createResourceID(source, namespace, id string, o Resource) ResourceID {
	t := reflect.TypeOf(o)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return ResourceID{
		ID:        id,
		Namespace: namespace,
		Type:      t.Name(),
		Source:    source,
	}
}

func (r *Registry) get(resourceID ResourceID, o Resource) (ok bool) {
	erw := r.resources[resourceID]
	if erw != nil {
		reflect.ValueOf(o).Elem().Set(reflect.ValueOf(erw.Resource).Elem())

		return true
	}

	return false
}

func (r *Registry) GetFieldDependencies(f fields.Field) []*ResourceWrapper {
	var ret []*ResourceWrapper

	if fh, ok := f.(fields.FieldDependencyHolder); ok {
		for _, dep := range fh.FieldDependencies() {
			rw := r.fieldMap[dep]
			if rw != nil {
				ret = append(ret, rw)
			}
		}
	}

	return ret
}

func (r *Registry) register(resourceID ResourceID, o Resource) (added bool, err error) {
	tinfo, ok := r.types[resourceID.Type]

	if !ok {
		err := r.RegisterType(o)
		if err != nil {
			return false, err
		}

		tinfo = r.types[resourceID.Type]
	}

	erw := r.resources[resourceID]
	if erw != nil && erw.Resource.IsRegistered() {
		if !r.opts.AllowDuplicates {
			return false, fmt.Errorf("resource already registered: %s", resourceID.ID)
		}

		reflect.ValueOf(o).Elem().Set(reflect.ValueOf(erw.Resource).Elem())

		return false, nil
	}

	o.SetState(ResourceStateNew)

	rw := &ResourceWrapper{
		ResourceID:   resourceID,
		Fields:       generateResourceFields(o, tinfo),
		Resource:     o,
		DependedBy:   make(map[*ResourceWrapper]struct{}),
		Dependencies: make(map[*ResourceWrapper]struct{}),
	}

	o.setRegistry(r)
	o.setWrapper(rw)
	o.setRegistered(true)

	err = setFieldDefaults(rw)
	if err != nil {
		return false, err
	}

	for _, f := range rw.Fields {
		if f.Type.Properties.Ignored {
			continue
		}

		if depRes, ok := r.fieldMap[f.Value.Interface()]; ok {
			rw.Dependencies[depRes] = struct{}{}
			depRes.DependedBy[rw] = struct{}{}

			f.Value.Set(reflect.ValueOf(fields.MakeProxyField(f.Value.Interface())))

			continue
		}

		if fh, ok := f.Value.Interface().(fields.FieldDependencyHolder); ok {
			for _, dep := range fh.FieldDependencies() {
				if depRes, ok := r.fieldMap[dep]; ok {
					rw.Dependencies[depRes] = struct{}{}
					depRes.DependedBy[rw] = struct{}{}
				}
			}
		}

		r.fieldMap[f.Value.Interface()] = rw
	}

	// Check if object has additional FieldDependencies.
	if fh, ok := o.(fields.FieldDependencyHolder); ok {
		for _, dep := range fh.FieldDependencies() {
			if depRes, ok := r.fieldMap[dep]; ok {
				rw.Dependencies[depRes] = struct{}{}
				depRes.DependedBy[rw] = struct{}{}
			}
		}
	}

	if erw != nil {
		err = rw.SetFieldValues(r.loadedResources[resourceID].Properties)
		if err != nil {
			return false, err
		}

		rw.Resource.SetState(erw.Resource.State())

		// Remove/fix existing dependencies<->depended by mapping.
		for k := range erw.Dependencies {
			delete(k.DependedBy, erw)

			rw.Dependencies[k] = struct{}{}
			k.DependedBy[rw] = struct{}{}
		}

		for k := range erw.DependedBy {
			delete(k.Dependencies, erw)

			rw.DependedBy[k] = struct{}{}
			k.Dependencies[rw] = struct{}{}
		}
	}

	r.resources[resourceID] = rw

	return true, nil
}

func (r *Registry) DeregisterAppResource(app *apiv1.App, id string, o Resource) error {
	resID := r.createResourceID(SourceApp, app.Id, id, o)
	return r.deregister(resID, o)
}

func (r *Registry) DeregisterDependencyResource(dep *apiv1.Dependency, id string, o Resource) error {
	resID := r.createResourceID(SourceDependency, dep.Id, id, o)
	return r.deregister(resID, o)
}

func (r *Registry) DeregisterPluginResource(scope, id string, o Resource) error {
	resID := r.createResourceID(SourcePlugin, scope, id, o)
	return r.deregister(resID, o)
}

func (r *Registry) deregister(resourceID ResourceID, o Resource) error {
	erw := r.resources[resourceID]
	if erw == nil || !erw.Resource.IsRegistered() {
		return fmt.Errorf("resource not registered: %s", resourceID.ID)
	}

	erw.Resource.setRegistered(false)

	if !erw.Resource.IsExisting() {
		delete(r.resources, resourceID)

		for k := range erw.Dependencies {
			delete(k.DependedBy, erw)
		}

		for k := range erw.DependedBy {
			delete(k.Dependencies, erw)
		}
	}

	reflect.ValueOf(o).Elem().Set(reflect.ValueOf(erw.Resource).Elem())

	return nil
}

func (r *Registry) Load(ctx context.Context, state []byte) error {
	if len(state) == 0 {
		return nil
	}

	var loaded []*ResourceSerialized

	err := json.Unmarshal(state, &loaded)
	if err != nil {
		return err
	}

	loadedMap := make(map[ResourceID]*ResourceSerialized)
	resourceMap := make(map[ResourceID]*ResourceWrapper)

	for _, e := range loaded {
		loadedMap[e.ResourceID] = e
	}

	for _, res := range r.resources {
		resourceMap[res.ResourceID] = res

		e, ok := loadedMap[res.ResourceID]
		if !ok {
			continue
		}

		delete(loadedMap, res.ResourceID)

		err := res.SetFieldValues(e.Properties)
		if err != nil {
			return err
		}

		if e.IsNew {
			res.Resource.SetState(ResourceStateNew)
		} else {
			res.Resource.SetState(ResourceStateExisting)
		}
	}

	// Process missing resources.
	missing := make(map[ResourceID]*ResourceWrapper)

	for _, v := range loadedMap {
		rti, ok := r.types[v.Type]
		if !ok {
			return fmt.Errorf("unknown resource type found: %s", v.Type)
		}

		obj := reflect.New(rti.Type)
		res := obj.Interface().(Resource)

		rw := &ResourceWrapper{
			ResourceID:   v.ResourceID,
			Fields:       generateResourceFields(res, rti),
			Resource:     res,
			DependedBy:   make(map[*ResourceWrapper]struct{}),
			Dependencies: make(map[*ResourceWrapper]struct{}),
		}

		res.setRegistry(r)
		res.setWrapper(rw)
		res.setRegistered(false)

		err := setFieldDefaults(rw)
		if err != nil {
			return err
		}

		err = rw.SetFieldValues(v.Properties)
		if err != nil {
			return err
		}

		resourceMap[v.ResourceID] = rw

		if v.IsNew {
			rw.Resource.SetState(ResourceStateNew)
		} else {
			rw.Resource.SetState(ResourceStateExisting)
		}

		missing[rw.ResourceID] = rw
	}

	// Fill dependencies now that resourceMap is filled.
	for _, v := range loadedMap {
		rw := resourceMap[v.ResourceID]

		for _, d := range v.DependedBy {
			dep, ok := resourceMap[d]
			if !ok {
				return fmt.Errorf("dependency missing: %s", d)
			}

			rw.DependedBy[dep] = struct{}{}
			dep.Dependencies[rw] = struct{}{}
		}

		for _, d := range v.Dependencies {
			dep, ok := resourceMap[d]
			if !ok {
				return fmt.Errorf("dependency missing: %s", d)
			}

			rw.Dependencies[dep] = struct{}{}
			dep.DependedBy[rw] = struct{}{}
		}
	}

	r.loadedResources = loadedMap

	for id, rw := range missing {
		r.resources[id] = rw
	}

	return nil
}

func (r *Registry) Process(ctx context.Context, meta interface{}) error {
	// Remove unregistered that are already defined with unique id.
	resourceUniqueIDMap := make(map[string]*ResourceWrapper)

	for _, rw := range r.resources {
		if rw.Resource.IsRegistered() {
			if rr, ok := rw.Resource.(ResourceReference); ok && rr.ReferenceID() != "" {
				resourceUniqueIDMap[rr.ReferenceID()] = rw
			}

			continue
		}
	}

	for id, rw := range r.resources {
		if rw.Resource.IsRegistered() {
			continue
		}

		// obsoleteID always == id when processing unregistered resource
		obsoleteID, err := mergeUniqueResource(rw, resourceUniqueIDMap)
		if err != nil {
			return err
		}

		if obsoleteID != nil {
			delete(r.resources, id)

			for k := range rw.Dependencies {
				delete(k.DependedBy, rw)
			}

			for k := range rw.DependedBy {
				delete(k.Dependencies, rw)
			}
		}
	}

	// Init where needed.
	err := r.init(ctx, meta)
	if err != nil {
		return err
	}

	// Read where needed.
	if r.opts.Read {
		err = r.read(ctx, meta)
		if err != nil {
			return err
		}
	}

	return nil
}

func mergeUniqueResource(res *ResourceWrapper, resourceUniqueIDMap map[string]*ResourceWrapper) (obsoleteID *ResourceID, err error) {
	// Merge potentially obsolete resources with newly registered.
	rr, ok := res.Resource.(ResourceReference)
	if !ok {
		return nil, nil
	}

	uniqID := rr.ReferenceID()
	if uniqID == "" {
		return nil, nil
	}

	existing, ok := resourceUniqueIDMap[uniqID]

	if ok {
		if existing.Resource.IsRegistered() == res.Resource.IsRegistered() {
			return nil, fmt.Errorf("multiple resources registered with same unique ID! one is: %s, another: %s",
				res.ResourceID, existing.ResourceID)
		}

		if res.Resource.IsRegistered() {
			resourceUniqueIDMap[uniqID] = res
			return &existing.ResourceID, nil
		}

		return &res.ResourceID, nil
	}

	resourceUniqueIDMap[uniqID] = res

	return nil, nil
}

func (r *Registry) init(ctx context.Context, meta interface{}) error {
	return r.processInOrder(ctx, defaultConcurrency, func(res *ResourceWrapper) error {
		if rr, ok := res.Resource.(ResourceIniter); ok {
			return rr.Init(ctx, meta, r.opts)
		}

		return nil
	})
}

func (r *Registry) read(ctx context.Context, meta interface{}) error {
	var (
		obsoleteIDs []*ResourceID
		mu          sync.Mutex
	)

	resourceUniqueIDMap := make(map[string]*ResourceWrapper)

	r.checkResources(r.resources)

	err := r.processInOrder(ctx, defaultConcurrency, func(res *ResourceWrapper) error {
		if res.IsSkipped {
			return nil
		}

		rr, ok := res.Resource.(ResourceReader)
		if !ok {
			return nil
		}

		// Merge potentially obsolete resources with newly registered.
		mu.Lock()
		defer mu.Unlock()

		// Skip reading objects that do not have unique id defined.
		if rr, ok := res.Resource.(ResourceReference); ok && rr.ReferenceID() == "" {
			return nil
		}

		obsoleteID, err := mergeUniqueResource(res, resourceUniqueIDMap)
		if err != nil {
			return err
		}
		if obsoleteID != nil {
			obsoleteIDs = append(obsoleteIDs, obsoleteID)
		}

		return rr.Read(ctx, meta)
	})
	if err != nil {
		return err
	}

	// Mark obsolete as deleted.
	for _, id := range obsoleteIDs {
		r.resources[*id].Resource.MarkAsDeleted()
	}

	return nil
}

func (r *Registry) processInOrder(ctx context.Context, concurrency int, f func(res *ResourceWrapper) error) error {
	var pool errgroup.Runner

	if concurrency > 0 {
		pool, ctx = errgroup.WithConcurrency(ctx, defaultConcurrency)
	} else {
		pool, ctx = errgroup.WithContext(ctx)
	}

	g, _ := errgroup.WithContext(ctx)

	resMap := make(map[*ResourceWrapper]*sync.WaitGroup, len(r.resources))

	for _, res := range r.resources {
		var wg sync.WaitGroup

		for dep := range res.Dependencies {
			if _, ok := r.resources[dep.ResourceID]; ok {
				wg.Add(1)
			}
		}

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
					if depWg, ok := resMap[dep]; ok {
						depWg.Done()
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
		val = fields.BoolUnsetOutput()

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

func unskipRecursiveDependencies(rw *ResourceWrapper) {
	if !rw.IsSkipped {
		return
	}

	rw.IsSkipped = false

	for dep := range rw.Dependencies {
		if rw.Resource.IsExisting() {
			continue
		}

		unskipRecursiveDependencies(dep)
	}
}

func (r *Registry) checkResources(resources map[ResourceID]*ResourceWrapper) {
	for _, rw := range resources {
		if rw.Source == SourceApp && r.skippedAppIDs[rw.Namespace] {
			rw.IsSkipped = true
		}
	}

	for _, rw := range resources {
		if !rw.IsSkipped {
			unskipRecursiveDependencies(rw)
		}
	}
}

func (r *Registry) Diff(ctx context.Context) ([]*Diff, error) {
	r.checkResources(r.resources)

	var mu sync.RWMutex

	diffMap := make(map[*ResourceWrapper]*Diff)

	// Process actual diff.
	err := r.processInOrder(ctx, -1, func(res *ResourceWrapper) error {
		if res.IsSkipped {
			return nil
		}

		// Add all missing resources as deletions.
		if !res.Resource.IsRegistered() {
			mu.Lock()
			deleteObjectTree(res, diffMap, true)
			mu.Unlock()

			return nil
		}

		if r.opts.Destroy {
			mu.Lock()
			deleteObjectTree(res, diffMap, false)
			mu.Unlock()

			return nil
		}

		mu.RLock()
		existing := diffMap[res]
		mu.RUnlock()

		if existing != nil && existing.Type == DiffTypeRecreate {
			return nil
		}

		d := r.calculateDiff(res)
		if d != nil {
			res.Resource.setDiff(d)

			mu.Lock()
			diffMap[res] = d

			if d.Type == DiffTypeRecreate {
				for dep := range res.DependedBy {
					dep.IsSkipped = false
				}
			}
			mu.Unlock()
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	var diff []*Diff

	for res, d := range diffMap {
		diff = append(diff, d)
		res.Resource.setDiff(d)
	}

	return diff, nil
}

func (r *Registry) Dump() ([]byte, error) {
	var resources []*ResourceWrapper

	for _, res := range r.resources {
		if res.Resource.IsDeleted() || res.Resource.SkipState() {
			continue
		}

		resources = append(resources, res)
	}

	sort.Slice(resources, func(i, j int) bool {
		return resources[i].Less(&resources[j].ResourceID)
	})

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

func handleDiffAction(ctx context.Context, meta interface{}, d *Diff, callback func(*apiv1.ApplyAction)) error {
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

			d.Object.UnsetAllCurrent()
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
	// During creation/update/process - wait for all dependencies to finish on create/update/process.
	if d.Type == DiffTypeCreate || d.Type == DiffTypeUpdate || d.Type == DiffTypeProcess || (d.Type == DiffTypeRecreate && step == 1) {
		for dep := range d.Object.Dependencies {
			resDiff := dep.Resource.Diff()

			// No need to wait for deletions, otherwise wait for dependency.
			if resDiff != nil && resDiff.Type != DiffTypeDelete {
				if err := resDiff.WaitContext(ctx, -1); err != nil {
					return err
				}
			}
		}

		return nil
	}

	// During deletion - wait for objects that depend on me to be deleted first.
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

func (r *Registry) Apply(ctx context.Context, meta interface{}, diff []*Diff, callback func(*apiv1.ApplyAction)) error {
	pool, ctx := errgroup.WithConcurrency(ctx, defaultConcurrency)

	for _, d := range diff {
		d.Object.Resource.setDiff(d)
	}

	// Add another cancel for errgroup context so that we can handle error from pool at all times.
	ctxErrgroup, cancelErrgroup := context.WithCancel(ctx)

	defer cancelErrgroup()

	g, _ := errgroup.WithContext(ctxErrgroup)

	var mu sync.Mutex

	cb := func(a *apiv1.ApplyAction) {
		if callback != nil {
			mu.Lock()
			callback(a)
			mu.Unlock()
		}
	}

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
					err := handleDiffAction(ctx, meta, d, cb)
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

	err = pool.Wait()

	// Cleanup resources.
	res := make(map[ResourceID]*ResourceWrapper, len(r.resources))

	for k, rw := range r.resources {
		if !rw.Resource.IsNew() {
			res[k] = rw
		}
	}

	r.resources = res

	return err
}

func (r *Registry) calculateDiff(rw *ResourceWrapper) *Diff {
	if rdc, ok := rw.Resource.(ResourceDiffCalculator); ok {
		typ := rdc.CalculateDiff()
		if typ != DiffTypeNone {
			return NewDiff(rw, typ, rw.FieldList())
		}

		return nil
	}

	if rbdh, ok := rw.Resource.(ResourceBeforeDiffHook); ok {
		rbdh.BeforeDiff()
	}

	if rw.Resource.IsNew() {
		typ := DiffTypeCreate

		return NewDiff(rw, typ, rw.FieldList())
	}

	forceNew := false

	var fieldsList []string

	for name, f := range rw.Fields {
		fieldChanged, fieldForceNew := r.calculateFieldDiff(f)
		if !fieldChanged {
			continue
		}

		fieldsList = append(fieldsList, name)

		if fieldForceNew {
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

	return NewDiff(rw, typ, fieldsList)
}

func (r *Registry) calculateFieldDiff(field *FieldInfo) (changed, forceNew bool) {
	if field.Type.Properties.Ignored {
		return false, false
	}

	v := field.Value.Interface().(fields.ValueTracker)

	// If field is a dep holder (depends on multiple fields) check each field if it is associated with a resource to be deleted/recreated.
	if fdh, ok := field.Value.Interface().(fields.FieldDependencyHolder); ok {
		for _, fd := range fdh.FieldDependencies() {
			if dep, ok := r.fieldMap[fd]; ok && dep.Resource.Diff() != nil {
				if dep.Resource.Diff().Type == DiffTypeRecreate || (dep.Resource.IsExisting() && field.Value.Interface().(fields.Field).IsOutput()) {
					return true, field.Type.Properties.ForceNew || field.Type.Properties.HardLink
				}
			}
		}
	}

	return v.IsChanged(), field.Type.Properties.ForceNew
}
