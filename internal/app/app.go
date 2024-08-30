package app

import (
	"github.com/lRhythm/shortener/internal/config"
	"github.com/lRhythm/shortener/internal/logs"
	"github.com/lRhythm/shortener/internal/service"
	"github.com/lRhythm/shortener/internal/storage"
	"github.com/lRhythm/shortener/internal/transport/rest"
	"os"
	"os/signal"
	"syscall"
)

func Start() {
	logger := logs.New()

	cfg, err := config.New()
	if err != nil {
		logger.Fatal(err)
	}

	db, err := storage.NewMemory(cfg.File())
	if err != nil {
		logger.Fatal(err)
	}
	defer func(stor *storage.Memory) {
		_ = stor.Close()
	}(db)

	s, err := rest.New(
		logger,
		cfg,
		service.New(
			service.WithStorage(db),
		),
	)
	if err != nil {
		logger.Fatal(err)
	}

	sCh := make(chan os.Signal, 1)
	signal.Notify(sCh, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	go func() {
		_ = <-sCh
		logger.Info("server shutting down")
		_ = s.Shutdown()
	}()

	logger.Info("server started")
	if err = s.Listen(); err != nil {
		logger.Fatal(err)
	}

	logger.Info("server shut down")
}
