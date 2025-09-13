package main

import (
	"embed"
	"log"

	"github.com/DevonFarm/sales/farm"
	"github.com/DevonFarm/sales/horse"
	"github.com/DevonFarm/sales/server"
)

//go:embed templates assets
var templates embed.FS

func runServer() error {
	srvr, err := server.NewServer(templates)
	if err != nil {
		return err
	}

	horse.RegisterRoutes(srvr.App, srvr.DB, srvr.Auth)
	farm.RegisterRoutes(srvr.App, srvr.DB)

	return srvr.Listen(":4242")
}

func main() {
	if err := runServer(); err != nil {
		log.Fatal(err)
	}
}
