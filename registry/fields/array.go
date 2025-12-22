package fields

import (
	"reflect"
)

type arrayField interface {
	SetCurrent([]any)
	LookupCurrent() ([]any, bool)
	Current() []any
}

type ArrayInputField interface {
	arrayField
	InputField

	LookupWanted() ([]any, bool)
	Wanted() []any
	SetWanted([]any)
	Any() []any
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

func (f *ArrayBaseField) SetCurrent(i []any) {
	f.setCurrent(interfaceArrayToFieldArray(i))
}

func (f *ArrayBaseField) LookupCurrent() (v []any, ok bool) {
	if !f.currentDefined {
		return nil, f.currentDefined
	}

	val, ok := f.Serialize(f.currentVal).([]any)

	return val, ok
}

func (f *ArrayBaseField) SetWanted(i []any) {
	f.setWanted(interfaceArrayToFieldArray(i))
}

func (f *ArrayBaseField) LookupWanted() (v []any, ok bool) {
	if !f.wantedDefined {
		return nil, false
	}

	val, ok := f.Serialize(f.wanted()).([]any)

	return val, ok
}

func (f *ArrayBaseField) Wanted() []any {
	v, _ := f.LookupWanted()
	return v
}

func (f *ArrayBaseField) Current() []any {
	v, _ := f.LookupCurrent()
	return v
}

func (f *ArrayBaseField) Any() []any {
	cur, ok := f.LookupCurrent()
	if ok {
		return cur
	}

	return f.Wanted()
}

func (f *ArrayBaseField) Serialize(i any) any {
	if i == nil {
		return make([]any, 0)
	}

	c := i.([]Field) //nolint:errcheck
	m := make([]any, len(c))

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

func (f *ArrayBaseField) FieldDependencies() []any {
	if f.wanted() == nil {
		return nil
	}

	var deps []any

	for _, v := range f.wanted().([]Field) { //nolint:errcheck
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

func (f *ArrayBaseField) EmptyValue() any {
	var ret []any
	return ret
}

func interfaceArrayToFieldArray(in []any) []Field {
	o := make([]Field, len(in))

	for i, v := range in {
		o[i] = mapFieldFromInterface(v)
	}

	return o
}
