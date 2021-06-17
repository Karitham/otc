package periodic

import (
	"context"
	"time"

	"github.com/Karitham/otc/runner"
	"github.com/Karitham/otc/source"
	"github.com/Karitham/otc/storage"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

type args struct {
	tick time.Duration
}

// Command return periodic as a command runner
func Command() *cli.Command {
	a := &args{}

	return &cli.Command{
		Name:        "cron",
		Aliases:     []string{"c"},
		Description: "run periodically",
		Usage:       "cron <getter> <storer>",
		Before: func(c *cli.Context) error {
			runner.FromCtx(c.Context).Runner(a.Run)
			return nil
		},
		Flags: []cli.Flag{
			&cli.DurationFlag{
				Name:        "schedule",
				Aliases:     []string{"s"},
				EnvVars:     []string{"SCHEDULE_LOOP"},
				Destination: &a.tick,
				Value:       time.Minute,
			},
		},
	}
}

// Run periodically runs the getter and uploads it to the storer
func (a *args) Run(ctx context.Context, g source.Getter, s storage.Storer) error {
	ticker := time.NewTicker(a.tick)
	go func() {
		<-ctx.Done()
		ticker.Stop()
	}()

	for range ticker.C {
		if err := getStore(ctx, g, s); err != nil {
			return err
		}
		log.Info().Msg("data retrieved and stored")
	}

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
	return nil
}
