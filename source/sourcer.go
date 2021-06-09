package source

import (
	"context"
	"io"
)

// Getter represents any type that can be acquired
// The goal is to provide a simple interface to get data from.
//
// **Don't forget to close the resulting ReadCloser**
type Getter interface {
	Get(context.Context) (io.ReadCloser, error)
}
