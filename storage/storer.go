package storage

import (
	"io"
)

type Storer interface {
	Store(io.Reader) error
}
