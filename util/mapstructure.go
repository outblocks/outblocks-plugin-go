package util

import (
	"reflect"

	"github.com/mitchellh/mapstructure"
)

type Unmarshaler interface {
	UnmarshalMapstructure(interface{}) error
}

func MapstructureJSONDecode(in, out interface{}) error {
	return MapstructureDecode(in, out, "json")
}

func MapstructureDecode(in, out interface{}, tag string) error {
	cfg := &mapstructure.DecoderConfig{
		DecodeHook: func(from reflect.Value, to reflect.Value) (interface{}, error) {
			// If the destination implements the unmarshaling interface
			u, ok := to.Interface().(Unmarshaler)
			if !ok {
				return from.Interface(), nil
			}

			// If it is nil and a pointer, create and assign the target value first
			if to.IsNil() && to.Type().Kind() == reflect.Ptr {
				to.Set(reflect.New(to.Type().Elem()))
				u = to.Interface().(Unmarshaler)
			}

			// Call the custom unmarshaling method
			if err := u.UnmarshalMapstructure(from.Interface()); err != nil {
				return to.Interface(), err
			}

			return to.Interface(), nil
		},
		Metadata:         nil,
		Result:           out,
		TagName:          tag,
		WeaklyTypedInput: true,
	}

	decoder, err := mapstructure.NewDecoder(cfg)
	if err != nil {
		return err
	}

	return decoder.Decode(in)
}
