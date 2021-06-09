package pgdocker

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/rs/zerolog/log"
)

type Config struct {
	ContainerName string
	DbName        string
	DbUser        string
}

// Get implements source.Getter to retrieve the backup from inside the container direclty
//
// **Don't forget to close the resulting file**
func (c *Config) Get(ctx context.Context) (io.ReadCloser, error) {
	dumpCommand := fmt.Sprintf("pg_dump -U %s -d %s", c.DbUser, c.DbName)

	log.Trace().Str("docker_command", dumpCommand).Msg("running command")

	f, err := os.CreateTemp(os.TempDir(), c.ContainerName+"-*")
	if err != nil {
		return nil, fmt.Errorf("pgdocker: creating temp file %w", err)
	}

	cmd := exec.CommandContext(ctx, "docker", "exec", c.ContainerName, "/bin/bash", "-c", dumpCommand)
	cmd.Stderr = f
	cmd.Stdout = f

	if err = cmd.Run(); err != nil {
		return nil, fmt.Errorf("pgdocker: running command %w", err)
	}

	_, err = f.Seek(0, 0)

	return f, err
}

func New(ContainerName, DbName, DbUser string) *Config {
	return &Config{
		ContainerName: ContainerName,
		DbName:        DbName,
		DbUser:        DbUser,
	}
}
