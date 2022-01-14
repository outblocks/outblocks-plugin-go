package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/outblocks/outblocks-plugin-go/registry/fields"
	"github.com/outblocks/outblocks-plugin-go/util"
)

type ResourceState int

const (
	ResourceStateNew ResourceState = iota + 1
	ResourceStateExisting
	ResourceStateDeleted
)

var mutexKV = util.NewMutexKV()

type Resource interface {
	setDiff(*Diff)
	setRegistered(bool)
	setWrapper(*ResourceWrapper)

	GetName() string
	Diff() *Diff
	Wrapper() *ResourceWrapper
	SetState(ResourceState)
	State() ResourceState

	IsRegistered() bool

	IsNew() bool
	MarkAsNew()
	IsExisting() bool
	MarkAsExisting()
	IsDeleted() bool
	MarkAsDeleted()

	SkipState() bool
}

type ResourceReader interface {
	Read(ctx context.Context, meta interface{}) error
}

type ResourceReference interface {
	ReferenceID() string
}

type ResourceIniter interface {
	Init(ctx context.Context, meta interface{}, opts *Options) error
}

type ResourceProcessor interface {
	Process(ctx context.Context, meta interface{}) error
}

type ResourceCUD interface {
	Create(ctx context.Context, meta interface{}) error
	Update(ctx context.Context, meta interface{}) error
	Delete(ctx context.Context, meta interface{}) error
}

type ResourceDiffCalculator interface {
	CalculateDiff() DiffType
}

type ResourceBeforeDiffHook interface {
	BeforeDiff()
}

type ResourceTypeVerbose interface {
	GetType() string
}

type ResourceCriticalChecker interface {
	IsCritical(t DiffType, fieldList []string) bool
}

type ResourceBase struct {
	state      ResourceState
	wrapper    *ResourceWrapper
	diff       *Diff
	registered bool
}

func (b *ResourceBase) setDiff(v *Diff) { // nolint:unused
	b.diff = v
}

func (b *ResourceBase) Diff() *Diff {
	return b.diff
}

func (b *ResourceBase) setWrapper(v *ResourceWrapper) { // nolint:unused
	b.wrapper = v
}

func (b *ResourceBase) Wrapper() *ResourceWrapper {
	return b.wrapper
}

func (b *ResourceBase) SetState(v ResourceState) {
	b.state = v
}

func (b *ResourceBase) State() ResourceState {
	return b.state
}

func (b *ResourceBase) IsRegistered() bool {
	return b.registered
}

func (b *ResourceBase) setRegistered(v bool) { // nolint:unused
	b.registered = v
}

func (b *ResourceBase) IsNew() bool {
	return b.State() == ResourceStateNew
}

func (b *ResourceBase) IsExisting() bool {
	return b.State() == ResourceStateExisting
}

func (b *ResourceBase) IsDeleted() bool {
	return b.State() == ResourceStateDeleted
}

func (b *ResourceBase) MarkAsNew() {
	b.SetState(ResourceStateNew)
}

func (b *ResourceBase) MarkAsExisting() {
	b.SetState(ResourceStateExisting)
}

func (b *ResourceBase) MarkAsDeleted() {
	b.SetState(ResourceStateDeleted)
}

func (b *ResourceBase) SkipState() bool {
	return false
}

func (b *ResourceBase) Lock(k string) {
	mutexKV.Lock(k)
}

func (b *ResourceBase) Unlock(k string) {
	mutexKV.Unlock(k)
}

type ResourceID struct {
	ID        string `json:"id"`
	Namespace string `json:"namespace"`
	Type      string `json:"type"`
	Source    string `json:"source"`
}

func (rid *ResourceID) Less(rid2 *ResourceID) bool {
	if rid.Source != rid2.Source {
		return rid.Source < rid2.Source
	}

	if rid.Namespace != rid2.Namespace {
		return rid.Namespace < rid2.Namespace
	}

	if rid.ID != rid2.ID {
		return rid.ID < rid2.ID
	}

	return rid.Type < rid2.Type
}

