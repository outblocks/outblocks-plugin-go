package fields

type IntInputField interface {
	InputField

	SetCurrent(int)
	LookupCurrent() (int, bool)
	LookupWanted() (int, bool)
	GetCurrent() int
	GetWanted() int
	GetAny() int
}

type IntOutputField interface {
	OutputField

	SetCurrent(int)
	LookupCurrent() (int, bool)
	GetCurrent() int
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

func (f *IntField) LookupWanted() (v int, ok bool) {
	if !f.wantedDefined {
		return 0, false
	}

	return f.wanted.(int), true
}

func (f *IntField) GetWanted() int {
	v, _ := f.LookupWanted()
	return v
}

func (f *IntField) GetCurrent() int {
	v, _ := f.LookupCurrent()
	return v
}

func (f *IntField) GetAny() int {
	any, defined := f.lookupAny()
	if !defined {
		return 0
	}

	return any.(int)
}
