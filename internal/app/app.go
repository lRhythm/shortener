/*
Package app - создание необходимых экземпляров пакетов, конфигурирование и запуск сервиса.
*/
package app

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/lRhythm/shortener/internal/config"
	"github.com/lRhythm/shortener/internal/logs"
	"github.com/lRhythm/shortener/internal/service"
	"github.com/lRhythm/shortener/internal/storage"
	"github.com/lRhythm/shortener/internal/transport/grpc"
	"github.com/lRhythm/shortener/internal/transport/rest"
)

// StartREST - запуск сервиса для взаимодействия с помощью REST api.
func StartREST() {
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
		logger.Info("rest server shutting down")
		if err = s.Shutdown(); err != nil {
			logger.Info("shutting down err:", err)
		}
	}()

	logger.Info("rest server started")
	if err = s.Listen(); err != nil {
		logger.Fatal(err)
	}

	logger.Info("rest server shut down")
}

// StartGRPC - запуск сервиса для взаимодействия с помощью gRPC.
func StartGRPC() {
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

	s, err := grpc.New(
		cfg,
		service.New(
			service.WithStorage(store),
		),
	)
	if err != nil {
		logger.Fatal(err)
	}

	listen, err := net.Listen("tcp", cfg.Host())
	if err != nil {
		logger.Fatal(err)
	}

	logger.Info("grpc server started")

	if err = s.Serve(listen); err != nil {
		logger.Fatal(err)
	}
}
