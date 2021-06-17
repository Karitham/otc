package runner

import (
	"context"
	"errors"
	"io"
	"os"

	"github.com/Karitham/otc/source"
	"github.com/Karitham/otc/storage"
)

// ctxKey is used to build a unique context key
type ctxKey struct{}

// K is a unique ctx key
var K ctxKey = ctxKey{}

// Aliases to simplify types
type (
	// runnerFunc is a simple alias for a runner function
	runnerFunc = func(context.Context, source.Getter, storage.Storer) error
	// gMiddle is a getter middleware
	gMiddle = func(source.Getter) source.Getter
	// sMiddle is a storer middleware
	sMiddle = func(storage.Storer) storage.Storer
)

// Runner is a default runner.
// It's used to wrapper runners inside middlewares
type Runner struct {
	getter       source.Getter
	storer       storage.Storer
	runner       runnerFunc
	gMiddlewares []gMiddle
	sMiddlewares []sMiddle
}

// Runner sets the default runner
func (d *Runner) Runner(r runnerFunc) *Runner {
	d.runner = r
	return d
}

// Storer sets the default storer
func (d *Runner) Storer(s storage.Storer) *Runner {
	d.storer = s
	return d
}

// Getter sets the default getter
func (d *Runner) Getter(g source.Getter) *Runner {
	d.getter = g
	return d
}

// GetterWith adds a middleware at the end of the getter chain
func (d *Runner) GetterWith(m ...gMiddle) *Runner {
	d.gMiddlewares = append(d.gMiddlewares, m...)
	return d
}

// StorerWith adds a middleware at the end of the storer chain
func (d *Runner) StorerWith(m ...sMiddle) *Runner {
	d.sMiddlewares = append(d.sMiddlewares, m...)
	return d
}

// Middlewares returns the middlewares
func (d *Runner) Middlewares() ([]gMiddle, []sMiddle) {
	return d.gMiddlewares, d.sMiddlewares
}

// FromCtx returns a DefaultRunner from a context
func FromCtx(ctx context.Context) *Runner {
	d, ok := ctx.Value(K).(*Runner)
	if !ok {
		return nil
	}
	return d
}

// Run applies the middlewares then runs the command
func (d *Runner) Run(ctx context.Context) error {
	if d.runner == nil {
		return nil
	}

	return d.runner(ctx, d, d)
}

// Get implements the source.Getter interface to be able to pass and apply middlewares at each call
func (d *Runner) Get(ctx context.Context) (io.ReadCloser, error) {
	if d.getter == nil {
		return nil, errors.New("no source provided")
	}
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
func (d *Runner) Store(r io.Reader) error {
	if d.storer == nil {
		_, err := io.Copy(os.Stdout, r)
		return err
	}
	if len(d.sMiddlewares) == 0 {
		return d.storer.Store(r)
	}

	s := d.sMiddlewares[len(d.sMiddlewares)-1](d.storer)
	for i := len(d.sMiddlewares) - 2; i >= 0; i-- {
		s = d.sMiddlewares[i](s)
	}

	return s.Store(r)
}
