package types

type ApplyAction struct {
	Type       PlanType `json:"type"`
	Namespace  string   `json:"namespace"`
	ObjectID   string   `json:"object_id"`
	ObjectType string   `json:"object_type"`
	ObjectName string   `json:"object_name"`
	Progress   int      `json:"progress"`
	Total      int      `json:"total"`
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
