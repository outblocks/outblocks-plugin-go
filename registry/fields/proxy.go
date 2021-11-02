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
	default:
		panic("unknown base field for proxy")
	}
}

type proxyBaseField struct {
	org interface{}
}

func (f *proxyBaseField) FieldDependencies() []interface{} {
	if v, ok := f.org.(FieldDependencyHolder); ok {
		return append(v.FieldDependencies(), f.org)
	}

	return []interface{}{f.org}
}

// String.
type proxyStringField struct {
	StringField
	*proxyBaseField
}

func newProxyStringField(org stringBaseField) *proxyStringField {
	return &proxyStringField{
		proxyBaseField: &proxyBaseField{
			org: org,
		},
	}
}

func (f *proxyStringField) LookupWanted() (v string, ok bool) {
	return f.org.(stringBaseField).LookupCurrent()
}

func (f *proxyStringField) Wanted() string {
	return f.org.(stringBaseField).Current()
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

// Bool.
type proxyBoolField struct {
	BoolField
	*proxyBaseField
}

func newProxyBoolField(org boolBaseField) *proxyBoolField {
	return &proxyBoolField{
		proxyBaseField: &proxyBaseField{
			org: org,
		},
	}
}

func (f *proxyBoolField) LookupWanted() (v, ok bool) {
	return f.org.(boolBaseField).LookupCurrent()
}

func (f *proxyBoolField) Wanted() bool {
	return f.org.(boolBaseField).Current()
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

// Int.
type proxyIntField struct {
	IntField
	*proxyBaseField
}

func newProxyIntField(org intBaseField) *proxyIntField {
	return &proxyIntField{
		proxyBaseField: &proxyBaseField{
			org: org,
		},
	}
}

func (f *proxyIntField) LookupWanted() (v int, ok bool) {
	return f.org.(intBaseField).LookupCurrent()
}

func (f *proxyIntField) Wanted() int {
	return f.org.(intBaseField).Current()
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

// Map.
type proxyMapField struct {
	MapField
	*proxyBaseField
}

func newProxyMapField(org mapBaseField) *proxyMapField {
	return &proxyMapField{
		proxyBaseField: &proxyBaseField{
			org: org,
		},
	}
}

func (f *proxyMapField) LookupWanted() (v map[string]interface{}, ok bool) {
	return f.org.(mapBaseField).LookupCurrent()
}

func (f *proxyMapField) Wanted() map[string]interface{} {
	return f.org.(mapBaseField).Current()
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

// Array.
type proxyArrayField struct {
	ArrayField
	*proxyBaseField
}

func newProxyArrayField(org arrayBaseField) *proxyArrayField {
	return &proxyArrayField{
		proxyBaseField: &proxyBaseField{
			org: org,
		},
	}
}

func (f *proxyArrayField) LookupWanted() (v []interface{}, ok bool) {
	return f.org.(arrayBaseField).LookupCurrent()
}

func (f *proxyArrayField) Wanted() []interface{} {
	return f.org.(arrayBaseField).Current()
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
