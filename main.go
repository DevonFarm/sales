package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"

	"github.com/DevonFarm/sales/auth"
	"github.com/DevonFarm/sales/database"
	"github.com/DevonFarm/sales/horse"
)

//go:embed templates assets
var templates embed.FS

func runServer() error {
	godotenv.Load()
	connString := os.Getenv("COCKROACH_DSN")
	if connString == "" {
		return fmt.Errorf("missing COCKROACH_DSN env var")
	}
	db, err := database.NewDBConn(connString)
	if err != nil {
		return fmt.Errorf("failed to connect to db: %w", err)
	}
	defer db.Close(context.Background())

	fs := http.FS(templates)
	engine := html.NewFileSystem(fs, ".html")
	app := fiber.New(fiber.Config{
		Views:       engine,
		ViewsLayout: "templates/layouts/main.html",
	})
	app.Use(logger.New())

	// Serve static assets from embedded filesystem
	app.Use("/assets", filesystem.New(filesystem.Config{
		Root:       fs,
		PathPrefix: "assets",
	}))

	horse.Routes(app, db)
	// Auth routes (Stytch magic links)
	stytch, err := auth.NewStytchFromEnv()
	if err != nil {
		return fmt.Errorf("stytch failed to configure: %w", err)
	}
	stytch.Register(app)
	return app.Listen(":4242")
}

func main() {
	if err := runServer(); err != nil {
		log.Fatal(err)
	}
}