type ResourceSerialized struct {
	ResourceID
	IsNew        bool                   `json:"is_new,omitempty"`
	ReferenceID  string                 `json:"ref_id,omitempty"`
	Properties   map[string]interface{} `json:"properties,omitempty"`
	Dependencies []ResourceID           `json:"dependencies,omitempty"`
	DependedBy   []ResourceID           `json:"depended_by,omitempty"`
}

type ResourceWrapper struct {
	ResourceID

	Fields       map[string]*FieldInfo
	DependedBy   map[*ResourceWrapper]struct{}
	Dependencies map[*ResourceWrapper]struct{}
	Resource     Resource
	IsSkipped    bool
}

func (w *ResourceWrapper) String() string {
	return fmt.Sprintf("ResourceWrapper<ID=%s,Type=%s,Ns=%s>", w.ID, w.Type, w.Namespace)
}

func (w *ResourceWrapper) SetFieldValues(props map[string]interface{}) error {
	for k, v := range props {
		f, ok := w.Fields[k]
		if !ok || v == nil {
			continue
		}

		err := fields.SetFieldValue(f.Value.Interface(), v)
		if err != nil {
			return err
		}
	}

	return nil
}

func (w *ResourceWrapper) FieldList() []string {
	var f []string

	for k := range w.Fields {
		f = append(f, k)
	}

	return f
}

func (w *ResourceWrapper) MarshalJSON() ([]byte, error) {
	props := make(map[string]interface{})

	for k, v := range w.Fields {
		if v.Type.Properties.Ignored {
			continue
		}

		f := v.Value.Interface().(fields.Field)

		val, ok := f.LookupCurrentRaw()
		if ok {
			props[k] = f.Serialize(val)
		}
	}

	var (
		dependedBy []ResourceID
		deps       []ResourceID
	)

	for d := range w.DependedBy {
		if d.Resource.State() == ResourceStateExisting {
			dependedBy = append(dependedBy, d.ResourceID)
		}
	}

	sort.Slice(dependedBy, func(i, j int) bool {
		return dependedBy[i].Less(&dependedBy[j])
	})

	for d := range w.Dependencies {
		if d.Resource.State() == ResourceStateExisting {
			deps = append(deps, d.ResourceID)
		}
	}

	sort.Slice(deps, func(i, j int) bool {
		return deps[i].Less(&deps[j])
	})

	var refID string

	if rr, ok := w.Resource.(ResourceReference); ok {
		refID = rr.ReferenceID()
	}

	return json.Marshal(ResourceSerialized{
		IsNew:        w.Resource.IsNew(),
		ResourceID:   w.ResourceID,
		ReferenceID:  refID,
		Properties:   props,
		Dependencies: deps,
		DependedBy:   dependedBy,
	})
}

func (w *ResourceWrapper) MarkAllWantedAsCurrent() {
	for _, f := range w.Fields {
		switch i := f.Value.Interface().(type) {
		case fields.StringInputField:
			w, ok := i.LookupWanted()
			if ok {
				i.SetCurrent(w)
			}
		case fields.BoolInputField:
			w, ok := i.LookupWanted()
			if ok {
				i.SetCurrent(w)
			}
		case fields.IntInputField:
			w, ok := i.LookupWanted()
			if ok {
				i.SetCurrent(w)
			}
		case fields.MapInputField:
			w, ok := i.LookupWanted()
			if ok {
				i.SetCurrent(w)
			}
		case fields.ArrayInputField:
			w, ok := i.LookupWanted()
			if ok {
				i.SetCurrent(w)
			}
		}
	}
}

func (w *ResourceWrapper) UnsetAllCurrent() {
	for _, f := range w.Fields {
		if f.Type.Properties.Ignored || f.Value.Interface() == nil {
			continue
		}

		f.Value.Interface().(fields.Field).UnsetCurrent()
	}
}
