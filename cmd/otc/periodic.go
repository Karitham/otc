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
		if err := GetStore(ctx, g, s); err != nil {
			return err
		}
		log.Info().Msg("finished storing")
	}

	return nil
}

// GetStore gets the value then stores it
func GetStore(ctx context.Context, g source.Getter, s storage.Storer) error {
	log.Debug().Msg("getting the data")
	r, err := g.Get(ctx)
	if err != nil {
		return err
	}
	defer r.Close()

	log.Debug().Msg("storing the data")
	if err := s.Store(r); err != nil {
		return err
	}
	log.Debug().Msg("done storing")
	return nil
}
