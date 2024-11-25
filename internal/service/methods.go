package service

import (
	"math/rand"
	"sync"
	"time"
)

// genKey - генерация подстроки со случайной последовательностью символов как сокращенный URL.
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

// genStrs - для реализации Fan-In.
func genStrs(strs ...string) chan string {
	outCh := make(chan string)
	go func() {
		defer close(outCh)
		for _, s := range strs {
			outCh <- s
		}
	}()
	return outCh
}

// genStrs - для реализации Fan-In.
func pushStr(inCh chan string) chan string {
	outCh := make(chan string)
	go func() {
		defer close(outCh)
		for s := range inCh {
			outCh <- s
		}
	}()
	return outCh
}

// genStrs - реализация Fan-In.
func fanInStr(chs ...chan string) chan string {
	outCh := make(chan string)
	var wg sync.WaitGroup
	output := func(c chan string) {
		for s := range c {
			outCh <- s
		}
		wg.Done()
	}
	wg.Add(len(chs))
	for _, c := range chs {
		go output(c)
	}
	go func() {
		wg.Wait()
		close(outCh)
	}()
	return outCh
}
