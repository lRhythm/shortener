package app

import (
	"github.com/lRhythm/shortener/internal/config"
	"github.com/lRhythm/shortener/internal/logs"
	"github.com/lRhythm/shortener/internal/service"
	"github.com/lRhythm/shortener/internal/storage"
	"github.com/lRhythm/shortener/internal/transport/rest"
)

func Start() {
	logger := logs.New()

	cfg, err := config.New()
	if err != nil {
		logger.Fatal(err)
	}

	s, err := rest.New(
		logger,
		cfg,
		service.New(
			service.WithStorage(storage.NewInMemory()),
		),
	)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Fatal(s.Listen())
}
