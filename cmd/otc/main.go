package main

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/Karitham/otc/source/pgdocker"
	"github.com/Karitham/otc/storage/dropbox"
	"github.com/urfave/cli/v2"
)

type arguments struct {
	filename      string
	token         string
	containerName string
	dbName        string
	dbUser        string
	cronLoop      time.Duration
	verbose       bool
}

func main() {
	var args arguments

	app := &cli.App{
		Name:  "otc",
		Usage: "Out To Cloud, run command and upload their results to cloud",
		Before: func(c *cli.Context) error {
			log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

			if args.verbose {
				log.Logger = log.Level(zerolog.TraceLevel)
			}

			return nil
		},
		Commands: []*cli.Command{
			{
				Name:        "cron",
				Aliases:     []string{"c"},
				Description: "run periodically",
				Action: func(c *cli.Context) error {
					s := dropbox.New(args.filename, args.token)
					g := pgdocker.New(args.containerName, args.dbName, args.dbUser)

					return Periodically(c.Context, args.cronLoop, s, g)
				},
				Flags: []cli.Flag{
					&cli.DurationFlag{
						Name:        "schedule",
						Aliases:     []string{"s"},
						EnvVars:     []string{"SCHEDULE_LOOP"},
						Destination: &args.cronLoop,
						Value:       time.Minute,
					},
					&cli.StringFlag{
						Name:        "filename",
						Aliases:     []string{"f"},
						EnvVars:     []string{"FILE"},
						DefaultText: "otc_" + time.Now().Format(time.Stamp),
						Destination: &args.filename,
					},
					&cli.StringFlag{
						Name:        "token",
						Aliases:     []string{"t"},
						EnvVars:     []string{"DB_TOKEN"},
						Destination: &args.token,
					},
					&cli.StringFlag{
						Name:        "user",
						Aliases:     []string{"u"},
						EnvVars:     []string{"DB_USER"},
						Value:       "postgres",
						Destination: &args.dbUser,
					},
					&cli.StringFlag{
						Name:        "database",
						Aliases:     []string{"d"},
						EnvVars:     []string{"DB_NAME"},
						Destination: &args.dbName,
					},
					&cli.StringFlag{
						Name:        "container",
						Aliases:     []string{"c"},
						Destination: &args.containerName,
					},
				},
			},
		},
		Flags: []cli.Flag{&cli.BoolFlag{
			Name:        "verbose",
			Aliases:     []string{"v"},
			Destination: &args.verbose,
		}},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal().Err(err).Msg("there was an error running your command")
	}
}
