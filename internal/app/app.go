package app

import (
	"os"
	"os/signal"
	"syscall"

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

	var store service.RepositoryInterface
	dns, ok := cfg.DSN()
	if ok {
		store, err = storage.NewDB(dns)
		if err != nil {
			logger.Fatal(err)
		}
		defer store.Close()
	} else {
		store, err = storage.NewMemory(cfg.File())
		if err != nil {
			logger.Fatal(err)
		}
		defer store.Close()
	}

	s, err := rest.New(
		logger,
		cfg,
		service.New(
			service.WithStorage(store),
		),
	)
	if err != nil {
		logger.Fatal(err)
	}

	sCh := make(chan os.Signal, 1)
	signal.Notify(sCh, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	go func() {
		<-sCh
		logger.Info("server shutting down")
		_ = s.Shutdown()
	}()

	logger.Info("server started")
	if err = s.Listen(); err != nil {
		logger.Fatal(err)
	}

	logger.Info("server shut down")
}
