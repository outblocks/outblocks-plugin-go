package registry

import (
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
	Ignored  bool
	ForceNew bool
}

func parseFieldPropertiesTag(tag string) *FieldProperties {
	ret := &FieldProperties{}
	taginfo := strings.Split(tag, ",")

	for _, t := range taginfo {
		switch t {
		case "-":
			ret.Ignored = true

		case "force_new":
			ret.ForceNew = true
		}
	}

	return ret
}