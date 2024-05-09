package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

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
