package types

type ApplyAction struct {
	Source     string   `json:"source"`
	Type       PlanType `json:"type"`
	Namespace  string   `json:"namespace"`
	ObjectID   string   `json:"object_id"`
	ObjectType string   `json:"object_type"`
	ObjectName string   `json:"object_name"`
	Progress   int      `json:"progress"`
	Total      int      `json:"total"`
}
