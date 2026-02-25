package image

import (
	"image"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"
)

// Fit mode constants
const (
	ImageFitContain = "contain"
	ImageFitCover   = "cover"
	ImageFitFill    = "fill"
	ImageFitInside  = "inside"
	ImageFitOutside = "outside"
)

// Format constants
const (
	FormatJPEG = "jpeg"
	FormatPNG  = "png"
	FormatGIF  = "gif"
	FormatBMP  = "bmp"
	FormatTIFF = "tiff"
	FormatWEBP = "webp"
)

// Image provides image processing capabilities
type Image struct {
	img    image.Image
	format string
}

// NewImage decodes an image from an io.Reader
func NewImage(r io.Reader) (*Image, error) {
	img, format, err := image.Decode(r)
	if err != nil {
		return nil, err
	}
	return &Image{img: img, format: format}, nil
}

// Open reads and decodes an image from a file path
func Open(path string) (*Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return NewImage(f)
}

// FromImage creates an Image from a standard library image.Image
func FromImage(img image.Image, format string) *Image {
	return &Image{img: img, format: format}
}

// Format returns the decoded format name
func (i *Image) Format() string {
	return i.format
}

// Width returns the image width in pixels
func (i *Image) Width() int {
	return i.img.Bounds().Dx()
}

// Height returns the image height in pixels
func (i *Image) Height() int {
	return i.img.Bounds().Dy()
}

// Bounds returns the image bounding rectangle
func (i *Image) Bounds() image.Rectangle {
	return i.img.Bounds()
}

// RawImage returns the underlying image.Image
func (i *Image) RawImage() image.Image {
	return i.img
}

// Clone creates a deep copy of the image
func (i *Image) Clone() *Image {
	bounds := i.img.Bounds()
	dst := image.NewRGBA(bounds)
	draw.Draw(dst, bounds, i.img, bounds.Min, draw.Src)
	return &Image{img: dst, format: i.format}
}
