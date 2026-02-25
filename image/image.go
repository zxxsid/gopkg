package image

import (
	"image"
	"io"
)

const (
	ImageFitContain = "contain"
	ImageFitCover   = "cover"
	ImageFitFill    = "fill"
	ImageFitInside  = "inside"
	ImageFitOutside = "outside"
)

type Image struct {
	src    image.Image
	dst    image.Image
	format string
	w, h   int
}

func NewImage(r io.Reader) (*Image, error) {
	img, format, err := image.Decode(r)
	if err != nil {
		return nil, err
	}
	return &Image{src: img, format: format}, nil
}

