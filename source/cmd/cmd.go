package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/Karitham/otc/runner"
	"github.com/Karitham/otc/source"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

// check impl
var _ source.Getter = (*Args)(nil)

// Args is the configuration of the getter,
type Args struct {
	command string
	args    []string
}

func Command() *cli.Command {
	args := &Args{}
	var s string

	return &cli.Command{
		Name:  "cmd",
		Usage: "run a command and use stdout/stderr as the source",
		Before: func(c *cli.Context) error {
			args.args = strings.Fields(s)

			m, ok := c.Context.Value(runner.K).(*runner.Default)
			if !ok {
				log.Error().Msg("Invalid runner provided")
				c.Done()
			}
			c.Context = context.WithValue(c.Context, runner.K, m.Getter(args))
			return nil
		},
		Action: func(*cli.Context) error { return nil },
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "command",
				Aliases:     []string{"c"},
				Destination: &args.command,
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "args",
				Aliases:     []string{"a"},
				Destination: &s,
			},
		},
	}
}

// Get implements source.Getter to retrieve the backup from inside the container direclty
func (a *Args) Get(ctx context.Context) (io.ReadCloser, error) {
	log.Trace().Str("command", a.command).Msg("running command")

	f, err := os.CreateTemp(os.TempDir(), "opc_*")
	if err != nil {
		return nil, fmt.Errorf("command: creating temp file %w", err)
	}

	cmd := exec.CommandContext(ctx, a.command, a.args...)
	cmd.Stderr = f
	cmd.Stdout = f

	if err = cmd.Run(); err != nil {
		return nil, fmt.Errorf("command: running %s: %w", a.command, err)
	}
	log.Trace().Msg("command ran successfully")

	_, err = f.Seek(0, 0)
	return f, err
}
