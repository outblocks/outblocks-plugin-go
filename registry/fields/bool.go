package fields

type BoolInputField interface {
	InputField

	SetCurrent(bool)
	LookupCurrent() (bool, bool)
	LookupWanted() (bool, bool)
	Current() bool
	Wanted() bool
	Any() bool
}

type BoolOutputField interface {
	OutputField

	SetCurrent(bool)
	LookupCurrent() (bool, bool)
	Current() bool
}

type BoolField struct {
	FieldBase
}

func Bool(val bool) BoolInputField {
	return &BoolField{FieldBase: BasicValue(val, false)}
}

func BoolUnset() BoolInputField {
	return &BoolField{FieldBase: BasicValueUnset(false)}
}

func BoolUnsetOutput() BoolOutputField {
	return &BoolField{FieldBase: BasicValueUnset(true)}
}

func BoolOutput(val bool) BoolOutputField {
	return &BoolField{FieldBase: BasicValue(val, true)}
}

func (f *BoolField) SetCurrent(i bool) {
	f.setCurrent(i)
}

func (f *BoolField) LookupCurrent() (v, ok bool) {
	if !f.currentDefined {
		return false, f.currentDefined
	}

	return f.current.(bool), true
}

func (f *BoolField) LookupWanted() (v, ok bool) {
	if !f.wantedDefined {
		return false, false
	}

	return f.wanted.(bool), true
}

func (f *BoolField) Wanted() bool {
	v, _ := f.LookupWanted()
	return v
}

func (f *BoolField) Current() bool {
	v, _ := f.LookupCurrent()
	return v
}

func (f *BoolField) Any() bool {
	any, defined := f.lookupAny()
	if !defined {
		return false
	}

	return any.(bool)
}
