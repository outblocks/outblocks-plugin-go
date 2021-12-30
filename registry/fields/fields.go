package fields

import (
	"fmt"
	"sync"
)

type ValueTracker interface {
	IsChanged() bool
	IsValid() bool
	Invalidate()
}

type FieldDependencyHolder interface {
	FieldDependencies() []interface{}
}

type Field interface {
	ValueTracker

	Serialize(interface{}) interface{}
	LookupCurrentRaw() (interface{}, bool)
	UnsetCurrent()
	IsOutput() bool
	EmptyValue() interface{}
}

type InputField interface {
	Field

	LookupWantedRaw() (interface{}, bool)
	UnsetWanted()
}

type OutputField interface {
	Field
}

type VerboseField interface {
	Verbose() string
}

type FieldBase struct {
	isOutput                      bool
	currentDefined, wantedDefined bool
	currentVal                    interface{}
	wantedVal                     interface{}
	wantedLazy                    func() interface{}
	wantedFunc                    func() interface{}
	invalidated                   bool

	once struct {
		lazyWanted sync.Once
	}
}

func BasicValue(n interface{}, output bool) FieldBase {
	if output {
		return FieldBase{
			currentDefined: true,
			currentVal:     n,
		}
	}

	return FieldBase{
		wantedDefined: true,
		wantedVal:     n,
	}
}

func BasicValueLazy(f func() interface{}) FieldBase {
	return FieldBase{
		wantedDefined: true,
		wantedLazy:    f,
	}
}

func BasicValueFunc(f func() interface{}) FieldBase {
	return FieldBase{
		wantedDefined: true,
		wantedFunc:    f,
	}
}

func BasicValueUnset(output bool) FieldBase {
	return FieldBase{
		isOutput: output,
	}
}

func (f *FieldBase) UnsetWanted() {
	f.wantedDefined = false
}

func (f *FieldBase) UnsetCurrent() {
	f.currentDefined = false
}

func (f *FieldBase) setCurrent(i interface{}) {
	f.currentDefined = true
	f.invalidated = false
	f.currentVal = i
}

func (f *FieldBase) setWanted(i interface{}) {
	f.wantedDefined = true
	f.wantedLazy = nil
	f.wantedFunc = nil
	f.wantedVal = i
}

func (f *FieldBase) SetWantedLazy(i func() interface{}) {
	f.once.lazyWanted = sync.Once{}
	f.wantedDefined = true
	f.wantedLazy = i
	f.wantedFunc = nil
}

func (f *FieldBase) wanted() interface{} {
	if f.wantedLazy != nil {
		f.once.lazyWanted.Do(func() {
			f.wantedVal = f.wantedLazy()
		})
	}

	if f.wantedFunc != nil {
		return f.wantedFunc()
	}

	return f.wantedVal
}

func (f *FieldBase) IsValid() bool {
	return !f.invalidated
}

func (f *FieldBase) Invalidate() {
	f.invalidated = true
}

func (f *FieldBase) IsChanged() bool {
	if f.invalidated {
		return true
	}

	if !f.wantedDefined {
		return false
	}

	if !f.currentDefined {
		return true
	}

	return f.currentVal != f.wanted()
}

func (f *FieldBase) IsOutput() bool {
	return f.isOutput
}

func (f *FieldBase) IsCurrentDefined() bool {
	return f.currentDefined
}

func (f *FieldBase) IsWantedDefined() bool {
	return f.wantedDefined
}

func (f *FieldBase) lookupAny() (interface{}, bool) {
	if f.currentDefined {
		return f.currentVal, true
	}

	return f.wanted(), f.wantedDefined
}

func (f *FieldBase) LookupCurrentRaw() (interface{}, bool) {
	return f.currentVal, f.currentDefined
}

func (f *FieldBase) LookupWantedRaw() (interface{}, bool) {
	return f.wanted(), f.wantedDefined
}

func (f *FieldBase) Serialize(i interface{}) interface{} {
	return i
}

func toInt(in interface{}) (v int, ok bool) {
	switch i := in.(type) {
	case float64:
		return int(i), true
	case int64:
		return int(i), true
	case int:
		return i, true
	default:
		return 0, false
	}
}

func SetFieldValue(f, v interface{}) error {
	switch val := f.(type) {
	case stringField:
		if _, ok := v.(string); !ok {
			return nil
		}

		val.SetCurrent(v.(string))

	case boolField:
		if _, ok := v.(bool); !ok {
			return nil
		}

		val.SetCurrent(v.(bool))

	case intField:
		out, ok := toInt(v)
		if !ok {
			return nil
		}

		val.SetCurrent(out)

	case mapField:
		if _, ok := v.(map[string]interface{}); !ok {
			return nil
		}

		val.SetCurrent(v.(map[string]interface{}))

	case arrayField:
		if _, ok := v.([]interface{}); !ok {
			return nil
		}

		val.SetCurrent(v.([]interface{}))

	default:
		return fmt.Errorf("unknown field type found: %+v", f)
	}

	return nil
}
