package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/iancoleman/strcase"
	"github.com/otiai10/copy"
)

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
		addr := ":4242"
		fs := http.FileServer(http.Dir("./output"))
		http.Handle("/", fs)
		if err := http.ListenAndServe(addr, nil); err != nil {
			return fmt.Errorf("failed to run dev server: %w", err)
		}
	}
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func buildHorses(tmplGlob *template.Template, outDir string) error {
	if err := os.MkdirAll(fmt.Sprintf("%s/horses", outDir), 0755); err != nil {
		return fmt.Errorf("failed to create horses dir: %w", err)
	}

	// temporary hardcoded horse
	horse := Horse{Name: "Link"}
	horse.NewImage(
		"link_headshot.jpeg",
		"link_headshot_thumb.jpeg",
		"Headshot of a silver bay Gypsian colt",
	)
	horse.NewImage(
		"link_standing.jpeg",
		"link_standing_thumb.jpeg",
		"A silver bay Gypsian colt standing in a field",
	)
	horse.NewImage(
		"link_stepping.jpeg",
		"link_stepping_thumb.jpeg",
		"A silver bay Gypsian colt, appearing to step near a Norwegian Fjord mare with her ears pinned",
	)
	horse.NewImage(
		"link_trotting.jpeg",
		"link_trotting_thumb.jpeg",
		"A silver bay Gypsian colt trotting in a field",
	)
	filename := horse.HTMLPath()
	f, err := os.Create(fmt.Sprintf("%s/%s", outDir, filename))
	if err != nil {
		return fmt.Errorf("failed to create %s: %w", filename, err)
	}
	if err := tmplGlob.ExecuteTemplate(f, "horse.html", horse); err != nil {
		return fmt.Errorf("failed to execute horse template: %w", err)
	}
	return nil
}

type Horse struct {
	Name   string
	Images []*Image
}

func (h *Horse) HTMLPath() string {
	return fmt.Sprintf("%s.html", strcase.ToSnake(h.Name))
}

func (h *Horse) NewImage(full, thumbnail, alt string) {
	prefix := fmt.Sprintf("assets/images/horses/%s", h.Name)
	img := &Image{
		Full:      fmt.Sprintf("%s/%s", prefix, full),
		Thumbnail: fmt.Sprintf("%s/%s", prefix, thumbnail),
		Alt:       alt,
	}
	if h.Images == nil {
		h.Images = make([]*Image, 0)
	}
	h.Images = append(h.Images, img)
}

type Image struct {
	Full      string
	Alt       string
	Thumbnail string
}
