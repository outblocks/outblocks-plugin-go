package plugin_go

import "github.com/outblocks/outblocks-plugin-go/types"

type ApplyRequest struct {
	Plan *types.Plan `json:"plan"`
}

func (r *ApplyRequest) Type() RequestType {
	return RequestTypeApply
}

type ApplyResponse struct {
	Actions []*types.ApplyAction `json:"actions,omitempty"`
}

func (r *ApplyResponse) Type() ResponseType {
	return ResponseTypeApply
}

type ApplyDoneResponse struct {
	Plugin types.PluginState  `json:"plugin"`
	Deploy *types.StateDeploy `json:"deploy_state"`
}

func (r *ApplyDoneResponse) Type() ResponseType {
	return ResponseTypeApplyDone
}
