package server

import (
	"embed"
	"fmt"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"

	"github.com/DevonFarm/sales/auth"
	"github.com/DevonFarm/sales/database"
)

type Server struct {
	App  *fiber.App
	DB   *database.DB
	Auth *auth.StytchAuth
}

// templateFS must contain the "templates" and "assets" directories and
// "templates/layouts/main.html" must exist.
func NewServer(templateFS embed.FS) (*Server, error) {
	godotenv.Load()
	connString := os.Getenv("COCKROACH_DSN")
	if connString == "" {
		return nil, fmt.Errorf("missing COCKROACH_DSN env var")
	}
	db, err := database.NewDBConn(connString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}

	fs := http.FS(templateFS)
	engine := html.NewFileSystem(fs, ".html")
	app := fiber.New(fiber.Config{
		Views:       engine,
		ViewsLayout: "templates/layouts/main.html",
	})
	app.Use(logger.New())

	// Auth routes
	stytch, err := auth.NewStytchFromEnv()
	if err != nil {
		return nil, fmt.Errorf("stytch failed to configure: %w", err)
	}
	stytch.Register(app, db)

	// Serve static assets from embedded filesystem
	app.Use("/assets", filesystem.New(filesystem.Config{
		Root:       fs,
		PathPrefix: "assets",
	}))

	return &Server{
		App:  app,
		DB:   db,
		Auth: stytch,
	}, nil
}

func (s *Server) Listen(addr string) error {
	return s.App.Listen(addr)
}
