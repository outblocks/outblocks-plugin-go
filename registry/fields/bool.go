package fields

type boolField interface {
	SetCurrent(bool)
	LookupCurrent() (bool, bool)
	Current() bool
}

type BoolInputField interface {
	boolField
	InputField

	LookupWanted() (bool, bool)
	Wanted() bool
	SetWanted(bool)
	Any() bool
}

type BoolOutputField interface {
	boolField
	OutputField

	Input() BoolInputField
}

type BoolBaseField struct {
	FieldBase
}

func Bool(val bool) BoolInputField {
	return &BoolBaseField{FieldBase: BasicValue(val, false)}
}

func BoolUnset() BoolInputField {
	return &BoolBaseField{FieldBase: BasicValueUnset(false)}
}

func BoolUnsetOutput() BoolOutputField {
	return &BoolBaseField{FieldBase: BasicValueUnset(true)}
}

func BoolOutput(val bool) BoolOutputField {
	return &BoolBaseField{FieldBase: BasicValue(val, true)}
}

func (f *BoolBaseField) SetCurrent(i bool) {
	f.setCurrent(i)
}

func (f *BoolBaseField) LookupCurrent() (v, ok bool) {
	if !f.currentDefined {
		return false, f.currentDefined
	}

	return f.currentVal.(bool), true
}

func (f *BoolBaseField) SetWanted(i bool) {
	f.setWanted(i)
}

func (f *BoolBaseField) LookupWanted() (v, ok bool) {
	if !f.wantedDefined {
		return false, false
	}

	return f.wanted().(bool), true
}

func (f *BoolBaseField) Wanted() bool {
	v, _ := f.LookupWanted()
	return v
}

func (f *BoolBaseField) Current() bool {
	v, _ := f.LookupCurrent()
	return v
}

func (f *BoolBaseField) Any() bool {
	any, defined := f.lookupAny()
	if !defined {
		return false
	}

	return any.(bool)
}

func (f *BoolBaseField) Input() BoolInputField {
	return f
}

func (f *BoolBaseField) EmptyValue() interface{} {
	return false
}
