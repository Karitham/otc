package discord

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/Karitham/webhook"
)

func TestClient_Store(t *testing.T) {
	url := os.Getenv("WEBHOOK_URL")
	if url == "" {
		t.Skip("no webhook url provided")
	}

	type fields struct {
		Hook     *webhook.Hook
		filename string
	}
	type ags struct {
		file io.Reader
	}
	tests := []struct {
		args    ags
		fields  fields
		name    string
		wantErr bool
	}{
		{
			args: ags{
				file: bytes.NewBuffer([]byte("hello world")),
			},
			fields: fields{
				Hook:     webhook.New(url),
				filename: "hello_world.txt",
			},
			name:    "hello world test",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &args{
				Hook:     tt.fields.Hook,
				filename: tt.fields.filename,
			}
			if err := c.Store(tt.args.file); (err != nil) != tt.wantErr {
				t.Errorf("Client.Store() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
