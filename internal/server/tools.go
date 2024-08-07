package server

import (
	"math/rand"
	"time"
)

func (s *Server) genID() string {
	var (
		charset    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
		charsetLen = len(charset)
		IDLen      = 8
	)
	rand.New(rand.NewSource(time.Now().UnixNano()))
	ID := make([]byte, IDLen)
	for i := range ID {
		ID[i] = charset[rand.Intn(charsetLen)]
	}
	return string(ID)
}
