package plugin_go

import "github.com/outblocks/outblocks-plugin-go/types"

type ApplyRequest struct {
	Plan       *types.Plan `json:"plan"`
	DeployPlan *types.Plan `json:"deploy,omitempty"`
	DNSPlan    *types.Plan `json:"dns,omitempty"`
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
	Plugin     types.PluginStateMap `json:"plugin"`
	DeployPlan *types.Plan          `json:"deploy,omitempty"`
	DNSPlan    *types.Plan          `json:"dns,omitempty"`
}

func (r *ApplyDoneResponse) Type() ResponseType {
	return ResponseTypeApplyDone
}
