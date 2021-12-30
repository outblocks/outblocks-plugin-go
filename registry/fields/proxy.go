package fields

import "reflect"

func MakeProxyField(i interface{}) interface{} {
	switch v := i.(type) {
	case stringField:
		return newProxyStringField(v)
	case boolField:
		return newProxyBoolField(v)
	case intField:
		return newProxyIntField(v)
	case mapField:
		return newProxyMapField(v)
	case arrayField:
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
	StringBaseField
	*proxyBaseField
}

func newProxyStringField(org stringField) *proxyStringField {
	return &proxyStringField{
		proxyBaseField: &proxyBaseField{
			org: org,
		},
	}
}

func (f *proxyStringField) LookupWanted() (v string, ok bool) {
	return f.org.(stringField).LookupCurrent()
}

func (f *proxyStringField) Wanted() string {
	return f.org.(stringField).Current()
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
	BoolBaseField
	*proxyBaseField
}

func newProxyBoolField(org boolField) *proxyBoolField {
	return &proxyBoolField{
		proxyBaseField: &proxyBaseField{
			org: org,
		},
	}
}

func (f *proxyBoolField) LookupWanted() (v, ok bool) {
	return f.org.(boolField).LookupCurrent()
}

func (f *proxyBoolField) Wanted() bool {
	return f.org.(boolField).Current()
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
	IntBaseField
	*proxyBaseField
}

func newProxyIntField(org intField) *proxyIntField {
	return &proxyIntField{
		proxyBaseField: &proxyBaseField{
			org: org,
		},
	}
}

func (f *proxyIntField) LookupWanted() (v int, ok bool) {
	return f.org.(intField).LookupCurrent()
}

func (f *proxyIntField) Wanted() int {
	return f.org.(intField).Current()
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
	MapBaseField
	*proxyBaseField
}

func newProxyMapField(org mapField) *proxyMapField {
	return &proxyMapField{
		proxyBaseField: &proxyBaseField{
			org: org,
		},
	}
}

func (f *proxyMapField) LookupWanted() (v map[string]interface{}, ok bool) {
	return f.org.(mapField).LookupCurrent()
}

func (f *proxyMapField) Wanted() map[string]interface{} {
	return f.org.(mapField).Current()
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
	ArrayBaseField
	*proxyBaseField
}

func newProxyArrayField(org arrayField) *proxyArrayField {
	return &proxyArrayField{
		proxyBaseField: &proxyBaseField{
			org: org,
		},
	}
}

func (f *proxyArrayField) LookupWanted() (v []interface{}, ok bool) {
	return f.org.(arrayField).LookupCurrent()
}

func (f *proxyArrayField) Wanted() []interface{} {
	return f.org.(arrayField).Current()
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
