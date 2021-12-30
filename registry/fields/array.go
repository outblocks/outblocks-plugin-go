package fields

import (
	"reflect"
)

type arrayField interface {
	SetCurrent([]interface{})
	LookupCurrent() ([]interface{}, bool)
	Current() []interface{}
}

type ArrayInputField interface {
	arrayField
	InputField

	LookupWanted() ([]interface{}, bool)
	Wanted() []interface{}
	SetWanted([]interface{})
	Any() []interface{}
}

type ArrayOutputField interface {
	arrayField
	OutputField

	Input() ArrayInputField
}

type ArrayBaseField struct {
	FieldBase
}

func Array(val []Field) ArrayInputField {
	return &ArrayBaseField{FieldBase: BasicValue(val, false)}
}

func ArrayUnset() ArrayInputField {
	return &ArrayBaseField{FieldBase: BasicValueUnset(false)}
}

func ArrayUnsetOutput() ArrayOutputField {
	return &ArrayBaseField{FieldBase: BasicValueUnset(true)}
}

func ArrayOutput(val []Field) ArrayOutputField {
	return &ArrayBaseField{FieldBase: BasicValue(val, true)}
}

func (f *ArrayBaseField) SetCurrent(i []interface{}) {
	f.setCurrent(interfaceArrayToFieldArray(i))
}

func (f *ArrayBaseField) LookupCurrent() (v []interface{}, ok bool) {
	if !f.currentDefined {
		return nil, f.currentDefined
	}

	return f.Serialize(f.currentVal).([]interface{}), true
}

func (f *ArrayBaseField) SetWanted(i []interface{}) {
	f.setWanted(interfaceArrayToFieldArray(i))
}

func (f *ArrayBaseField) LookupWanted() (v []interface{}, ok bool) {
	if !f.wantedDefined {
		return nil, false
	}

	return f.Serialize(f.wanted()).([]interface{}), true
}

func (f *ArrayBaseField) Wanted() []interface{} {
	v, _ := f.LookupWanted()
	return v
}

func (f *ArrayBaseField) Current() []interface{} {
	v, _ := f.LookupCurrent()
	return v
}

func (f *ArrayBaseField) Any() []interface{} {
	cur, ok := f.LookupCurrent()
	if ok {
		return cur
	}

	return f.Wanted()
}

func (f *ArrayBaseField) Serialize(i interface{}) interface{} {
	if i == nil {
		return make([]interface{}, 0)
	}

	c := i.([]Field)
	m := make([]interface{}, len(c))

	for i, v := range c {
		if v == nil {
			m[i] = v
		}

		val, ok := v.LookupCurrentRaw()
		if !ok {
			if ifield, ok := v.(InputField); ok {
				val, ok = ifield.LookupWantedRaw()
				if ok {
					m[i] = v.Serialize(val)
					continue
				}
			}
		}

		m[i] = v.Serialize(val)
	}

	return m
}

func (f *ArrayBaseField) FieldDependencies() []interface{} {
	if f.wanted() == nil {
		return nil
	}

	var deps []interface{}

	for _, v := range f.wanted().([]Field) {
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

func (f *ArrayBaseField) IsChanged() bool {
	if f.currentVal == nil || f.wanted() == nil || f.invalidated {
		return f.FieldBase.IsChanged()
	}

	cur := f.Current()
	wanted := f.Wanted()

	return !reflect.DeepEqual(cur, wanted)
}

func (f *ArrayBaseField) Input() ArrayInputField {
	return f
}

func (f *ArrayBaseField) EmptyValue() interface{} {
	var ret []interface{}
	return ret
}

func interfaceArrayToFieldArray(in []interface{}) []Field {
	o := make([]Field, len(in))

	for i, v := range in {
		o[i] = mapFieldFromInterface(v)
	}

	return o
}
