package communication

type StartRequest struct {
	Properties map[string]interface{} `json:"properties"`
}

func (r *StartRequest) Type() RequestType {
	return RequestTypeStart
}
