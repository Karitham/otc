package storage

import (
	"io"
)

// Storer represents anywhere you can store into
// The goal is to implement any form of storage,
// from file systems, to cloud
type Storer interface {
	Store(io.Reader) error
}
