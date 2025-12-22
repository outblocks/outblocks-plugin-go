package util

import (
	"google.golang.org/protobuf/types/known/structpb"
)

func MergeMaps(a ...map[string]any) map[string]any {
	if len(a) == 0 {
		return nil
	}

	out := make(map[string]any, len(a[0]))
	for k, v := range a[0] {
		out[k] = v
	}

	for _, m := range a[1:] {
		for k, v := range m {
			if v, ok := v.(map[string]any); ok {
				if bv, ok := out[k]; ok {
					if bv, ok := bv.(map[string]any); ok {
						out[k] = MergeMaps(bv, v)
						continue
					}
				}
			}

			out[k] = v
		}
	}

	return out
}

func MergeStringMaps(a ...map[string]string) map[string]string {
	if len(a) == 0 {
		return nil
	}

	out := make(map[string]string, len(a[0]))

	for _, m := range a {
		for k, v := range m {
			out[k] = v
		}
	}

	return out
}

func MustNewStruct(m map[string]any) *structpb.Struct {
	s, err := structpb.NewStruct(m)
	if err != nil {
		panic(err)
	}

	return s
}
