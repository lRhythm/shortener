package app

import (
	"github.com/lRhythm/shortener/internal/service"
	"github.com/lRhythm/shortener/internal/storage"
	"github.com/lRhythm/shortener/internal/transport/httptransport"
	"log"
)

func Start() {
	s, e := httptransport.New(
		service.New(
			service.WithStorage(storage.NewInMemory()),
		),
	)
	if e != nil {
		log.Fatal(e)
	}
	log.Fatal(s.Listen())
}
