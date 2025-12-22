package fields

import "reflect"

func MakeProxyField(i any) any {
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
	org any
}

func (f *proxyBaseField) FieldDependencies() []any {
	if v, ok := f.org.(FieldDependencyHolder); ok {
		return append(v.FieldDependencies(), f.org)
	}

	return []any{f.org}
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
	val, ok := f.org.(stringField)
	if !ok {
		return "", false
	}

	return val.LookupCurrent()
}

func (f *proxyStringField) Wanted() string {
	val, ok := f.LookupWanted()
	if !ok {
		return ""
	}

	return val
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
	return f.org.(boolField).LookupCurrent() //nolint:errcheck
}

func (f *proxyBoolField) Wanted() bool {
	return f.org.(boolField).Current() //nolint:errcheck
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
	return f.org.(intField).LookupCurrent() //nolint:errcheck
}

func (f *proxyIntField) Wanted() int {
	return f.org.(intField).Current() //nolint:errcheck
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

func (f *proxyMapField) WantedFieldMap() map[string]Field {
	return f.org.(mapField).CurrentFieldMap() //nolint:errcheck
}

func (f *proxyMapField) LookupWanted() (v map[string]any, ok bool) {
	return f.org.(mapField).LookupCurrent() //nolint:errcheck
}

func (f *proxyMapField) Wanted() map[string]any {
	return f.org.(mapField).Current() //nolint:errcheck
}

func (f *proxyMapField) Any() map[string]any {
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

func (f *proxyArrayField) LookupWanted() (v []any, ok bool) {
	return f.org.(arrayField).LookupCurrent() //nolint:errcheck
}

func (f *proxyArrayField) Wanted() []any {
	return f.org.(arrayField).Current() //nolint:errcheck
}

func (f *proxyArrayField) Any() []any {
	cur, ok := f.LookupCurrent()
	if ok {
		return cur
	}

	return f.Wanted()
}

func (f *proxyArrayField) IsChanged() bool {
	return !reflect.DeepEqual(f.Current(), f.Wanted())
}
