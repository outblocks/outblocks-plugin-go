package plugin

import "github.com/outblocks/outblocks-plugin-go/types"

type ApplyRequest struct {
	DeployBaseRequest
}

func (r *ApplyRequest) Type() RequestType {
	return RequestTypeApply
}

type ApplyResponse struct {
	Actions []*types.ApplyAction `json:"actions"`
}

func (r *ApplyResponse) Type() ResponseType {
	return ResponseTypeApply
}

type ApplyDoneResponse struct {
	PluginState      *types.PluginState                `json:"plugin_state"`
	AppStates        map[string]*types.AppState        `json:"app_states"`
	DependencyStates map[string]*types.DependencyState `json:"dep_states"`
}

func (r *ApplyDoneResponse) Type() ResponseType {
	return ResponseTypeApplyDone
}
