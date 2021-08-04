package fields

import (
	"fmt"
	"strings"
)

type stringBaseField interface {
	SetCurrent(string)
	LookupCurrent() (string, bool)
	Current() string
}

type StringInputField interface {
	stringBaseField
	InputField

	LookupWanted() (string, bool)
	Wanted() string
	Any() string
}

type StringOutputField interface {
	stringBaseField
	OutputField
}

type StringField struct {
	FieldBase
}

func String(val string) StringInputField {
	return &StringField{FieldBase: BasicValue(val, false)}
}

func StringUnset() StringInputField {
	return &StringField{FieldBase: BasicValueUnset(false)}
}

func StringUnsetOutput() StringOutputField {
	return &StringField{FieldBase: BasicValueUnset(true)}
}

func StringOutput(val string) StringOutputField {
	return &StringField{FieldBase: BasicValue(val, true)}
}

func (f *StringField) SetCurrent(i string) {
	f.setCurrent(i)
}

func (f *StringField) LookupCurrent() (v string, ok bool) {
	if !f.currentDefined {
		return "", false
	}

	return f.current.(string), true
}

func (f *StringField) LookupWanted() (v string, ok bool) {
	if !f.wantedDefined {
		return "", false
	}

	return f.wanted.(string), true
}

func (f *StringField) Wanted() string {
	v, _ := f.LookupWanted()
	return v
}

func (f *StringField) Current() string {
	v, _ := f.LookupCurrent()
	return v
}

func (f *StringField) Any() string {
	any, defined := f.lookupAny()
	if !defined {
		return ""
	}

	return any.(string)
}

type IStringField struct {
	StringField
}

func (f *IStringField) IsChanged() bool {
	if f.current == nil || f.wanted == nil {
		return f.StringField.IsChanged()
	}

	return strings.EqualFold(f.current.(string), f.wanted.(string))
}

func IString(val string) StringInputField {
	return &IStringField{StringField: StringField{BasicValue(val, false)}}
}

type SprintfField struct {
	FieldBase
	args []interface{}
	fmt  string
}

func Sprintf(format string, args ...interface{}) StringInputField {
	return &SprintfField{
		FieldBase: BasicValue(nil, false),
		fmt:       format,
		args:      args,
	}
}

func (f *SprintfField) Any() string {
	if f.currentDefined {
		return f.current.(string)
	}

	var args []interface{}

	for _, a := range f.args {
		v, ok := a.(Field)
		if !ok {
			args = append(args, a)
			continue
		}

		if !v.IsOutput() {
			a, ok = v.(InputField).LookupWantedRaw()
			if ok {
				args = append(args, a)
				continue
			}
		}

		a, ok = v.LookupCurrentRaw()
		if !ok {
			panic("cannot get value of argument")
		}

		args = append(args, a)
	}

	return fmt.Sprintf(f.fmt, args...)
}

func (f *SprintfField) LookupWanted() (string, bool) {
	if !f.wantedDefined {
		return "", false
	}

	var args []interface{}

	for _, a := range f.args {
		v, ok := a.(InputField)
		if !ok {
			args = append(args, a)
			continue
		}

		if !v.IsOutput() {
			a, ok = v.LookupWantedRaw()
			if ok {
				args = append(args, a)
				continue
			}
		}

		a, ok = v.LookupCurrentRaw()
		if !ok {
			return "", false
		}

		args = append(args, a)
	}

	return fmt.Sprintf(f.fmt, args...), true
}

func (f *SprintfField) LookupWantedRaw() (interface{}, bool) {
	return f.LookupWanted()
}

func (f *SprintfField) Wanted() string {
	v, _ := f.LookupWanted()
	return v
}

func (f *SprintfField) Current() string {
	v, _ := f.LookupCurrent()
	return v
}

func (f *SprintfField) FieldDependencies() []interface{} {
	var deps []interface{}

	for _, a := range f.args {
		_, ok := a.(Field)
		if ok {
			if fh, ok := a.(FieldDependencyHolder); ok {
				deps = append(deps, fh.FieldDependencies()...)

				continue
			}

			deps = append(deps, a)
		}
	}

	return deps
}

func (f *SprintfField) LookupCurrent() (v string, ok bool) {
	if f.currentDefined {
		return f.current.(string), true
	}

	var args []interface{}

	for _, a := range f.args {
		v, ok := a.(InputField)
		if !ok {
			args = append(args, a)
			continue
		}

		a, ok = v.LookupCurrentRaw()
		if !ok {
			return "", false
		}

		args = append(args, a)
	}

	return fmt.Sprintf(f.fmt, args...), true
}

func (f *SprintfField) LookupCurrentRaw() (v interface{}, ok bool) {
	return f.LookupCurrent()
}

func (f *SprintfField) SetCurrent(i string) {
	f.setCurrent(i)
}

func (f *SprintfField) IsChanged() bool {
	return f.Current() != f.Wanted()
}
