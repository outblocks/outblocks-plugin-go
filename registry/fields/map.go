package fields

import (
	"fmt"
	"reflect"
)

type mapField interface {
	SetCurrent(map[string]any)
	LookupCurrent() (map[string]any, bool)
	Current() map[string]any
	WantedFieldMap() map[string]Field
	CurrentFieldMap() map[string]Field
}

type MapInputField interface {
	mapField
	InputField

	LookupWanted() (map[string]any, bool)
	Wanted() map[string]any
	SetWanted(map[string]any)
	Any() map[string]any
}

type MapOutputField interface {
	mapField
	OutputField

	Input() MapInputField
}

type MapBaseField struct {
	FieldBase
}

func Map(val map[string]Field) MapInputField {
	return &MapBaseField{FieldBase: BasicValue(val, false)}
}

func MapLazy(f func() map[string]Field) MapInputField {
	return &MapBaseField{FieldBase: BasicValueLazy(func() any { return f() })}
}

func MapUnset() MapInputField {
	return &MapBaseField{FieldBase: BasicValueUnset(false)}
}

func MapUnsetOutput() MapOutputField {
	return &MapBaseField{FieldBase: BasicValueUnset(true)}
}

func MapOutput(val map[string]Field) MapOutputField {
	return &MapBaseField{FieldBase: BasicValue(val, true)}
}

func (f *MapBaseField) WantedFieldMap() map[string]Field {
	if !f.wantedDefined {
		return nil
	}

	return f.wanted().(map[string]Field) //nolint:errcheck
}

func (f *MapBaseField) CurrentFieldMap() map[string]Field {
	if !f.currentDefined {
		return nil
	}

	return f.currentVal.(map[string]Field) //nolint:errcheck
}

func (f *MapBaseField) SetCurrent(i map[string]any) {
	f.setCurrent(interfaceMapToFieldMap(i))
}

func (f *MapBaseField) LookupCurrent() (v map[string]any, ok bool) {
	if !f.currentDefined {
		return nil, f.currentDefined
	}

	val, ok := f.Serialize(f.currentVal).(map[string]any)

	return val, ok
}

func (f *MapBaseField) SetWanted(i map[string]any) {
	f.setWanted(interfaceMapToFieldMap(i))
}

func (f *MapBaseField) LookupWanted() (v map[string]any, ok bool) {
	if !f.wantedDefined {
		return nil, false
	}

	val, ok := f.Serialize(f.wanted()).(map[string]any)

	return val, ok
}

func (f *MapBaseField) Wanted() map[string]any {
	v, _ := f.LookupWanted()
	return v
}

func (f *MapBaseField) Current() map[string]any {
	v, _ := f.LookupCurrent()
	return v
}

func (f *MapBaseField) Any() map[string]any {
	cur, ok := f.LookupCurrent()
	if ok {
		return cur
	}

	return f.Wanted()
}

func (f *MapBaseField) Serialize(i any) any {
	m := make(map[string]any)

	if i == nil {
		return m
	}

	for k, v := range i.(map[string]Field) { //nolint:errcheck
		if v == nil {
			m[k] = v

			continue
		}

		val, ok := v.LookupCurrentRaw()
		if !ok {
			if ifield, ok := v.(InputField); ok {
				val, ok = ifield.LookupWantedRaw()
				if ok {
					m[k] = v.Serialize(val)
					continue
				}
			}
		}

		m[k] = v.Serialize(val)
	}

	return m
}

func (f *MapBaseField) FieldDependencies() []any {
	if f.wanted() == nil {
		return nil
	}

	var deps []any

	for _, v := range f.wanted().(map[string]Field) { //nolint:errcheck
		if v == nil {
			continue
		}

		if fh, ok := v.(FieldDependencyHolder); ok {
			deps = append(deps, fh.FieldDependencies()...)

			continue
		}

		deps = append(deps, v)
	}

	return deps
}

func (f *MapBaseField) IsChanged() bool {
	if f.currentVal == nil || f.wanted() == nil || f.invalidated {
		return f.FieldBase.IsChanged()
	}

	cur := f.Current()
	wanted := f.Wanted()

	return !reflect.DeepEqual(cur, wanted)
}

func (f *MapBaseField) Input() MapInputField {
	return f
}

func mapFieldFromInterface(i any) Field {
	if i == nil {
		return nil
	}

	switch v := i.(type) {
	case int:
		o := IntUnset()
		o.SetCurrent(v)

		return o
	case string:
		o := StringUnset()
		o.SetCurrent(v)

		return o
	case bool:
		o := BoolUnset()
		o.SetCurrent(v)

	case map[string]any:
		o := MapUnset()
		o.SetCurrent(v)

		return o
	}

	panic(fmt.Sprintf("unmappable field: %+v", i))
}

func interfaceMapToFieldMap(i map[string]any) map[string]Field {
	o := make(map[string]Field)

	for k, v := range i {
		o[k] = mapFieldFromInterface(v)
	}

	return o
}

func (f *MapBaseField) EmptyValue() any {
	var ret map[string]any
	return ret
}
