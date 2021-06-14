package runner

import (
	"context"
	"io"

	"github.com/Karitham/otc/source"
	"github.com/Karitham/otc/storage"
)

type CtxKey struct{}

var K CtxKey = CtxKey{}

type RunnerFunc = func(context.Context, source.Getter, storage.Storer) error

type (
	gMiddle = func(source.Getter) source.Getter
	sMiddle = func(storage.Storer) storage.Storer
)

// Default is a default runner.
// It's used to wrapper runners inside middlewares
type Default struct {
	getter       source.Getter
	storer       storage.Storer
	runner       RunnerFunc
	gMiddlewares []gMiddle
	sMiddlewares []sMiddle
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
func (d *Default) GetterWith(m ...gMiddle) *Default {
	d.gMiddlewares = append(d.gMiddlewares, m...)
	return d
}

// With adds a middleware at the end of the chain
func (d *Default) StorerWith(m ...sMiddle) *Default {
	d.sMiddlewares = append(d.sMiddlewares, m...)
	return d
}

// Middlewares returns the middlewares
func (d *Default) Middlewares() ([]gMiddle, []sMiddle) {
	return d.gMiddlewares, d.sMiddlewares
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

	return d.runner(ctx, d, d)
}

// Get implements the source.Getter interface to be able to pass and apply middlewares at each call
func (d *Default) Get(ctx context.Context) (io.ReadCloser, error) {
	if len(d.gMiddlewares) == 0 {
		return d.getter.Get(ctx)
	}

	g := d.gMiddlewares[len(d.gMiddlewares)-1](d.getter)
	for i := len(d.gMiddlewares) - 2; i >= 0; i-- {
		g = d.gMiddlewares[i](g)
	}

	return g.Get(ctx)
}

// Store implements the store.Storer interface to be able to pass and apply middlewares at each call
func (d *Default) Store(r io.Reader) error {
	if len(d.sMiddlewares) == 0 {
		return d.storer.Store(r)
	}

	s := d.sMiddlewares[len(d.sMiddlewares)-1](d.storer)
	for i := len(d.sMiddlewares) - 2; i >= 0; i-- {
		s = d.sMiddlewares[i](s)
	}

	return s.Store(r)
}
