package dropbox

import (
	"io"
	"net/http"
	"time"

	"github.com/Karitham/otc/storage"
	"github.com/tj/go-dropbox"
	"github.com/tj/go-dropy"
)

// check impl
var _ = storage.Storer((*Client)(nil))

// Client holds all that's required to upload to dropbox
type Client struct {
	DP       *dropy.Client
	filename string
}

// Store implements storage.Storer
func (c *Client) Store(file io.Reader) error {
	err := c.DP.Upload("/"+c.filename, file)
	if err != nil {
		return err
	}
	return nil
}

// New returns a new client
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
