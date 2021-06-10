package cmd

import (
	"github.com/urfave/cli/v2"
)

// OTC, where to register your new commands
type OTC struct {
	Storers []*cli.Command
	Getters []*cli.Command
	Runners []*cli.Command
}

// RegisterStorer registers a new storer
func (o *OTC) RegisterStorer(s ...*cli.Command) {
	o.Storers = append(o.Storers, s...)
}

// RegisterGetter registers a new getter
func (o *OTC) RegisterGetter(s ...*cli.Command) {
	o.Getters = append(o.Getters, s...)
}

// RegisterRunner registers a new ruunner
func (o *OTC) RegisterRunner(s ...*cli.Command) {
	o.Runners = append(o.Runners, s...)
}

// NestSubCommands give us the ability to use a command after the other without caring about order
func makeCommands(runner []*cli.Command, storer []*cli.Command, getter []*cli.Command) []*cli.Command {
	for i := range storer {
		storer[i].Subcommands = append(storer[i].Subcommands, getter...)
	}
	for i := range getter {
		getter[i].Subcommands = append(getter[i].Subcommands, storer...)
	}
	gs := append(storer, getter...)

	for i := range runner {
		runner[i].Subcommands = gs
	}

	return runner
}

// Commands returns the registered commands
func (o *OTC) Commands() []*cli.Command {
	return makeCommands(o.Runners, o.Storers, o.Getters)
}
