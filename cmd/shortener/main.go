/*
Package main - HTTP сервис сокращения URL.
*/
package main

import (
	"fmt"

	"github.com/lRhythm/shortener/internal/app"
)

const na = "N/A"

// Флаги линковщика.
var buildVersion, buildDate, buildCommit = na, na, na

func main() {
	// Пример запуска из корневой директории проекта:
	// go run -ldflags "-X main.buildVersion=v1.0.1 -X 'main.buildDate=$(date +'%Y/%m/%d %H:%M:%S')'" ./cmd/shortener/main.go
	fmt.Printf("Build version: %s\nBuild date: %s\nBuild commit: %s\n", buildVersion, buildDate, buildCommit)

	app.Start()
}
