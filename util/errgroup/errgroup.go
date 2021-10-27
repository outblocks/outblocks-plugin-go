package errgroup

import (
	"context"

	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

type Runner interface {
	Go(f func() error)
	Wait() error
}

type Group struct {
	*errgroup.Group

	ctx context.Context
	sem *semaphore.Weighted
}

func WithContext(ctx context.Context) (Runner, context.Context) {
	return errgroup.WithContext(ctx)
}

func WithConcurrency(ctx context.Context, concurrency int) (Runner, context.Context) {
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
