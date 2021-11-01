package registry

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/outblocks/outblocks-plugin-go/registry/fields"
	"github.com/outblocks/outblocks-plugin-go/util"
)

const (
	SourceApp        = "app"
	SourceDependency = "dependency"
	SourcePlugin     = "plugin"
)

type ResourceState int

const (
	ResourceStateNew ResourceState = iota + 1
	ResourceStateExisting
	ResourceStateDeleted
)

var mutexKV = util.NewMutexKV()

type Resource interface {
	GetName() string
	SetDiff(*Diff)
	Diff() *Diff
	SetState(ResourceState)
	State() ResourceState

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

type ResourceUnique interface {
	UniqueID() string
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

type ResourceBase struct {
	state ResourceState
	diff  *Diff
}

func (b *ResourceBase) SetDiff(v *Diff) {
	b.diff = v
}

func (b *ResourceBase) Diff() *Diff {
	return b.diff
}

func (b *ResourceBase) SetState(v ResourceState) {
	b.state = v
}

func (b *ResourceBase) State() ResourceState {
	return b.state
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

type ResourceSerialized struct {
	ResourceID
	Properties   map[string]interface{} `json:"properties,omitempty"`
	Dependencies []ResourceID           `json:"dependencies,omitempty"`
}

type ResourceWrapper struct {
	ResourceID

	Fields       map[string]*FieldInfo
	DependedBy   map[*ResourceWrapper]struct{}
	Dependencies map[*ResourceWrapper]struct{}
	Resource     Resource
	IsRegistered bool
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

	var deps []ResourceID

	for d := range w.Dependencies {
		if d.Resource.State() == ResourceStateExisting {
			deps = append(deps, d.ResourceID)
		}
	}

	return json.Marshal(ResourceSerialized{
		ResourceID:   w.ResourceID,
		Properties:   props,
		Dependencies: deps,
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
