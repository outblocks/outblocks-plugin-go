package fields

import "reflect"

func MakeProxyField(i interface{}) interface{} {
	switch v := i.(type) {
	case stringBaseField:
		return newProxyStringField(v)
	case boolBaseField:
		return newProxyBoolField(v)
	case intBaseField:
		return newProxyIntField(v)
	case mapBaseField:
		return newProxyMapField(v)
	case arrayBaseField:
		return newProxyArrayField(v)
	}

	return nil
}

// String.
type proxyStringField struct {
	StringField

	org stringBaseField
}

func newProxyStringField(org stringBaseField) *proxyStringField {
	return &proxyStringField{
		org: org,
	}
}

func (f *proxyStringField) LookupWanted() (v string, ok bool) {
	return f.org.LookupCurrent()
}

func (f *proxyStringField) Wanted() string {
	return f.org.Current()
}

func (f *proxyStringField) Any() string {
	cur, ok := f.LookupCurrent()
	if ok {
		return cur
	}

	return f.Wanted()
}

func (f *proxyStringField) IsChanged() bool {
	return f.Current() != f.Wanted()
}

func (f *proxyStringField) FieldDependencies() []interface{} {
	if v, ok := f.org.(FieldHolder); ok {
		return v.FieldDependencies()
	}

	return nil
}

// Bool.
type proxyBoolField struct {
	BoolField

	org boolBaseField
}

func newProxyBoolField(org boolBaseField) *proxyBoolField {
	return &proxyBoolField{
		org: org,
	}
}

func (f *proxyBoolField) LookupWanted() (v, ok bool) {
	return f.org.LookupCurrent()
}

func (f *proxyBoolField) Wanted() bool {
	return f.org.Current()
}

func (f *proxyBoolField) Any() bool {
	cur, ok := f.LookupCurrent()
	if ok {
		return cur
	}

	return f.Wanted()
}

func (f *proxyBoolField) IsChanged() bool {
	return f.Current() != f.Wanted()
}

func (f *proxyBoolField) FieldDependencies() []interface{} {
	if v, ok := f.org.(FieldHolder); ok {
		return v.FieldDependencies()
	}

	return nil
}

// Int.
type proxyIntField struct {
	IntField

	org intBaseField
}

func newProxyIntField(org intBaseField) *proxyIntField {
	return &proxyIntField{
		org: org,
	}
}

func (f *proxyIntField) LookupWanted() (v int, ok bool) {
	return f.org.LookupCurrent()
}

func (f *proxyIntField) Wanted() int {
	return f.org.Current()
}

func (f *proxyIntField) Any() int {
	cur, ok := f.LookupCurrent()
	if ok {
		return cur
	}

	return f.Wanted()
}

func (f *proxyIntField) IsChanged() bool {
	return reflect.DeepEqual(f.Current(), f.Wanted())
}

func (f *proxyIntField) FieldDependencies() []interface{} {
	if v, ok := f.org.(FieldHolder); ok {
		return v.FieldDependencies()
	}

	return nil
}

// Map.
type proxyMapField struct {
	MapField

	org mapBaseField
}

func newProxyMapField(org mapBaseField) *proxyMapField {
	return &proxyMapField{
		org: org,
	}
}

func (f *proxyMapField) LookupWanted() (v map[string]interface{}, ok bool) {
	return f.org.LookupCurrent()
}

func (f *proxyMapField) Wanted() map[string]interface{} {
	return f.org.Current()
}

func (f *proxyMapField) Any() map[string]interface{} {
	cur, ok := f.LookupCurrent()
	if ok {
		return cur
	}

	return f.Wanted()
}

func (f *proxyMapField) IsChanged() bool {
	return !reflect.DeepEqual(f.Current(), f.Wanted())
}

func (f *proxyMapField) FieldDependencies() []interface{} {
	if v, ok := f.org.(FieldHolder); ok {
		return v.FieldDependencies()
	}

	return nil
}

// Array.
type proxyArrayField struct {
	ArrayField

	org arrayBaseField
}

func newProxyArrayField(org arrayBaseField) *proxyArrayField {
	return &proxyArrayField{
		org: org,
	}
}

func (f *proxyArrayField) LookupWanted() (v []interface{}, ok bool) {
	return f.org.LookupCurrent()
}

func (f *proxyArrayField) Wanted() []interface{} {
	return f.org.Current()
}

func (f *proxyArrayField) Any() []interface{} {
	cur, ok := f.LookupCurrent()
	if ok {
		return cur
	}

	return f.Wanted()
}

func (f *proxyArrayField) IsChanged() bool {
	return !reflect.DeepEqual(f.Current(), f.Wanted())
}

func (f *proxyArrayField) FieldDependencies() []interface{} {
	if v, ok := f.org.(FieldHolder); ok {
		return v.FieldDependencies()
	}

	return nil
}
