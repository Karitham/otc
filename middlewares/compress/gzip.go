package compress

import (
	"bytes"
	"compress/gzip"
	"context"
	"io"

	"github.com/Karitham/otc/cmd"
	"github.com/Karitham/otc/runner"
	"github.com/Karitham/otc/source"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func Command() cmd.Middleware {
	g := &Gzip{}

	return cmd.Middleware{
		Func: func(bf cli.BeforeFunc) cli.BeforeFunc {
			return func(c *cli.Context) error {
				runner.FromCtx(c.Context).GetterWith(g.Compress(c.Context))
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
func (g *Gzip) Compress(ctx context.Context) func(getter source.Getter) source.Getter {
	return func(getter source.Getter) source.Getter {
		if !g.enabled {
			return getter
		}
		log.Trace().Msg("Entered gz middleware")

		buf := &bytes.Buffer{}

		got, err := getter.Get(ctx)
		if err != nil {
			return getter
		}
		defer got.Close()

		gz, err := gzip.NewWriterLevel(buf, g.level)
		if err != nil {
			return getter
		}
		_, err = io.Copy(gz, got)
		if err != nil {
			return getter
		}
		gz.Close()

		g.r = io.NopCloser(buf)

		log.Trace().Msg("Left gz middleware")

		return g
	}
}

func (g *Gzip) Get(context.Context) (io.ReadCloser, error) {
	return g.r, nil
}
