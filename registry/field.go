package registry

import (
	"fmt"
	"reflect"
	"strings"
)

type FieldTypeInfo struct {
	ReflectType reflect.StructField
	Properties  *FieldProperties
	Default     string
	DefaultSet  bool
}

type FieldInfo struct {
	Type  *FieldTypeInfo
	Value reflect.Value
}

type FieldProperties struct {
	Ignored           bool // ignore from state
	ForceNew          bool // any change of this field forces new resource
	Computed          bool // computed field disallows user input and is created by resource itself
	PropagateRecreate bool // two step update fields support more than 1 update (as they check dependencies if they are recreated)
}

func parseFieldPropertiesTag(tag string) *FieldProperties {
	ret := &FieldProperties{}

	if tag == "" {
		return ret
	}

	taginfo := strings.Split(tag, ",")

	for _, t := range taginfo {
		switch t {
		case "-":
			ret.Ignored = true

		case "computed":
			ret.Computed = true

		case "force_new":
			ret.ForceNew = true

		case "propagate_recreate":
			ret.PropagateRecreate = true
		default:
			panic(fmt.Sprintf("unknown field properties tag: %s", t))
		}
	}

	return ret
}
