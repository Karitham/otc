package cmd

import (
	"context"
	"io"
	"testing"
)

func TestGet(t *testing.T) {
	c := Args{
		args:    []string{"hello world"},
		command: "echo",
	}

	r, err := c.Get(context.Background())
	if err != nil {
		t.Fatalf("error getting the data: %s", err)
	}
	defer r.Close()

	b, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("error reading the data: %s", err)
	}
	if string(b) == "hello world" {
		t.Fatal("no data retrieved")
	}
}
