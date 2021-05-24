package types

type ApplyAction struct {
	Object      string `json:"object"`
	Description string `json:"description"`
	Progress    int    `json:"progress"`
	Total       int    `json:"total"`
}

func (a *ApplyAction) ProgressIncBy(cnt int) *ApplyAction {
	b := *a
	b.Progress += cnt

	return &b
}

func (a *ApplyAction) ProgressInc() *ApplyAction {
	return a.ProgressIncBy(1)
}
