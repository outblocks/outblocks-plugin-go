package types

type ApplyAction struct {
	Object      string `json:"object"`
	Description string `json:"description"`
	IsDone      bool   `json:"is_done"`
}
