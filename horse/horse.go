package horse

import (
	"fmt"

	"github.com/iancoleman/strcase"
)

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
