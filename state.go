package plugin_go

import "github.com/outblocks/outblocks-plugin-go/types"

type GetStateRequest struct {
	StateType  string                 `json:"type"`
	Env        string                 `json:"env"`
	Properties map[string]interface{} `json:"properties"`
	Lock       bool                   `json:"lock"`
}

func (r *GetStateRequest) Type() RequestType {
	return RequestTypeGetState
}

type GetStateResponse struct {
	State  *types.StateData   `json:"state"`
	Source *types.StateSource `json:"source"`
}

func (r *GetStateResponse) Type() ResponseType {
	return ResponseTypeGetState
}

type SaveStateRequest struct {
	State      *types.StateData       `json:"state"`
	StateType  string                 `json:"type"`
	Env        string                 `json:"env"`
	Properties map[string]interface{} `json:"properties"`
}

func (r *SaveStateRequest) Type() RequestType {
	return RequestTypeSaveState
}

type SaveStateResponse struct {
	EmptyResponse
}
