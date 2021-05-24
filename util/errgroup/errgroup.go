package errgroup

import (
	"context"

	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

type Group struct {
	*errgroup.Group

	ctx context.Context
	sem *semaphore.Weighted
}

func WithConcurrency(ctx context.Context, concurrency int) (*Group, context.Context) {
	sem := semaphore.NewWeighted(int64(concurrency))
	gr, ctx := errgroup.WithContext(ctx)

	return &Group{Group: gr, sem: sem, ctx: ctx}, ctx
}

func (g *Group) Go(f func() error) {
	g.Group.Go(func() error {
		err := g.sem.Acquire(g.ctx, 1)
		if err != nil {
			return err
		}
		defer g.sem.Release(1)

		return f()
	})
}
