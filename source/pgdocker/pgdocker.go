package pgdocker

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"

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
	log.Trace().Str("docker_command", dumpCommand).Msg("running command")

	f, err := os.CreateTemp(os.TempDir(), a.ContainerName+"-*")
	if err != nil {
		return nil, fmt.Errorf("pgdocker: creating temp file %w", err)
	}

	cmd := exec.CommandContext(ctx, "docker", "exec", a.ContainerName, "/bin/bash", "-c", dumpCommand)
	cmd.Stderr = f
	cmd.Stdout = f

	if err = cmd.Run(); err != nil {
		return nil, fmt.Errorf("pgdocker: running command %w", err)
	}
	log.Trace().Msg("command ran successfully")

	_, err = f.Seek(0, 0)
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
