package cmd

import (
	"github.com/urfave/cli/v2"
)

// OTC, where to register your new commands
type OTC struct {
	storers []*cli.Command
	getters []*cli.Command
	runners []*cli.Command

	middlewares []Middleware
}

// Middleware is a command middleware
type Middleware struct {
	Func    func(cli.BeforeFunc) cli.BeforeFunc
	Command *cli.BoolFlag
	Args    []cli.Flag
}

// RegisterStorer registers a new storer
func (o *OTC) RegisterStorer(s ...*cli.Command) {
	o.storers = append(o.storers, s...)
}

// RegisterGetter registers a new getter
func (o *OTC) RegisterGetter(g ...*cli.Command) {
	o.getters = append(o.getters, g...)
}

// RegisterRunner registers a new ruunner
func (o *OTC) RegisterRunner(r ...*cli.Command) {
	o.runners = append(o.runners, r...)
}

// RegisterMiddleware registers a middleware
func (o *OTC) RegisterMiddleware(m ...Middleware) {
	o.middlewares = append(o.middlewares, m...)
}

// NestSubCommands give us the ability to use a command after the other without caring about order
func makeCommands(runners []*cli.Command, storers []*cli.Command, getters []*cli.Command, middlewares []Middleware) []*cli.Command {
	for i := range storers {
		storers[i].Subcommands = append(storers[i].Subcommands, getters...)
	}
	for i := range getters {
		getters[i].Subcommands = append(getters[i].Subcommands, storers...)
	}

	gs := append(storers, getters...)

	for i := range runners {
		runners[i].Subcommands = gs
		for j := range middlewares {
			runners[i].Before = middlewares[j].Func(runners[i].Before)
			runners[i].Flags = append(runners[i].Flags, middlewares[j].Command)
			runners[i].Flags = append(runners[i].Flags, middlewares[j].Args...)
		}
	}

	return runners
}

// Commands returns the registered commands
func (o *OTC) Commands() []*cli.Command {
	return makeCommands(o.runners, o.storers, o.getters, o.middlewares)
}
