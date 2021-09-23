package plugin

type ProjectInitRequest struct {
	Name          string                 `json:"name"`
	DeployPlugins []string               `json:"deploy_plugins"`
	RunPlugins    []string               `json:"run_plugins"`
	Args          map[string]interface{} `json:"args"`
}

func (r *ProjectInitRequest) Type() RequestType {
	return RequestTypeInit
}

type ProjectInitResponse struct {
	Properties map[string]interface{} `json:"properties"`
}

func (r *ProjectInitResponse) Type() ResponseType {
	return ResponseTypeInit
}
