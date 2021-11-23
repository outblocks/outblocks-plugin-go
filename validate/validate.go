package validate

import (
	"github.com/outblocks/outblocks-plugin-go/types"
	"google.golang.org/protobuf/types/known/structpb"
)

func Any(m map[string]*structpb.Value, key, msg string) (val interface{}, err error) {
	if v, ok := m[key]; ok {
		return v.AsInterface(), nil
	}

	return nil, types.NewValidationError(key, msg)
}

func String(m map[string]*structpb.Value, key, msg string) (val string, err error) {
	v, err := Any(m, key, msg)
	if err != nil {
		return "", err
	}

	if v, ok := v.(string); ok {
		return v, nil
	}

	return "", types.NewValidationError(key, msg)
}

func OptionalString(def string, m map[string]*structpb.Value, key, msg string) (val string, err error) {
	if _, ok := m[key]; !ok {
		return def, nil
	}

	return String(m, key, msg)
}
