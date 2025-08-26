package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
	"github.com/otiai10/copy"

	"github.com/DevonFarm/sales/database"
	"github.com/DevonFarm/sales/horse"
)

//go:embed templates assets
var templates embed.FS

// builds static content
func build() error {
	t, err := template.ParseGlob("templates/**.html")
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}
	outDir := "output"
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return fmt.Errorf("failed to create output dir: %w", err)
	}
	f, err := os.Create(fmt.Sprintf("%s/index.html", outDir))
	if err != nil {
		return fmt.Errorf("failed to create index.html: %w", err)
	}
	if err := t.ExecuteTemplate(f, "index.html", nil); err != nil {
		return fmt.Errorf("failed to execute index template: %w", err)
	}
	if err := buildHorses(t, outDir); err != nil {
		return fmt.Errorf("failed to build horses templates: %w", err)
	}
	if err := copy.Copy("assets", fmt.Sprintf("%s/assets", outDir)); err != nil {
		return fmt.Errorf("failed to copy assets dir: %w", err)
	}
	return nil
}

func run() error {
	var dev bool
	flag.BoolVar(&dev, "dev", false, "Run live dev server")
	if err := build(); err != nil {
		return fmt.Errorf("failed to build static content: %w", err)
	}
	flag.Parse()
	if dev {
		if err := runHTTPServer(); err != nil {
			return fmt.Errorf("failed to run HTTP server: %w", err)
		}
	}
	return nil
}

func runHTTPServer() error {
	addr := ":4242"
	fs := http.FileServer(http.Dir("./output"))
	http.Handle("/", fs)
	fmt.Println("Starting HTTP server on", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		return fmt.Errorf("failed to run HTTP server: %w", err)
	}
	return nil
}

func runAPI() error {
	if err := godotenv.Load(); err != nil {
		log.Printf("failed to load .env: %v", err)
	}
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
		Views: engine,
	})

	// Serve static assets from embedded filesystem
	app.Use("/assets", filesystem.New(filesystem.Config{
		Root:       fs,
		PathPrefix: "assets",
	}))

	horse.Routes(app, db)

	// h := horse.NewHorse(
	// 	"Test Horse",
	// 	"A horse used for testing",
	// 	time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
	// 	horse.GenderStallion,
	// )
	// if err := h.Save(context.Background(), db); err != nil {
	// 	return fmt.Errorf("failed to save horse: %w", err)
	// }
	return app.Listen(":4242")
}

func main() {
	if err := runAPI(); err != nil {
		log.Fatal(err)
	}
}

func buildHorses(tmplGlob *template.Template, outDir string) error {
	if err := os.MkdirAll(fmt.Sprintf("%s/horses", outDir), 0755); err != nil {
		return fmt.Errorf("failed to create horses dir: %w", err)
	}

	// temporary hardcoded horse
	h := horse.Horse{Name: "Link"}
	h.NewImage(
		"link_headshot.jpeg",
		"link_headshot_thumb.jpeg",
		"Headshot of a silver bay Gypsian colt",
	)
	h.NewImage(
		"link_standing.jpeg",
		"link_standing_thumb.jpeg",
		"A silver bay Gypsian colt standing in a field",
	)
	h.NewImage(
		"link_stepping.jpeg",
		"link_stepping_thumb.jpeg",
		"A silver bay Gypsian colt, appearing to step near a Norwegian Fjord mare with her ears pinned",
	)
	h.NewImage(
		"link_trotting.jpeg",
		"link_trotting_thumb.jpeg",
		"A silver bay Gypsian colt trotting in a field",
	)
	h.NewImage(
		"link_standing.jpeg",
		"link_standing_thumb.jpeg",
		"A silver bay Gypsian colt standing in a field",
	)
	h.NewImage(
		"link_stepping.jpeg",
		"link_stepping_thumb.jpeg",
		"A silver bay Gypsian colt, appearing to step near a Norwegian Fjord mare with her ears pinned",
	)
	h.NewImage(
		"link_trotting.jpeg",
		"link_trotting_thumb.jpeg",
		"A silver bay Gypsian colt trotting in a field",
	)
	filename := h.HTMLPath()
	f, err := os.Create(fmt.Sprintf("%s/%s", outDir, filename))
	if err != nil {
		return fmt.Errorf("failed to create %s: %w", filename, err)
	}
	if err := tmplGlob.ExecuteTemplate(f, "horse.html", h); err != nil {
		return fmt.Errorf("failed to execute horse template: %w", err)
	}
	return nil
}
