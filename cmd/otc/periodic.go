package main

import (
	"context"
	"time"

	"github.com/Karitham/otc/source"
	"github.com/Karitham/otc/storage"
	"github.com/rs/zerolog/log"
)

func Periodically(ctx context.Context, t time.Duration, s storage.Storer, g source.Getter) error {
	ticker := time.NewTicker(t)
	go func() {
		<-ctx.Done()
		ticker.Stop()
	}()

	for range ticker.C {
		log.Info().Msg("getting the data")
		r, err := g.Get(ctx)
		if err != nil {
			return err
		}

		log.Info().Msg("storing the data")
		if err := s.Store(r); err != nil {
			return err
		}

		log.Info().Msg("done storing")
	}
	return nil
}
