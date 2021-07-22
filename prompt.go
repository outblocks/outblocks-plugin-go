package plugin

type PromptConfirmation struct {
	Message string `json:"message"`
	Default bool   `json:"default"`
}

func (r *PromptConfirmation) Type() ResponseType {
	return ResponseTypePromptConfirmation
}

type PromptSelect struct {
	Message string   `json:"message"`
	Options []string `json:"options"`
	Default string   `json:"default"`
}

func (r *PromptSelect) Type() ResponseType {
	return ResponseTypePromptSelect
}

type PromptInput struct {
	Message string `json:"message"`
	Default string `json:"default"`
}

func (r *PromptInput) Type() ResponseType {
	return ResponseTypePromptInput
}

type PromptConfirmationAnswer struct {
	Confirmed bool `json:"confirmed"`
}

func (r *PromptConfirmationAnswer) Type() RequestType {
	return RequestTypePromptConfirmation
}

type PromptInputAnswer struct {
	Answer string `json:"answer"`
}

func (r *PromptInputAnswer) Type() RequestType {
	return RequestTypePromptAnswer
}
