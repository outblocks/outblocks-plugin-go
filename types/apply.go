package types

type TargetType int

const (
	TargetTypeApp = iota + 1
	TargetTypeDependency
)

type ApplyAction struct {
	TargetID    string     `json:"target_id"`
	TargetType  TargetType `json:"target_type"`
	Object      string     `json:"object"`
	Description string     `json:"description"`
	Progress    int        `json:"progress"`
	Total       int        `json:"total"`
}

func (a *ApplyAction) WithProgressIncBy(cnt int) *ApplyAction {
	b := *a
	b.Progress += cnt

	return &b
}

func (a *ApplyAction) WithProgressInc() *ApplyAction {
	return a.WithProgressIncBy(1)
}

func (a *ApplyAction) WithProgress(p int) *ApplyAction {
	b := *a
	b.Progress = p

	return &b
}

func (a *ApplyAction) WithDesc(str string) *ApplyAction {
	b := *a
	b.Description = str

	return &b
}
