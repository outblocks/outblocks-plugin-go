package types

import (
	"context"
)

type Indexed interface {
	Key() string
}

type ApplyCallbackFunc func(desc string)
type ApplyActionCallbackFunc func(obj, desc string, idx, progress, total int)

type Planner interface {
	Plan(ctx context.Context, key string, dest interface{}, verify bool) (*PlanAction, error)
	Apply(ctx context.Context, ops []*PlanActionOperation, callback ApplyCallbackFunc) error
}
