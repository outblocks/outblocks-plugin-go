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

		if v, ok := field.LookupCurrentRaw(); ok {
			params = append(params, v)
			continue
		}

		switch input := f.(type) {
		case StringInputField:
			params = append(params, input.Any())
		case BoolInputField:
			params = append(params, input.Any())
		case IntInputField:
			params = append(params, input.Any())
		case MapInputField:
			params = append(params, input.Any())
		case ArrayInputField:
			params = append(params, input.Any())
		default:
			return ""
		}
	}

	return fmt.Sprintf(format, params...)
}
