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
		Filename string
	}
	type args struct {
		file io.Reader
	}
	tests := []struct {
		args    args
		fields  fields
		name    string
		wantErr bool
	}{
		{
			args: args{
				file: bytes.NewBuffer([]byte("hello world")),
			},
			fields: fields{
				Hook:     webhook.New(url),
				Filename: "hello_world.txt",
			},
			name:    "hello world test",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Args{
				Hook:     tt.fields.Hook,
				Filename: tt.fields.Filename,
			}
			if err := c.Store(tt.args.file); (err != nil) != tt.wantErr {
				t.Errorf("Client.Store() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
