package dropbox

import (
	"io"
	"net/http"
	"time"

	"github.com/Karitham/otc/runner"
	"github.com/tj/go-dropbox"
	"github.com/tj/go-dropy"
	"github.com/urfave/cli/v2"
)

// args holds all that's required to upload to dropbox
type args struct {
	DP *dropy.Client

	// flags
	filename string
	token    string
	timeout  time.Duration
}

// Command returns dropbox as a storer
func Command() *cli.Command {
	a := &args{}

	return &cli.Command{
		Name:  "dropbox",
		Usage: "store in dropbox",
		Before: func(c *cli.Context) error {
			a.DP = dropy.New(
				dropbox.New(
					&dropbox.Config{
						HTTPClient: &http.Client{
							Timeout: a.timeout,
						}, AccessToken: a.token,
					},
				),
			)

			runner.FromCtx(c.Context).Storer(a)
			return nil
		},
		Action: func(*cli.Context) error { return nil },
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "filename",
				Aliases:     []string{"f"},
				EnvVars:     []string{"DROPBOX_FILENAME"},
				DefaultText: "otc_" + time.Now().Format(time.Kitchen),
				Destination: &a.filename,
			},
			&cli.StringFlag{
				Name:        "token",
				Aliases:     []string{"t"},
				EnvVars:     []string{"DROPBOX_TOKEN"},
				Destination: &a.token,
			},
			&cli.DurationFlag{
				Name:        "timeout",
				Value:       5 * time.Minute,
				Destination: &a.timeout,
			},
		},
	}
}

// Store implements storage.Storer
func (a *args) Store(file io.Reader) error {
	err := a.DP.Upload("/"+a.filename, file)
	if err != nil {
		return err
	}
	return nil
}
