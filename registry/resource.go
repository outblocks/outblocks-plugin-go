package registry

import (
	"context"
	"encoding/json"

	"github.com/outblocks/outblocks-plugin-go/registry/fields"
)

type Resource interface {
	GetName() string
	SetNew(bool)
	IsNew() bool
	Read(ctx context.Context, meta interface{}) error
	Create(ctx context.Context, meta interface{}) error
	Update(ctx context.Context, meta interface{}) error
	Delete(ctx context.Context, meta interface{}) error
}

type ResourceTypeVerbose interface {
	GetType() string
}

type ResourceBase struct {
	new bool
}

func (b *ResourceBase) IsNew() bool {
	return b.new
}

func (b *ResourceBase) SetNew(v bool) {
	b.new = v
}

type ResourceID struct {
	ID        string `json:"id"`
	Namespace string `json:"namespace"`
	Type      string `json:"type"`
}

type ResourceSerialized struct {
	ResourceID
	Properties   map[string]interface{} `json:"properties,omitempty"`
	DependedBy   []ResourceID           `json:"depended_by,omitempty"`
	Dependencies []ResourceID           `json:"dependencies,omitempty"`
}

type ResourceWrapper struct {
	ResourceID

	Fields       map[string]*FieldInfo         `json:"-"`
	DependedBy   map[*ResourceWrapper]struct{} `json:"-"`
	Dependencies map[*ResourceWrapper]struct{} `json:"-"`
	Resource     Resource                      `json:"-"`
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
		if !d.Resource.IsNew() {
			dependedBy = append(dependedBy, d.ResourceID)
		}
	}

	for d := range w.Dependencies {
		if !d.Resource.IsNew() {
			deps = append(deps, d.ResourceID)
		}
	}

	return json.Marshal(ResourceSerialized{
		ResourceID: ResourceID{
			ID:        w.ID,
			Namespace: w.Namespace,
			Type:      w.Type,
		},
		Properties:   props,
		DependedBy:   dependedBy,
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