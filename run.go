package plugin

import "github.com/outblocks/outblocks-plugin-go/types"

type RunRequest struct {
	Apps         []*types.App           `json:"apps"`
	Dependencies []*types.Dependency    `json:"dependencies"`
	Args         map[string]interface{} `json:"args"`
}

func (r *RunRequest) Type() RequestType {
	return RequestTypeRun
}

type RunDoneResponse struct {
	AppStates        map[string]*types.LocalAccessInfo `json:"app_states"`
	DependencyStates map[string]*types.LocalAccessInfo `json:"dep_states"`
}

func (r *RunDoneResponse) Type() ResponseType {
	return ResponseTypeRunDone
}
