package plugin_go

type PromptAnswerRequest struct {
	// TODO: prompt support
}

func (r *PromptAnswerRequest) Type() RequestType {
	return RequestTypePromptAnswer
}
