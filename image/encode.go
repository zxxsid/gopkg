package image

import (
	"bytes"
	"fmt"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/image/bmp"
	"golang.org/x/image/tiff"
)

// EncodeOptions controls encoding parameters
type EncodeOptions struct {
	Quality int
}

// EncodeOption is a functional option for encoding
type EncodeOption func(*EncodeOptions)

// Quality sets the JPEG encoding quality (1-100, default 85)
func Quality(q int) EncodeOption {
	return func(opts *EncodeOptions) {
		opts.Quality = q
	}
}

func defaultEncodeOptions() *EncodeOptions {
	return &EncodeOptions{Quality: 85}
}

// Encode writes the image to w in the specified format
func (i *Image) Encode(w io.Writer, format string, opts ...EncodeOption) error {
	o := defaultEncodeOptions()
	for _, opt := range opts {
		opt(o)
	}

	format = strings.ToLower(format)
	switch format {
	case FormatJPEG, "jpg":
		return jpeg.Encode(w, i.img, &jpeg.Options{Quality: o.Quality})
	case FormatPNG:
		return png.Encode(w, i.img)
	case FormatGIF:
		return gif.Encode(w, i.img, nil)
	case FormatBMP:
		return bmp.Encode(w, i.img)
	case FormatTIFF:
		return tiff.Encode(w, i.img, nil)
	default:
		return fmt.Errorf("unsupported encoding format: %s", format)
	}
}

// Save encodes the image and writes it to a file, format is inferred from the file extension
func (i *Image) Save(path string, opts ...EncodeOption) error {
	format := formatFromPath(path)
	if format == "" {
		return fmt.Errorf("cannot determine format from file extension: %s", path)
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return i.Encode(f, format, opts...)
}

// ToBytes encodes the image to a byte slice in the specified format
func (i *Image) ToBytes(format string, opts ...EncodeOption) ([]byte, error) {
	var buf bytes.Buffer
	if err := i.Encode(&buf, format, opts...); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func formatFromPath(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".jpg", ".jpeg":
		return FormatJPEG
	case ".png":
		return FormatPNG
	case ".gif":
		return FormatGIF
	case ".bmp":
		return FormatBMP
	case ".tiff", ".tif":
		return FormatTIFF
	default:
		return ""
	}
}
