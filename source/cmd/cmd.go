package cmd

import (
	"bytes"
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
var _ source.Getter = (*args)(nil)

// args is the configuration of the getter,
type args struct {
	command string
	args    []string
}

// Command returns cmd as a command getter
func Command() *cli.Command {
	a := &args{}
	var s string

	return &cli.Command{
		Name:  "cmd",
		Usage: "run a command and use stdout as the source",
		Before: func(c *cli.Context) error {
			a.args = strings.Fields(s)
			runner.FromCtx(c.Context).Getter(a)
			return nil
		},
		Action: func(*cli.Context) error { return nil },
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "command",
				Aliases:     []string{"c"},
				Destination: &a.command,
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
func (a *args) Get(ctx context.Context) (io.ReadCloser, error) {
	log.Trace().Str("command", a.command).Strs("args", a.args).Msg("running command")

	f, err := os.CreateTemp(os.TempDir(), "opc_*")
	if err != nil {
		return nil, fmt.Errorf("command: creating temp file %w", err)
	}

	stderrBuf := &bytes.Buffer{}

	cmd := exec.CommandContext(ctx, a.command, a.args...)
	cmd.Stdout = f
	cmd.Stderr = stderrBuf

	cmdErr := cmd.Run()

	_, err = f.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	if cmdErr != nil {
		return nil, fmt.Errorf("cmd: running command %w. Output: %s", cmdErr, stderrBuf.String())
	}
	log.Trace().Msg("command ran successfully")

	return f, err
}
