package rest

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"

	"github.com/lRhythm/shortener/internal/config"
	"github.com/lRhythm/shortener/internal/logs"
	"github.com/lRhythm/shortener/internal/service"
	"github.com/lRhythm/shortener/internal/storage"
)

func ExampleServer_PingHandler() {
	cfg, _ := config.New()
	db, _ := storage.NewMemory(cfg.File())
	s, _ := New(logs.New(), cfg, service.New(service.WithStorage(db)))
	f := fiber.New()
	f.Get("/ping", s.PingHandler)
	log.Fatal(s.Listen())
}

func ExampleServer_GetHandler() {
	cfg, _ := config.New()
	db, _ := storage.NewMemory(cfg.File())
	s, _ := New(logs.New(), cfg, service.New(service.WithStorage(db)))
	f := fiber.New()
	f.Get(fmt.Sprintf("/:%s", pathParamID), s.GetHandler)
	log.Fatal(s.Listen())
}

func ExampleServer_APIUserUrlsGetHandler() {
	cfg, _ := config.New()
	db, _ := storage.NewMemory(cfg.File())
	s, _ := New(logs.New(), cfg, service.New(service.WithStorage(db)))
	f := fiber.New()
	f.Get("/api/user/urls", s.authenticateMiddleware, s.APIUserUrlsGetHandler)
	log.Fatal(s.Listen())
}

func ExampleServer_APIUserUrlsDeleteHandler() {
	cfg, _ := config.New()
	db, _ := storage.NewMemory(cfg.File())
	s, _ := New(logs.New(), cfg, service.New(service.WithStorage(db)))
	f := fiber.New()
	f.Delete("/api/user/urls", s.authenticateMiddleware, s.APIUserUrlsDeleteHandler)
	log.Fatal(s.Listen())
}

func ExampleServer_APICreateHandler() {
	cfg, _ := config.New()
	db, _ := storage.NewMemory(cfg.File())
	s, _ := New(logs.New(), cfg, service.New(service.WithStorage(db)))
	f := fiber.New()
	f.Post("/api/shorten", s.registerMiddleware, s.APICreateHandler)
	log.Fatal(s.Listen())
}

func ExampleServer_APICreateBatchHandler() {
	cfg, _ := config.New()
	db, _ := storage.NewMemory(cfg.File())
	s, _ := New(logs.New(), cfg, service.New(service.WithStorage(db)))
	f := fiber.New()
	f.Post("/api/shorten/batch", s.registerMiddleware, s.APICreateBatchHandler)
	log.Fatal(s.Listen())
}

func ExampleServer_CreateHandler() {
	cfg, _ := config.New()
	db, _ := storage.NewMemory(cfg.File())
	s, _ := New(logs.New(), cfg, service.New(service.WithStorage(db)))
	f := fiber.New()
	f.Post("/", s.registerMiddleware, s.CreateHandler)
	log.Fatal(s.Listen())
}
