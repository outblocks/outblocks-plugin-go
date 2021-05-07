package types

type App struct {
	Name       string                 `json:"name"`
	Type       string                 `json:"type"`
	URL        string                 `json:"url"`
	Deploy     string                 `json:"deploy"`
	Needs      map[string]AppNeed     `json:"needs"`
	Properties map[string]interface{} `json:"properties"`
}

type AppNeed struct {
	Properties map[string]interface{} `json:"properties"`
}
