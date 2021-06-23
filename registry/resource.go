package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

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

type FieldTypeInfo struct {
	ReflectType reflect.StructField
	Properties  *FieldProperties
	Default     string
	DefaultSet  bool
}

type FieldInfo struct {
	Type  *FieldTypeInfo
	Value reflect.Value
}

type ResourceID struct {
	ID        string `json:"id"`
	Namespace string `json:"namespace"`
	Type      string `json:"type"`
}

type ResourceData struct {
	ResourceID
	Properties   map[string]interface{} `json:"properties"`
	DependedBy   []ResourceID           `json:"depended_by"`
	Dependencies []ResourceID           `json:"dependencies"`
}

type ResourceWrapper struct {
	ResourceID

	Fields       map[string]*FieldInfo `json:"-"`
	DependedBy   []*ResourceWrapper    `json:"-"`
	Dependencies []*ResourceWrapper    `json:"-"`
	Resource     Resource              `json:"-"`
}

func (w *ResourceWrapper) SetFieldValues(props map[string]interface{}) error {
	for k, v := range props {
		f, ok := w.Fields[k]
		if !ok || v == nil {
			continue
		}

		switch val := f.Value.Interface().(type) {
		case fields.StringInputField:
			val.SetCurrent(v.(string))
		case fields.BoolInputField:
			val.SetCurrent(v.(bool))
		case fields.IntInputField:
			switch i := v.(type) {
			case float64:
				val.SetCurrent(int(i))
			case int64:
				val.SetCurrent(int(i))
			case int:
				val.SetCurrent(i)
			default:
				return fmt.Errorf("unknown int field input found: %s", v)
			}
		case fields.MapInputField:
			val.SetCurrent(v.(map[string]interface{}))
		default:
			return fmt.Errorf("unknown field type found: %s", f.Value.Type())
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

		props[k] = v.Value.Interface().(fields.Field).SerializeValue()
	}

	var (
		dependedBy []ResourceID
		deps       []ResourceID
	)

	for _, d := range w.DependedBy {
		if !d.Resource.IsNew() {
			dependedBy = append(dependedBy, d.ResourceID)
		}
	}

	for _, d := range w.Dependencies {
		if !d.Resource.IsNew() {
			deps = append(deps, d.ResourceID)
		}
	}

	return json.Marshal(ResourceData{
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
		i, ok := f.Value.Interface().(fields.InputField)
		if !ok {
			continue
		}

		i.SetWantedAsCurrent()
	}
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
