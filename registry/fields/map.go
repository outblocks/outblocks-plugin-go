package fields

import (
	"reflect"
)

type MapInputField interface {
	InputField

	SetCurrent(map[string]interface{})
	LookupCurrent() (map[string]interface{}, bool)
	LookupWanted() (map[string]interface{}, bool)
	GetCurrent() map[string]interface{}
	GetWanted() map[string]interface{}
	GetAny() map[string]interface{}
}

type MapOutputField interface {
	OutputField

	SetCurrent(map[string]interface{})
	LookupCurrent() (map[string]interface{}, bool)
	GetCurrent() map[string]interface{}
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

	return fieldMapToInterfaceMap(f.current.(map[string]Field)), true
}

func (f *MapField) LookupWanted() (v map[string]interface{}, ok bool) {
	if !f.wantedDefined {
		return nil, false
	}

	return fieldMapToInterfaceMap(f.wanted.(map[string]Field)), true
}

func (f *MapField) GetWanted() map[string]interface{} {
	v, _ := f.LookupWanted()
	return v
}

func (f *MapField) GetCurrent() map[string]interface{} {
	v, _ := f.LookupCurrent()
	return v
}

func (f *MapField) GetAny() map[string]interface{} {
	any, defined := f.lookupAny()
	if !defined {
		return nil
	}

	return fieldMapToInterfaceMap(any.(map[string]Field))
}

func (f *MapField) SerializeValue() interface{} {
	if !f.currentDefined {
		return nil
	}

	m := make(map[string]interface{})

	for k, v := range f.current.(map[string]Field) {
		m[k] = v.SerializeValue()
	}

	return m
}

func (f *MapField) FieldDependencies() []interface{} {
	if f.current == nil {
		return nil
	}

	var deps []interface{}

	for _, v := range f.current.(map[string]Field) {
		_, ok := v.(InputField)
		if ok {
			deps = append(deps, v)
		}
	}

	return deps
}

func (f *MapField) IsChanged() bool {
	if f.current == nil || f.wanted == nil {
		return f.FieldBase.IsChanged()
	}

	cur := f.current.(map[string]Field)
	wanted := f.wanted.(map[string]Field)

	return !reflect.DeepEqual(fieldMapToInterfaceMap(cur), fieldMapToInterfaceMap(wanted))
}

func mapFieldFromInterface(i interface{}) Field {
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

		return o
	}

	panic("unmappable field")
}

func fieldMapToInterfaceMap(i map[string]Field) map[string]interface{} {
	o := make(map[string]interface{})

	for k, v := range i {
		o[k] = v.SerializeValue()
	}

	return o
}

func interfaceMapToFieldMap(i map[string]interface{}) map[string]Field {
	o := make(map[string]Field)

	for k, v := range i {
		o[k] = mapFieldFromInterface(v)
	}

	return o
}
