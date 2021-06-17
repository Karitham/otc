package discord

import (
	"io"
	"time"

	"github.com/Karitham/otc/runner"
	"github.com/Karitham/webhook"
	"github.com/urfave/cli/v2"
)

// Client holds all that's required to upload to dropbox
type args struct {
	Hook *webhook.Hook

	// flags
	filename string
	url      string
}

// Command returns discord as a storer
func Command() *cli.Command {
	a := &args{}

	return &cli.Command{
		Name:  "discord",
		Usage: "store in a discord channel via webhook",
		Before: func(c *cli.Context) error {
			runner.FromCtx(c.Context).Storer(a)

			a.Hook = webhook.New(a.url)
			return nil
		},
		Action: func(*cli.Context) error { return nil },
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "url",
				Destination: &a.url,
				Aliases:     []string{"u"},
				EnvVars:     []string{"WEBHOOK_URL"},
			},
			&cli.StringFlag{
				Name:        "file",
				Aliases:     []string{"f"},
				EnvVars:     []string{"FILENAME"},
				Destination: &a.filename,
				Value:       "otc_" + time.Now().Format(time.Kitchen),
			},
		},
	}
}

// Store implements storage.Storer
func (c *args) Store(file io.Reader) error {
	c.Hook.With(&webhook.Webhook{
		Files: []webhook.Attachment{
			{
				Body:     file,
				Filename: c.filename,
			},
		},
	})
	return c.Hook.Run()
}
