package pgdocker

import (
	"context"
	"io"
	"os"
	"testing"
)

func TestGet(t *testing.T) {
	cname := os.Getenv("CONTAINER_NAME")
	if cname == "" {
		t.Fatal("no name provided")
	}
	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		t.Fatal("no name provided")
	}
	dbuser := os.Getenv("DB_USER")
	if dbuser == "" {
		dbuser = "postgres"
	}
	c := New(cname, dbname, dbuser)

	r, err := c.Get(context.Background())
	if err != nil {
		t.Fatalf("error getting the data: %s", err)
	}
	defer r.Close()

	b, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("error reading the data: %s", err)
	}
	if string(b) == "" {
		t.Fatal("no data retrieved")
	}
}
