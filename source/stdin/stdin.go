package stdin

import (
	"context"
	"io"
	"os"

	"github.com/Karitham/otc/runner"
	"github.com/urfave/cli/v2"
)

type stdin struct{}

func Command() *cli.Command {
	return &cli.Command{
		Name:  "stdin",
		Usage: "run a command and use stdin as the source",
		Before: func(c *cli.Context) error {
			runner.FromCtx(c.Context).Getter(stdin{})
			return nil
		},
		Action: func(*cli.Context) error { return nil },
	}
}

// Get implements source.Getter to retrieve the backup from inside the container direclty
func (stdin) Get(ctx context.Context) (io.ReadCloser, error) {
	return io.NopCloser(os.Stdin), nil
}
