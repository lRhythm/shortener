package service

import (
	"math/rand"
	"time"
)

func (c *Client) genKey() string {
	var (
		charset    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
		charsetLen = len(charset)
		keyLen     = 8
	)
	rand.New(rand.NewSource(time.Now().UnixNano()))
	key := make([]byte, keyLen)
	for i := range key {
		key[i] = charset[rand.Intn(charsetLen)]
	}
	return string(key)
}
