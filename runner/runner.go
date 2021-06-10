package runner

import (
	"context"

	"github.com/Karitham/otc/source"
	"github.com/Karitham/otc/storage"
)

type CtxKey struct{}

var K CtxKey = CtxKey{}

// A Runner runs the final function
type Runner interface {
	Run(context.Context, source.Getter, storage.Storer) error
}

// Default is a default runner.
// It's used to wrapper runners inside middlewares
type Default struct {
	getter      source.Getter
	storer      storage.Storer
	runner      Runner
	middlewares []func(Runner) Runner
}

// Runner sets the default runner
func (d *Default) Runner(r Runner) *Default {
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
func (d *Default) With(m ...func(Runner) Runner) *Default {
	d.middlewares = append(d.middlewares, m...)
	return d
}

// Run applies the middlewares then runs the command
func (d *Default) Run(ctx context.Context) error {
	if len(d.middlewares) == 0 {
		return d.runner.Run(ctx, d.getter, d.storer)
	}

	m := d.middlewares[len(d.middlewares)-1](d.runner)
	for i := len(d.middlewares) - 2; i >= 0; i-- {
		m = d.middlewares[i](m)
	}

	return m.Run(ctx, d.getter, d.storer)
}

// NoOp is a NoOp runner
type NoOp struct{}

func (NoOp) Run(context.Context, source.Getter, storage.Storer) error { return nil }
