package pgdocker

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
// for our particular case, it's this simple struct but it could be
// larger depending on needs
type Args struct {
	ContainerName string
	DbName        string
	DbUser        string
}

func Command() *cli.Command {
	args := &Args{}

	return &cli.Command{
		Name:  "pgdocker",
		Usage: "run a pg_dump in a docker container",
		Before: func(c *cli.Context) error {
			runner.FromCtx(c.Context).Getter(args)
			return nil
		},
		Action: func(*cli.Context) error { return nil },
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "user",
				Aliases:     []string{"u"},
				EnvVars:     []string{"DB_USER"},
				Value:       "postgres",
				Destination: &args.DbUser,
			},
			&cli.StringFlag{
				Name:        "database",
				Aliases:     []string{"d"},
				EnvVars:     []string{"DB_NAME"},
				Destination: &args.DbName,
			},
			&cli.StringFlag{
				Name:        "container",
				Aliases:     []string{"c"},
				EnvVars:     []string{"CONTAINER_NAME"},
				Destination: &args.ContainerName,
			},
		},
	}
}

// Get implements source.Getter to retrieve the backup from inside the container direclty
//
// **Don't forget to close the resulting file**
func (a *Args) Get(ctx context.Context) (io.ReadCloser, error) {
	dumpCommand := fmt.Sprintf("pg_dump -U %s -d %s", a.DbUser, a.DbName)
	log.Trace().Str("docker_command", dumpCommand).Str("container_name", a.ContainerName).Msg("running command")

	f, err := os.CreateTemp(os.TempDir(), a.ContainerName+"-*")
	if err != nil {
		return nil, fmt.Errorf("pgdocker: creating temp file %w", err)
	}

	cmd := exec.CommandContext(ctx, "docker", "exec", a.ContainerName, "/bin/bash", "-c", dumpCommand)
	cmd.Stderr = f
	cmd.Stdout = f

	cmdErr := cmd.Run()

	_, err = f.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	if cmdErr != nil {
		buf := make([]byte, 2048)
		_, readErr := f.Read(buf)
		str := strings.Trim(string(buf), "\x00")

		if readErr != nil && readErr != io.EOF {
			return nil, fmt.Errorf("pgdocker: running command %w", cmdErr)
		}

		return nil, fmt.Errorf("pgdocker: running command %w. Output: %s", cmdErr, str)
	}
	log.Trace().Msg("command ran successfully")

	return f, err
}

// New returns a new pgdocker config object
func New(ContainerName, DbName, DbUser string) *Args {
	return &Args{
		ContainerName: ContainerName,
		DbName:        DbName,
		DbUser:        DbUser,
	}
}
