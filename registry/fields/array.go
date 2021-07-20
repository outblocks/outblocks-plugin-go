package fields

import (
	"reflect"
)

type arrayBaseField interface {
	SetCurrent([]interface{})
	LookupCurrent() ([]interface{}, bool)
	Current() []interface{}
}

type ArrayInputField interface {
	arrayBaseField
	InputField

	LookupWanted() ([]interface{}, bool)
	Wanted() []interface{}
	Any() []interface{}
}

type ArrayOutputField interface {
	arrayBaseField
	OutputField
}

type ArrayField struct {
	FieldBase
}

func Array(val []Field) ArrayInputField {
	return &ArrayField{FieldBase: BasicValue(val, false)}
}

func ArrayUnset() ArrayInputField {
	return &ArrayField{FieldBase: BasicValueUnset(false)}
}

func ArrayUnsetOutput() ArrayOutputField {
	return &ArrayField{FieldBase: BasicValueUnset(true)}
}

func ArrayOutput(val []Field) ArrayOutputField {
	return &ArrayField{FieldBase: BasicValue(val, true)}
}

func (f *ArrayField) SetCurrent(i []interface{}) {
	f.setCurrent(interfaceArrayToFieldArray(i))
}

func (f *ArrayField) LookupCurrent() (v []interface{}, ok bool) {
	if !f.currentDefined {
		return nil, f.currentDefined
	}

	return f.Serialize(f.current).([]interface{}), true
}

func (f *ArrayField) LookupWanted() (v []interface{}, ok bool) {
	if !f.wantedDefined {
		return nil, false
	}

	return f.Serialize(f.wanted).([]interface{}), true
}

func (f *ArrayField) Wanted() []interface{} {
	v, _ := f.LookupWanted()
	return v
}

func (f *ArrayField) Current() []interface{} {
	v, _ := f.LookupCurrent()
	return v
}

func (f *ArrayField) Any() []interface{} {
	cur, ok := f.LookupCurrent()
	if ok {
		return cur
	}

	return f.Wanted()
}

func (f *ArrayField) Serialize(i interface{}) interface{} {
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

func (f *ArrayField) FieldDependencies() []interface{} {
	if f.wanted == nil {
		return nil
	}

	var deps []interface{}

	for _, v := range f.wanted.([]Field) {
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

func (f *ArrayField) IsChanged() bool {
	if f.current == nil || f.wanted == nil {
		return f.FieldBase.IsChanged()
	}

	cur := f.Current()
	wanted := f.Wanted()

	return !reflect.DeepEqual(cur, wanted)
}

func interfaceArrayToFieldArray(in []interface{}) []Field {
	o := make([]Field, len(in))

	for i, v := range in {
		o[i] = mapFieldFromInterface(v)
	}

	return o
}
