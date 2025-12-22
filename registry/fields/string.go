package fields

import (
	"fmt"
	"strings"

	"github.com/outblocks/outblocks-plugin-go/util"
)

type stringField interface {
	SetCurrent(string)
	LookupCurrent() (string, bool)
	Current() string
}

type StringInputField interface {
	stringField
	InputField

	LookupWanted() (string, bool)
	Wanted() string
	SetWanted(string)
	Any() string
}

type StringOutputField interface {
	stringField
	OutputField

	Input() StringInputField
}

type StringBaseField struct {
	FieldBase
}

func String(val string) StringInputField {
	return &StringBaseField{FieldBase: BasicValue(val, false)}
}

func StringUnset() StringInputField {
	return &StringBaseField{FieldBase: BasicValueUnset(false)}
}

func StringUnsetOutput() StringOutputField {
	return &StringBaseField{FieldBase: BasicValueUnset(true)}
}

func StringOutput(val string) StringOutputField {
	return &StringBaseField{FieldBase: BasicValue(val, true)}
}

func (f *StringBaseField) SetCurrent(i string) {
	f.setCurrent(i)
}

func (f *StringBaseField) LookupCurrent() (v string, ok bool) {
	if !f.currentDefined {
		return "", false
	}

	val, ok := f.currentVal.(string)

	return val, ok
}

func (f *StringBaseField) SetWanted(i string) {
	f.setWanted(i)
}

func (f *StringBaseField) LookupWanted() (v string, ok bool) {
	if !f.wantedDefined {
		return "", false
	}

	val, ok := f.wanted().(string)

	return val, ok
}

func (f *StringBaseField) Wanted() string {
	v, _ := f.LookupWanted()
	return v
}

func (f *StringBaseField) Current() string {
	v, _ := f.LookupCurrent()
	return v
}

func (f *StringBaseField) Any() string {
	val, defined := f.lookupAny()
	if !defined {
		return ""
	}

	return val.(string) //nolint:errcheck
}

func (f *StringBaseField) Input() StringInputField {
	return f
}

func (f *StringBaseField) EmptyValue() any {
	return ""
}

type IStringBaseField struct {
	StringBaseField
}

func (f *IStringBaseField) IsChanged() bool {
	if f.currentVal == nil || f.wanted() == nil || f.invalidated {
		return f.StringBaseField.IsChanged()
	}

	return strings.EqualFold(f.currentVal.(string), f.wanted().(string)) //nolint:errcheck
}

func IString(val string) StringInputField {
	return &IStringBaseField{StringBaseField: StringBaseField{BasicValue(val, false)}}
}

type SprintfBaseField struct {
	FieldBase
	args []any
	fmt  string
}

func Sprintf(format string, args ...any) StringInputField {
	return &SprintfBaseField{
		FieldBase: BasicValue(nil, false),
		fmt:       format,
		args:      args,
	}
}

func (f *SprintfBaseField) Any() string {
	if f.currentDefined {
		return f.currentVal.(string) //nolint:errcheck
	}

	var args []any

	for _, a := range f.args {
		v, ok := a.(Field)
		if !ok {
			args = append(args, a)
			continue
		}

		if !v.IsOutput() {
			a, ok = v.(InputField).LookupWantedRaw() //nolint:errcheck
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

func (f *SprintfBaseField) LookupWanted() (string, bool) {
	if !f.wantedDefined {
		return "", false
	}

	var args []any

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

func (f *SprintfBaseField) SetWanted(i string) {
	f.setWanted(i)
}

func (f *SprintfBaseField) LookupWantedRaw() (any, bool) {
	return f.LookupWanted()
}

func (f *SprintfBaseField) Wanted() string {
	v, _ := f.LookupWanted()
	return v
}

func (f *SprintfBaseField) Current() string {
	v, _ := f.LookupCurrent()
	return v
}

func (f *SprintfBaseField) FieldDependencies() []any {
	var deps []any

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

func (f *SprintfBaseField) LookupCurrent() (v string, ok bool) {
	if f.currentDefined {
		return f.currentVal.(string), true //nolint:errcheck
	}

	var args []any

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

func (f *SprintfBaseField) LookupCurrentRaw() (v any, ok bool) {
	return f.LookupCurrent()
}

func (f *SprintfBaseField) SetCurrent(i string) {
	f.setCurrent(i)
}

func (f *SprintfBaseField) IsChanged() bool {
	return f.Current() != f.Wanted()
}

func (f *SprintfBaseField) EmptyValue() any {
	return ""
}

// String field with lazy initialization. If current value is defined in state, it is used instead.

type LazyStringBaseField struct {
	StringBaseField

	newValue func() any
}

func (f *LazyStringBaseField) SetCurrent(i string) {
	f.StringBaseField.SetCurrent(i)
	f.SetWanted(i)
}

func (f *LazyStringBaseField) UnsetCurrent() {
	f.StringBaseField.UnsetCurrent()
	f.SetWantedLazy(f.newValue)
}

func LazyString(newValue func() string) StringInputField {
	f := &LazyStringBaseField{
		newValue: func() any {
			return newValue()
		},
	}
	f.SetWantedLazy(f.newValue)

	return f
}

// Random string field with lazy initialization. If current value is defined in state, it is used instead.

type RandomStringField struct {
	LazyStringBaseField

	prefix, suffix                 string
	lower, upper, numeric, special bool
	length                         int
}

func (f *RandomStringField) Verbose() string {
	if cur, ok := f.LookupCurrent(); ok {
		return cur
	}

	return fmt.Sprintf("%s*%s", f.prefix, f.suffix)
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
	f.newValue = func() any {
		return f.prefix + util.RandomStringCustom(f.lower, f.upper, f.numeric, f.special, f.length) + f.suffix
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

func VerboseString(f stringField) string {
	if vf, ok := f.(VerboseField); ok {
		return vf.Verbose()
	}

	if sif, ok := f.(StringInputField); ok {
		return sif.Any()
	}

	return f.Current()
}
