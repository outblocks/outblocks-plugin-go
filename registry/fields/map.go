package fields

import (
	"fmt"
	"reflect"
)

type mapBaseField interface {
	SetCurrent(map[string]interface{})
	LookupCurrent() (map[string]interface{}, bool)
	Current() map[string]interface{}
}

type MapInputField interface {
	mapBaseField
	InputField

	LookupWanted() (map[string]interface{}, bool)
	Wanted() map[string]interface{}
	Any() map[string]interface{}
}

type MapOutputField interface {
	mapBaseField
	OutputField

	SetCurrent(map[string]interface{})
	LookupCurrent() (map[string]interface{}, bool)
	Current() map[string]interface{}
}

type MapField struct {
	FieldBase
}

func Map(val map[string]Field) MapInputField {
	return &MapField{FieldBase: BasicValue(val, false)}
}

func MapUnset() MapInputField {
	return &MapField{FieldBase: BasicValueUnset(false)}
}

func MapUnsetOutput() MapOutputField {
	return &MapField{FieldBase: BasicValueUnset(true)}
}

func MapOutput(val map[string]Field) MapOutputField {
	return &MapField{FieldBase: BasicValue(val, true)}
}

func (f *MapField) SetCurrent(i map[string]interface{}) {
	f.setCurrent(interfaceMapToFieldMap(i))
}

func (f *MapField) LookupCurrent() (v map[string]interface{}, ok bool) {
	if !f.currentDefined {
		return nil, f.currentDefined
	}

	return f.Serialize(f.current).(map[string]interface{}), true
}

func (f *MapField) LookupWanted() (v map[string]interface{}, ok bool) {
	if !f.wantedDefined {
		return nil, false
	}

	return f.Serialize(f.wanted).(map[string]interface{}), true
}

func (f *MapField) Wanted() map[string]interface{} {
	v, _ := f.LookupWanted()
	return v
}

func (f *MapField) Current() map[string]interface{} {
	v, _ := f.LookupCurrent()
	return v
}

func (f *MapField) Any() map[string]interface{} {
	cur, ok := f.LookupCurrent()
	if ok {
		return cur
	}

	return f.Wanted()
}

func (f *MapField) Serialize(i interface{}) interface{} {
	m := make(map[string]interface{})

	if i == nil {
		return m
	}

	for k, v := range i.(map[string]Field) {
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

func (f *MapField) FieldDependencies() []interface{} {
	if f.wanted == nil {
		return nil
	}

	var deps []interface{}

	for _, v := range f.wanted.(map[string]Field) {
		if v == nil {
			continue
		}

		if fh, ok := v.(FieldHolder); ok {
			deps = append(deps, fh.FieldDependencies()...)

			continue
		}

		deps = append(deps, v)
	}

	return deps
}

func (f *MapField) IsChanged() bool {
	if f.current == nil || f.wanted == nil {
		return f.FieldBase.IsChanged()
	}

	cur := f.Current()
	wanted := f.Wanted()

	return !reflect.DeepEqual(cur, wanted)
}

func mapFieldFromInterface(i interface{}) Field {
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

	case map[string]interface{}:
		o := MapUnset()
		o.SetCurrent(v)

		return o
	}

	panic(fmt.Sprintf("unmappable field: %+v", i))
}

func interfaceMapToFieldMap(i map[string]interface{}) map[string]Field {
	o := make(map[string]Field)

	for k, v := range i {
		o[k] = mapFieldFromInterface(v)
	}

	return o
}
