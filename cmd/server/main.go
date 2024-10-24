package main

import (
	"context"

	"log"
	"net/http"
	_ "net/http/pprof"

	"gophKeeper/internal/server/app"

	_ "github.com/lib/pq"
)

var buildVersion string
var buildDate string
var buildCommit string

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil)) // Запустите сервер pprof
	}()
	app.RunApp(context.Background(), nil, nil,
		app.BuildMetadata{Version: buildVersion, Date: buildDate, Commit: buildCommit})
}
