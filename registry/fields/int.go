package fields

type intBaseField interface {
	SetCurrent(int)
	LookupCurrent() (int, bool)
	Current() int
}

type IntInputField interface {
	intBaseField
	InputField

	LookupWanted() (int, bool)
	Wanted() int
	SetWanted(int)
	Any() int
}

type IntOutputField interface {
	intBaseField
	OutputField
}

type IntField struct {
	FieldBase
}

func Int(val int) IntInputField {
	return &IntField{FieldBase: BasicValue(val, false)}
}

func IntUnset() IntInputField {
	return &IntField{FieldBase: BasicValueUnset(false)}
}

func IntUnsetOutput() IntOutputField {
	return &IntField{FieldBase: BasicValueUnset(true)}
}

func IntOutput(val int) IntOutputField {
	return &IntField{FieldBase: BasicValue(val, true)}
}

func (f *IntField) SetCurrent(i int) {
	f.setCurrent(i)
}

func (f *IntField) LookupCurrent() (v int, ok bool) {
	if !f.currentDefined {
		return 0, f.currentDefined
	}

	return f.current.(int), true
}

func (f *IntField) SetWanted(i int) {
	f.setWanted(i)
}

func (f *IntField) LookupWanted() (v int, ok bool) {
	if !f.wantedDefined {
		return 0, false
	}

	return f.wanted.(int), true
}

func (f *IntField) Wanted() int {
	v, _ := f.LookupWanted()
	return v
}

func (f *IntField) Current() int {
	v, _ := f.LookupCurrent()
	return v
}

func (f *IntField) Any() int {
	any, defined := f.lookupAny()
	if !defined {
		return 0
	}

	return any.(int)
}
