package plugin_go

import "github.com/outblocks/outblocks-plugin-go/types"

type ApplyRequest struct {
	Apps         []*types.AppPlan        `json:"apps,omitempty"`
	Dependencies []*types.DependencyPlan `json:"dependencies,omitempty"`
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
