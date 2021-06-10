package discord

import (
	"context"
	"io"
	"time"

	"github.com/Karitham/otc/runner"
	"github.com/Karitham/webhook"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

// Client holds all that's required to upload to dropbox
type Args struct {
	Hook     *webhook.Hook
	Filename string
	url      string
}

func Command() *cli.Command {
	args := &Args{}

	return &cli.Command{
		Name:  "discord",
		Usage: "store in a discord channel via webhook",
		Before: func(c *cli.Context) error {
			m, ok := c.Context.Value(runner.K).(*runner.Default)
			if !ok {
				log.Error().Msg("Invalid runner provided")
				c.Done()
			}

			args.Hook = webhook.New(args.url)
			c.Context = context.WithValue(c.Context, runner.K, m.Storer(args))
			return nil
		},
		Action: func(_ *cli.Context) error { return nil },
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "url",
				Destination: &args.url,
				Aliases:     []string{"u"},
				EnvVars:     []string{"WEBHOOK_URL"},
			},
			&cli.StringFlag{
				Name:        "file",
				Aliases:     []string{"f"},
				EnvVars:     []string{"FILENAME"},
				Destination: &args.Filename,
				Value:       "otc_" + time.Now().Format(time.Kitchen),
			},
		},
	}
}

// Store implements storage.Storer
func (c *Args) Store(file io.Reader) error {
	c.Hook.Webhook.Files = []webhook.Attachment{
		{
			Body:     file,
			Filename: c.Filename,
		},
	}
	return c.Hook.Run()
}
