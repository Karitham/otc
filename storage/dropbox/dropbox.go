package dropbox

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/Karitham/otc/runner"
	"github.com/rs/zerolog/log"
	"github.com/tj/go-dropbox"
	"github.com/tj/go-dropy"
	"github.com/urfave/cli/v2"
)

// Args holds all that's required to upload to dropbox
type Args struct {
	DP       *dropy.Client
	filename string
	token    string
}

func Command() *cli.Command {
	args := &Args{}

	return &cli.Command{
		Name:  "dropbox",
		Usage: "store in dropbox",
		Before: func(c *cli.Context) error {
			args.DP = dropy.New(
				dropbox.New(
					&dropbox.Config{
						HTTPClient: &http.Client{
							Timeout: 5 * time.Minute,
						}, AccessToken: args.token,
					},
				),
			)

			m, ok := c.Context.Value(runner.K).(*runner.Default)
			if !ok {
				log.Error().Msg("Invalid runner provided")
				c.Done()
			}

			c.Context = context.WithValue(c.Context, runner.K, m.Storer(args))
			return nil
		},
		Action: func(*cli.Context) error { return nil },
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "filename",
				Aliases:     []string{"f"},
				EnvVars:     []string{"DROPBOX_FILENAME"},
				DefaultText: "otc_" + time.Now().Format(time.Kitchen),
				Destination: &args.filename,
			},
			&cli.StringFlag{
				Name:        "token",
				Aliases:     []string{"t"},
				EnvVars:     []string{"DROPBOX_TOKEN"},
				Destination: &args.token,
			},
		},
	}
}

// Store implements storage.Storer
func (a *Args) Store(file io.Reader) error {
	err := a.DP.Upload("/"+a.filename, file)
	if err != nil {
		return err
	}
	return nil
}
