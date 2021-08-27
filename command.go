package plugin

type CommandRequest struct {
	Command string                 `json:"cmd"`
	Args    map[string]interface{} `json:"args"`
}

func (r *CommandRequest) Type() RequestType {
	return RequestTypeCommand
}
