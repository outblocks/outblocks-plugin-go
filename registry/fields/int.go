package fields

type intField interface {
	SetCurrent(int)
	LookupCurrent() (int, bool)
	Current() int
}

type IntInputField interface {
	intField
	InputField

	LookupWanted() (int, bool)
	Wanted() int
	SetWanted(int)
	Any() int
}

type IntOutputField interface {
	intField
	OutputField

	Input() IntInputField
}

type IntBaseField struct {
	FieldBase
}

func Int(val int) IntInputField {
	return &IntBaseField{FieldBase: BasicValue(val, false)}
}

func IntUnset() IntInputField {
	return &IntBaseField{FieldBase: BasicValueUnset(false)}
}

func IntUnsetOutput() IntOutputField {
	return &IntBaseField{FieldBase: BasicValueUnset(true)}
}

func IntOutput(val int) IntOutputField {
	return &IntBaseField{FieldBase: BasicValue(val, true)}
}

func (f *IntBaseField) SetCurrent(i int) {
	f.setCurrent(i)
}

func (f *IntBaseField) LookupCurrent() (v int, ok bool) {
	if !f.currentDefined {
		return 0, f.currentDefined
	}

	return f.currentVal.(int), true
}

func (f *IntBaseField) SetWanted(i int) {
	f.setWanted(i)
}

func (f *IntBaseField) LookupWanted() (v int, ok bool) {
	if !f.wantedDefined {
		return 0, false
	}

	return f.wanted().(int), true
}

func (f *IntBaseField) Wanted() int {
	v, _ := f.LookupWanted()
	return v
}

func (f *IntBaseField) Current() int {
	v, _ := f.LookupCurrent()
	return v
}

func (f *IntBaseField) Any() int {
	any, defined := f.lookupAny()
	if !defined {
		return 0
	}

	return any.(int)
}

func (f *IntBaseField) Input() IntInputField {
	return f
}

func (f *IntBaseField) EmptyValue() interface{} {
	return 0
}
