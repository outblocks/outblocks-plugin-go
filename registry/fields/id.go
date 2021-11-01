package fields

import (
	"fmt"
)

func GenerateID(format string, a ...interface{}) string {
	var params []interface{}

	for _, f := range a {
		field, ok := f.(Field)
		if !ok {
			params = append(params, f)
			continue
		}

		if !field.IsValid() {
			return ""
		}

		if v, ok2 := field.LookupCurrentRaw(); ok2 {
			params = append(params, v)
			continue
		}

		var param interface{}

		switch input := f.(type) {
		case StringInputField:
			param, ok = input.LookupWanted()
		case BoolInputField:
			param, ok = input.LookupWanted()
		case IntInputField:
			param, ok = input.LookupWanted()
		case MapInputField:
			param, ok = input.LookupWanted()
		case ArrayInputField:
			param, ok = input.LookupWanted()
		default:
			return ""
		}

		if !ok {
			return ""
		}

		params = append(params, param)
	}

	return fmt.Sprintf(format, params...)
}
