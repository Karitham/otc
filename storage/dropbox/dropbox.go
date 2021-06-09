package dropbox

import (
	"io"
	"net/http"
	"time"

	"github.com/tj/go-dropbox"
	"github.com/tj/go-dropy"
)

type Client struct {
	DP       *dropy.Client
	filename string
}

func (c *Client) Store(file io.Reader) error {
	err := c.DP.Upload("/"+c.filename, file)
	if err != nil {
		return err
	}
	return nil
}

func New(filename string, token string) *Client {
	return &Client{
		filename: filename,
		DP: dropy.New(dropbox.New(&dropbox.Config{
			HTTPClient: &http.Client{
				Timeout: 5 * time.Minute,
			},
			AccessToken: token,
		},
		)),
	}
}
