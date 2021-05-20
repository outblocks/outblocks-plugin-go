package validate

import plugin "github.com/outblocks/outblocks-plugin-go"

func Any(m map[string]interface{}, key, msg string) (res plugin.Response, val interface{}) {
	if v, ok := m[key]; ok {
		return nil, v
	}

	return &plugin.ValidationErrorResponse{
		Path:  key,
		Error: msg,
	}, nil
}

func String(m map[string]interface{}, key, msg string) (res plugin.Response, val string) {
	res, v := Any(m, key, msg)
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

func OptionalString(def string, m map[string]interface{}, key, msg string) (res plugin.Response, val string) {
	if _, ok := m[key]; !ok {
		return nil, def
	}

	return String(m, key, msg)
}
