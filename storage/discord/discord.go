package discord

import (
	"io"

	"github.com/Karitham/otc/storage"
	"github.com/Karitham/webhook"
)

// check impl
var _ = storage.Storer((*Client)(nil))

// Client holds all that's required to upload to dropbox
type Client struct {
	Hook     *webhook.Hook
	Filename string
}

// Store implements storage.Storer
func (c *Client) Store(file io.Reader) error {
	c.Hook.Webhook.Files = []webhook.Attachment{
		{
			Body:     file,
			Filename: c.Filename,
		},
	}
	return c.Hook.Run()
}

// New returns a new client
func New(url string) *Client {
	return &Client{
		Hook: webhook.New(url),
	}
}
