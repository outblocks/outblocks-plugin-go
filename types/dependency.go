package types

type Dependency struct {
	Type       string                 `json:"type"`
	Properties map[string]interface{} `json:"properties"`
}
