package communication

type CommandRequest struct {
	Command string   `json:"cmd"`
	Args    []string `json:"args"`
}

func (r *CommandRequest) Type() RequestType {
	return RequestTypeCommand
}
