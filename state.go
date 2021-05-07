package communication

type GetStateRequest struct {
	StateType  string                 `json:"type"`
	Properties map[string]interface{} `json:"properties"`
}

func (r *GetStateRequest) Type() RequestType {
	return RequestTypeGetState
}

type GetStateResponse struct {
	DataResponse
}

type SaveStateRequest struct {
	StateType  string                 `json:"type"`
	Properties map[string]interface{} `json:"properties"`
}

func (r *SaveStateRequest) Type() RequestType {
	return RequestTypeSaveState
}

type SaveStateResponse struct {
	EmptyResponse
}
