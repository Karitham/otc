package command

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/Karitham/otc/source"
	"github.com/rs/zerolog/log"
)

// check impl
var _ = source.Getter((*Config)(nil))

// Config is the configuration of the getter,
// for our particular case, it's this simple struct but it could be
// larger depending on needs
type Config struct {
	command string
	args    []string
}

// Get implements source.Getter to retrieve the backup from inside the container direclty
func (c *Config) Get(ctx context.Context) (io.ReadCloser, error) {
	log.Trace().Str("command", c.command).Msg("running command")

	f, err := os.CreateTemp(os.TempDir(), "opc_*")
	if err != nil {
		return nil, fmt.Errorf("command: creating temp file %w", err)
	}

	cmd := exec.CommandContext(ctx, c.command, c.args...)
	cmd.Stderr = f
	cmd.Stdout = f

	if err = cmd.Run(); err != nil {
		return nil, fmt.Errorf("command: running %s: %w", c.command, err)
	}
	log.Trace().Msg("command ran successfully")

	_, err = f.Seek(0, 0)
	return f, err
}

// New returns a new config object
func New(command string, args ...string) *Config {
	return &Config{
		command: command,
		args:    args,
	}
}
