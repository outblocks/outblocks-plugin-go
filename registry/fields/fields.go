package fields

type ValueTracker interface {
	IsChanged() bool
}

type FieldHolder interface {
	FieldDependencies() []interface{}
}

type Field interface {
	ValueTracker

	SerializeValue() interface{}
	LookupCurrentRaw() (interface{}, bool)
	IsOutput() bool
}

type InputField interface {
	Field

	LookupWantedRaw() (interface{}, bool)
	SetWantedAsCurrent()
}

type OutputField interface {
	Field
}

type FieldBase struct {
	isOutput                      bool
	currentDefined, wantedDefined bool
	current, wanted               interface{}
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

func (f *FieldBase) setCurrent(i interface{}) {
	f.currentDefined = true
	f.current = i
}

func (f *FieldBase) IsChanged() bool {
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

func (f *FieldBase) SerializeValue() interface{} {
	return f.current
}

func (f *FieldBase) SetWantedAsCurrent() {
	f.current = f.wanted
	f.currentDefined = f.wantedDefined
}
