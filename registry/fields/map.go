package fields

import (
	"fmt"
	"reflect"
)

type mapField interface {
	SetCurrent(map[string]interface{})
	LookupCurrent() (map[string]interface{}, bool)
	Current() map[string]interface{}
}

type MapInputField interface {
	mapField
	InputField

	LookupWanted() (map[string]interface{}, bool)
	Wanted() map[string]interface{}
	SetWanted(map[string]interface{})
	Any() map[string]interface{}
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
	return &MapBaseField{FieldBase: BasicValueLazy(func() interface{} { return f() })}
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

func (f *MapBaseField) SetCurrent(i map[string]interface{}) {
	f.setCurrent(interfaceMapToFieldMap(i))
}

func (f *MapBaseField) LookupCurrent() (v map[string]interface{}, ok bool) {
	if !f.currentDefined {
		return nil, f.currentDefined
	}

	return f.Serialize(f.currentVal).(map[string]interface{}), true
}

func (f *MapBaseField) SetWanted(i map[string]interface{}) {
	f.setWanted(interfaceMapToFieldMap(i))
}

func (f *MapBaseField) LookupWanted() (v map[string]interface{}, ok bool) {
	if !f.wantedDefined {
		return nil, false
	}

	return f.Serialize(f.wanted()).(map[string]interface{}), true
}

func (f *MapBaseField) Wanted() map[string]interface{} {
	v, _ := f.LookupWanted()
	return v
}

func (f *MapBaseField) Current() map[string]interface{} {
	v, _ := f.LookupCurrent()
	return v
}

func (f *MapBaseField) Any() map[string]interface{} {
	cur, ok := f.LookupCurrent()
	if ok {
		return cur
	}

	return f.Wanted()
}

func (f *MapBaseField) Serialize(i interface{}) interface{} {
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

func (f *MapBaseField) FieldDependencies() []interface{} {
	if f.wanted() == nil {
		return nil
	}

	var deps []interface{}

	for _, v := range f.wanted().(map[string]Field) {
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

func (f *MapBaseField) EmptyValue() interface{} {
	var ret map[string]interface{}
	return ret
}
