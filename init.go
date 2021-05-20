package plugin

type InitRequest struct {
	Properties map[string]interface{} `json:"properties"`
}

func (r *InitRequest) Type() RequestType {
	return RequestTypeInit
}
