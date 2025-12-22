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
	FieldDependencies() []any
}

type Field interface { //nolint:iface
	ValueTracker

	Serialize(any) any
	LookupCurrentRaw() (any, bool)
	UnsetCurrent()
	IsOutput() bool
	EmptyValue() any
}

type InputField interface {
	Field

	LookupWantedRaw() (any, bool)
	UnsetWanted()
}

type OutputField interface { //nolint:iface
	Field
}

type VerboseField interface {
	Verbose() string
}

type FieldBase struct {
	isOutput                      bool
	currentDefined, wantedDefined bool
	currentVal                    any
	wantedVal                     any
	wantedLazy                    func() any
	wantedFunc                    func() any
	invalidated                   bool

	once struct {
		lazyWanted sync.Once
	}
}

func BasicValue(n any, output bool) FieldBase {
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

func BasicValueLazy(f func() any) FieldBase {
	return FieldBase{
		wantedDefined: true,
		wantedLazy:    f,
	}
}

func BasicValueFunc(f func() any) FieldBase {
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

func (f *FieldBase) setCurrent(i any) {
	f.currentDefined = true
	f.invalidated = false
	f.currentVal = i
}

func (f *FieldBase) setWanted(i any) {
	f.wantedDefined = true
	f.wantedLazy = nil
	f.wantedFunc = nil
	f.wantedVal = i
}

func (f *FieldBase) SetWantedLazy(i func() any) {
	f.once.lazyWanted = sync.Once{}
	f.wantedDefined = true
	f.wantedLazy = i
	f.wantedFunc = nil
}

func (f *FieldBase) wanted() any {
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

func (f *FieldBase) lookupAny() (any, bool) {
	if f.currentDefined {
		return f.currentVal, true
	}

	return f.wanted(), f.wantedDefined
}

func (f *FieldBase) LookupCurrentRaw() (any, bool) {
	return f.currentVal, f.currentDefined
}

func (f *FieldBase) LookupWantedRaw() (any, bool) {
	return f.wanted(), f.wantedDefined
}

func (f *FieldBase) Serialize(i any) any {
	return i
}

func toInt(in any) (v int, ok bool) {
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

func SetFieldValue(f, v any) error {
	switch val := f.(type) {
	case stringField:
		if _, ok := v.(string); !ok {
			return nil
		}

		val.SetCurrent(v.(string)) //nolint:errcheck

	case boolField:
		if _, ok := v.(bool); !ok {
			return nil
		}

		val.SetCurrent(v.(bool)) //nolint:errcheck

	case intField:
		out, ok := toInt(v)
		if !ok {
			return nil
		}

		val.SetCurrent(out)

	case mapField:
		if _, ok := v.(map[string]any); !ok {
			return nil
		}

		val.SetCurrent(v.(map[string]any)) //nolint:errcheck

	case arrayField:
		if _, ok := v.([]any); !ok {
			return nil
		}

		val.SetCurrent(v.([]any)) //nolint:errcheck

	default:
		return fmt.Errorf("unknown field type found: %+v", f)
	}

	return nil
}
