package fields

import (
	"fmt"
	"strings"

	"github.com/outblocks/outblocks-plugin-go/util"
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
	SetWanted(string)
	Any() string
}

type StringOutputField interface {
	stringBaseField
	OutputField

	Input() StringInputField
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

	return f.currentVal.(string), true
}

func (f *StringField) SetWanted(i string) {
	f.setWanted(i)
}

func (f *StringField) LookupWanted() (v string, ok bool) {
	if !f.wantedDefined {
		return "", false
	}

	return f.wanted().(string), true
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

func (f *StringField) Input() StringInputField {
	return f
}

func (f *StringField) EmptyValue() interface{} {
	return ""
}

type IStringField struct {
	StringField
}

func (f *IStringField) IsChanged() bool {
	if f.currentVal == nil || f.wanted() == nil || f.invalidated {
		return f.StringField.IsChanged()
	}

	return strings.EqualFold(f.currentVal.(string), f.wanted().(string))
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
		return f.currentVal.(string)
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
		if !ok || !v.IsValid() {
			a = v.EmptyValue()
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
		if !ok || !v.IsValid() {
			return "", false
		}

		args = append(args, a)
	}

	return fmt.Sprintf(f.fmt, args...), true
}

func (f *SprintfField) SetWanted(i string) {
	f.setWanted(i)
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
		return f.currentVal.(string), true
	}

	var args []interface{}

	for _, a := range f.args {
		v, ok := a.(InputField)
		if !ok {
			args = append(args, a)
			continue
		}

		a, ok = v.LookupCurrentRaw()
		if !ok || !v.IsValid() {
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

func (f *SprintfField) EmptyValue() interface{} {
	return ""
}

// Random string field with lazy initialization. If current value is defined in state, it is used instead.

type RandomStringField struct {
	StringField

	prefix, suffix                 string
	lower, upper, numeric, special bool
	length                         int
}

func (f *RandomStringField) SetCurrent(i string) {
	f.StringField.SetCurrent(i)
	f.SetWanted(i)
}

func (f *RandomStringField) UnsetCurrent() {
	f.StringField.UnsetCurrent()
	f.SetWantedLazy(f.newValue)
}

func (f *RandomStringField) newValue() interface{} {
	return f.prefix + util.RandomStringCustom(f.lower, f.upper, f.numeric, f.special, f.length) + f.suffix
}

func (f *RandomStringField) IsChanged() bool {
	if f.StringField.IsChanged() {
		return true
	}

	cur := f.Current()
	if !strings.HasSuffix(cur, f.suffix) || !strings.HasPrefix(cur, f.prefix) {
		return true
	}

	return false
}

func randomString(prefix, suffix string, lower, upper, numeric, special bool, length int) StringInputField {
	f := &RandomStringField{
		prefix:  prefix,
		suffix:  suffix,
		lower:   lower,
		upper:   upper,
		numeric: numeric,
		special: special,
		length:  length,
	}
	f.SetWantedLazy(f.newValue)

	return f
}

func RandomString(lower, upper, numeric, special bool, length int) StringInputField {
	return randomString("", "", lower, upper, numeric, special, length)
}

func RandomStringWithPrefix(prefix string, lower, upper, numeric, special bool, length int) StringInputField {
	return randomString(prefix, "", lower, upper, numeric, special, length)
}

func RandomStringWithSuffix(suffix string, lower, upper, numeric, special bool, length int) StringInputField {
	return randomString("", suffix, lower, upper, numeric, special, length)
}
