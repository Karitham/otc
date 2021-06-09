package source

import (
	"context"
	"io"
)

type Getter interface {
	Get(context.Context) (io.ReadCloser, error)
}
