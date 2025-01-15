/*
Package main - HTTP/2 (gRPC) сервис сокращения URL.
*/
package main

import "github.com/lRhythm/shortener/internal/app"

func main() {
	app.StartGRPC()
}
