package compress

import (
	"bytes"
	"compress/gzip"
	"context"
	"io"

	"github.com/Karitham/otc/cmd"
	"github.com/Karitham/otc/runner"
	"github.com/Karitham/otc/source"
	"github.com/Karitham/otc/storage"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func Command() cmd.Middleware {
	g := &Gzip{}

	return cmd.Middleware{
		Func: func(bf cli.BeforeFunc) cli.BeforeFunc {
			return func(c *cli.Context) error {
				runner.FromCtx(c.Context).With(g.Compress)
				return bf(c)
			}
		},
		Command: &cli.BoolFlag{
			Name:        "gz",
			Destination: &g.enabled,
			Usage:       "gz implements compression for gzip outputs",
		},
		Args: []cli.Flag{
			&cli.IntFlag{
				Name:        "level",
				Aliases:     []string{"l"},
				Value:       -1,
				Usage:       "compression level between 0-none and 9-best, default is -1",
				Destination: &g.level,
			},
		},
	}
}

// Gzip enables compression
type Gzip struct {
	r io.ReadCloser

	// flgas
	enabled bool
	level   int
}

// Compress implements a simple flate middleware
func (g *Gzip) Compress(r runner.RunnerFunc) runner.RunnerFunc {
	return func(ctx context.Context, get source.Getter, s storage.Storer) error {
		if !g.enabled {
			return r(ctx, get, s)
		}

		log.Trace().Msg("Entered gz middleware")

		buf := &bytes.Buffer{}

		got, err := get.Get(ctx)
		if err != nil {
			return err
		}
		defer got.Close()

		gz, err := gzip.NewWriterLevel(buf, g.level)
		if err != nil {
			return err
		}
		io.Copy(gz, got)
		gz.Close()

		g.r = io.NopCloser(buf)

		log.Trace().Msg("Left gz middleware")
		return r(ctx, g, s)
	}
}

func (g *Gzip) Get(context.Context) (io.ReadCloser, error) {
	return g.r, nil
}
