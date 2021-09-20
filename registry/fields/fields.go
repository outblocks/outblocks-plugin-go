package fields

import "fmt"

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
}

type InputField interface {
	Field

	LookupWantedRaw() (interface{}, bool)
	UnsetWanted()
}

type OutputField interface {
	Field
}

type FieldBase struct {
	isOutput                      bool
	currentDefined, wantedDefined bool
	current, wanted               interface{}
	invalidated                   bool
}

func BasicValue(n interface{}, output bool) FieldBase {
	if output {
		return FieldBase{
			currentDefined: true,
			current:        n,
		}
	}

	return FieldBase{
		wantedDefined: true,
		wanted:        n,
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
	f.current = i
}

func (f *FieldBase) setWanted(i interface{}) {
	f.wantedDefined = true
	f.wanted = i
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

	return f.current != f.wanted
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
		return f.current, true
	}

	return f.wanted, f.wantedDefined
}

func (f *FieldBase) LookupCurrentRaw() (interface{}, bool) {
	return f.current, f.currentDefined
}

func (f *FieldBase) LookupWantedRaw() (interface{}, bool) {
	return f.wanted, f.wantedDefined
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
	case stringBaseField:
		if _, ok := v.(string); !ok {
			return nil
		}

		val.SetCurrent(v.(string))

	case boolBaseField:
		if _, ok := v.(bool); !ok {
			return nil
		}

		val.SetCurrent(v.(bool))

	case intBaseField:
		out, ok := toInt(v)
		if !ok {
			return nil
		}

		val.SetCurrent(out)

	case mapBaseField:
		if _, ok := v.(map[string]interface{}); !ok {
			return nil
		}

		val.SetCurrent(v.(map[string]interface{}))

	case arrayBaseField:
		if _, ok := v.([]interface{}); !ok {
			return nil
		}

		val.SetCurrent(v.([]interface{}))

	default:
		return fmt.Errorf("unknown field type found: %+v", f)
	}

	return nil
}
