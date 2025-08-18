package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/otiai10/copy"

	"github.com/DevonFarm/sales/horse"
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
		fmt.Println("Starting dev server on", addr)
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
