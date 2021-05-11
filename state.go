package plugin_go

type GetStateRequest struct {
	StateType  string                 `json:"type"`
	Properties map[string]interface{} `json:"properties"`
	Lock       bool                   `json:"lock"`
}

func (r *GetStateRequest) Type() RequestType {
	return RequestTypeGetState
}

type GetStateResponse struct {
	DataResponse `json:",inline"`
	LockInfo     []byte `json:"lockinfo"`
}

func (r *GetStateResponse) Type() ResponseType {
	return ResponseTypeGetState
}

type SaveStateRequest struct {
	LockInfo   []byte                 `json:"lockinfo"`
	StateType  string                 `json:"type"`
	Properties map[string]interface{} `json:"properties"`
}

func (r *SaveStateRequest) Type() RequestType {
	return RequestTypeSaveState
}

type SaveStateResponse struct {
	EmptyResponse
}
