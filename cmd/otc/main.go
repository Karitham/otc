package main

import (
	"context"
	"os"

	"github.com/Karitham/otc/cmd"
	"github.com/Karitham/otc/middlewares/compress"
	"github.com/Karitham/otc/runner"
	"github.com/Karitham/otc/runner/once"
	"github.com/Karitham/otc/runner/periodic"
	command "github.com/Karitham/otc/source/cmd"
	"github.com/Karitham/otc/source/pgdocker"
	"github.com/Karitham/otc/source/stdin"
	"github.com/Karitham/otc/storage/discord"
	"github.com/Karitham/otc/storage/dropbox"
	"github.com/Karitham/otc/storage/stdout"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func main() {
	otc := cmd.OTC{}
	otc.RegisterGetter(
		pgdocker.Command(),
		command.Command(),
		stdin.Command(),
	)
	otc.RegisterStorer(
		discord.Command(),
		dropbox.Command(),
		stdout.Command(),
	)
	otc.RegisterRunner(
		periodic.Command(),
		once.Command(),
	)
	otc.RegisterMiddleware(
		compress.Command(),
	)

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	log.Logger = log.Level(zerolog.InfoLevel)
	var verbose bool

	r := &runner.Default{}

	app := &cli.App{
		Name:    "otc",
		Usage:   "over to cloud - Run once or periodically, pick and store!",
		Version: "0.1",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "verbose",
				Aliases:     []string{"V"},
				Destination: &verbose,
				Value:       false,
			},
		},
		Before: func(c *cli.Context) error {
			if verbose {
				log.Logger = log.Level(zerolog.TraceLevel)
			}
			c.Context = context.WithValue(c.Context, runner.K, r)

			return nil
		},
		Commands: otc.Commands(),
		After: func(c *cli.Context) error {
			m, ok := c.Context.Value(runner.K).(*runner.Default)
			if !ok {
				log.Error().Msg("Invalid runner provided")
				return nil
			}

			return m.Run(c.Context)
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("there was an error running your command")
	}
}
