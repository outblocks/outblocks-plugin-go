package validate

import plugin "github.com/outblocks/outblocks-plugin-go"

func ValidateAny(m map[string]interface{}, key, msg string) (plugin.Response, interface{}) {
	if v, ok := m[key]; ok {
		return nil, v
	}

	return &plugin.ValidationErrorResponse{
		Path:  key,
		Error: msg,
	}, nil

}

func ValidateString(m map[string]interface{}, key, msg string) (plugin.Response, string) {
	res, v := ValidateAny(m, key, msg)
	if res != nil {
		return res, ""
	}

	if v, ok := v.(string); ok {
		return nil, v
	}

	return &plugin.ValidationErrorResponse{
		Path:  key,
		Error: msg,
	}, ""
}
