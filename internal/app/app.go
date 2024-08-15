package app

import (
	"github.com/lRhythm/shortener/internal/config"
	"github.com/lRhythm/shortener/internal/service"
	"github.com/lRhythm/shortener/internal/storage"
	"github.com/lRhythm/shortener/internal/transport/httptransport"
	"log"
)

func Start() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}
	s, err := httptransport.New(
		cfg,
		service.New(
			service.WithStorage(storage.NewInMemory()),
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(s.Listen())
}
