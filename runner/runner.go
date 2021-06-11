package runner

import (
	"context"

	"github.com/Karitham/otc/source"
	"github.com/Karitham/otc/storage"
)

type CtxKey struct{}

var K CtxKey = CtxKey{}

// A Wither runs the final function
type Wither interface {
	With(m ...Middleware) Wither
}

type RunnerFunc = func(context.Context, source.Getter, storage.Storer) error

type Middleware = func(RunnerFunc) RunnerFunc

// Default is a default runner.
// It's used to wrapper runners inside middlewares
type Default struct {
	getter      source.Getter
	storer      storage.Storer
	runner      RunnerFunc
	middlewares []Middleware
}

// Runner sets the default runner
func (d *Default) Runner(r RunnerFunc) *Default {
	d.runner = r
	return d
}

// Storer sets the default storer
func (d *Default) Storer(s storage.Storer) *Default {
	d.storer = s
	return d
}

// Getter sets the default getter
func (d *Default) Getter(g source.Getter) *Default {
	d.getter = g
	return d
}

// With adds a middleware at the end of the chain
func (d *Default) With(m ...Middleware) *Default {
	d.middlewares = append(d.middlewares, m...)
	return d
}

func FromCtx(ctx context.Context) *Default {
	d, ok := ctx.Value(K).(*Default)
	if !ok {
		return nil
	}
	return d
}

// Run applies the middlewares then runs the command
func (d *Default) Run(ctx context.Context) error {
	if d.runner == nil {
		return nil
	}
	if len(d.middlewares) == 0 {
		return d.runner(ctx, d.getter, d.storer)
	}

	m := d.middlewares[len(d.middlewares)-1](d.runner)
	for i := len(d.middlewares) - 2; i >= 0; i-- {
		m = d.middlewares[i](m)
	}

	return m(ctx, d.getter, d.storer)
}
