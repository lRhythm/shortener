package app

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var storage map[string]string

func Start() {
	storage = make(map[string]string)
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, handlerResolver)
	log.Fatal(http.ListenAndServe(`localhost:8080`, mux))
}

func handlerResolver(resp http.ResponseWriter, req *http.Request) {
	m := req.Method
	pathID := strings.TrimPrefix(req.URL.Path, "/")
	if m == http.MethodGet && len(pathID) > 0 {
		getURL(resp, pathID)
		return
	}
	if m == http.MethodPost && len(pathID) == 0 {
		createURL(resp, req)
		return
	}
	resp.WriteHeader(http.StatusBadRequest)
}

func getURL(resp http.ResponseWriter, pathID string) {
	u, ok := storage[pathID]
	if !ok {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}
	resp.Header().Set("Location", u)
	resp.WriteHeader(http.StatusTemporaryRedirect)
}

func createURL(resp http.ResponseWriter, req *http.Request) {
	b, e := io.ReadAll(req.Body)
	if e != nil {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}
	u := string(b)
	_, e = url.ParseRequestURI(u)
	if e != nil {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}
	k := genID()
	storage[k] = u
	resp.WriteHeader(http.StatusCreated)
	_, _ = resp.Write([]byte(fmt.Sprintf("http://%s/%s", req.Host, k)))
}

func genID() string {
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
