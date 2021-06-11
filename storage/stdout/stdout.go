package stdout

import (
	"io"
	"os"

	"github.com/Karitham/otc/runner"
	"github.com/urfave/cli/v2"
)

type stdout struct{}

func Command() *cli.Command {
	return &cli.Command{
		Name:  "stdout",
		Usage: "run a command and use stdout as the output",
		Before: func(c *cli.Context) error {
			runner.FromCtx(c.Context).Storer(stdout{})
			return nil
		},
		Action: func(*cli.Context) error { return nil },
	}
}

// Get implements source.Getter to retrieve the backup from inside the container direclty
func (stdout) Store(s io.Reader) error {
	_, err := io.Copy(os.Stdout, s)
	return err
}
