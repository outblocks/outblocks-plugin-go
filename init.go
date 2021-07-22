package plugin

type InitRequest struct {
	Name          string                 `json:"name"`
	DeployPlugins []string               `json:"deploy_plugins"`
	RunPlugins    []string               `json:"run_plugins"`
	Args          map[string]interface{} `json:"args"`
}

func (r *InitRequest) Type() RequestType {
	return RequestTypeInit
}

type InitResponse struct {
	Properties map[string]interface{} `json:"properties"`
}

func (r *InitResponse) Type() ResponseType {
	return ResponseTypeInit
}
