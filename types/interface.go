package types

import (
	"context"
)

type Indexed interface {
	Key() string
}

type ApplyCallbackFunc func(desc string)
type ApplyActionCallbackFunc func(obj, desc string, idx, progress, total int)

// callback as a part of resourcedata

type Planner interface {
	Diff(ctx context.Context, data interface{}) (*PlanAction, error)
	Read(ctx context.Context, data interface{})
	Update(ctx context.Context, data interface{})
	Delete(ctx context.Context, data interface{})
	Create(ctx context.Context, data interface{})
	Apply(ctx context.Context) error
}

// changes:
// each resource has a schema
// everything is a list of resources
