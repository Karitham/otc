package once

import (
	"context"
	"errors"

	"github.com/Karitham/otc/runner"
	"github.com/Karitham/otc/source"
	"github.com/Karitham/otc/storage"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

// args are the command args
//
// Here there are none, but they follow the same pattern
// as other commands for clarity
type args struct{}

// Command return once as a Runner
func Command() *cli.Command {
	a := &args{}

	return &cli.Command{
		Name:        "once",
		Usage:       "once <getter> <storer>",
		Aliases:     []string{"o"},
		Description: "run once",
		Before: func(c *cli.Context) error {
			runner.FromCtx(c.Context).Runner(a.Run)
			return nil
		},
	}
}

// Run periodically runs the getter and uploads it to the storer
func (a *args) Run(ctx context.Context, g source.Getter, s storage.Storer) error {
	if g == nil {
		return errors.New("no getter provided")
	}
	if s == nil {
		return errors.New("no storer provided")
	}

	if err := getStore(ctx, g, s); err != nil {
		return err
	}

	log.Info().Msg("data retrieved and stored")

	return nil
}

// getStore gets the value then stores it
func getStore(ctx context.Context, g source.Getter, s storage.Storer) error {
	log.Debug().Msg("getting the data")
	r, err := g.Get(ctx)
	if err != nil {
		return err
	}
	defer r.Close()

	log.Debug().Msg("storing the data")
	if err := s.Store(r); err != nil {
		return err
	}
	log.Debug().Msg("done storing")
	return nil
}
