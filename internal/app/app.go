package app

import (
	"github.com/lRhythm/shortener/internal/server"
	"github.com/lRhythm/shortener/internal/storage"
	"log"
)

func Start() {
	log.Fatal(server.New(storage.NewInMemory()).Listen())
}
