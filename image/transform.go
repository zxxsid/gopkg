package image

import (
	"image"
	"image/color"
	"image/draw"
)

// Rotate90 rotates the image 90 degrees clockwise
func (i *Image) Rotate90() *Image {
	bounds := i.img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	dst := image.NewRGBA(image.Rect(0, 0, h, w))

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			dst.Set(h-1-(y-bounds.Min.Y), x-bounds.Min.X, i.img.At(x, y))
		}
	}

	return &Image{img: dst, format: i.format}
}

// Rotate180 rotates the image 180 degrees
func (i *Image) Rotate180() *Image {
	bounds := i.img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	dst := image.NewRGBA(image.Rect(0, 0, w, h))

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			dst.Set(w-1-(x-bounds.Min.X), h-1-(y-bounds.Min.Y), i.img.At(x, y))
		}
	}

	return &Image{img: dst, format: i.format}
}

// Rotate270 rotates the image 270 degrees clockwise (90 degrees counter-clockwise)
func (i *Image) Rotate270() *Image {
	bounds := i.img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	dst := image.NewRGBA(image.Rect(0, 0, h, w))

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			dst.Set(y-bounds.Min.Y, w-1-(x-bounds.Min.X), i.img.At(x, y))
		}
	}

	return &Image{img: dst, format: i.format}
}

// FlipH flips the image horizontally
func (i *Image) FlipH() *Image {
	bounds := i.img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	dst := image.NewRGBA(image.Rect(0, 0, w, h))

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			dst.Set(w-1-x, y, i.img.At(x+bounds.Min.X, y+bounds.Min.Y))
		}
	}

	return &Image{img: dst, format: i.format}
}

// FlipV flips the image vertically
func (i *Image) FlipV() *Image {
	bounds := i.img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	dst := image.NewRGBA(image.Rect(0, 0, w, h))

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			dst.Set(x, h-1-y, i.img.At(x+bounds.Min.X, y+bounds.Min.Y))
		}
	}

	return &Image{img: dst, format: i.format}
}

// Grayscale converts the image to grayscale
func (i *Image) Grayscale() *Image {
	bounds := i.img.Bounds()
	gray := image.NewGray(bounds)
	draw.Draw(gray, bounds, i.img, bounds.Min, draw.Src)
	return &Image{img: gray, format: i.format}
}

// Opacity returns a copy with the alpha channel set to the given value (0.0 ~ 1.0)
func (i *Image) Opacity(alpha float64) *Image {
	if alpha < 0 {
		alpha = 0
	}
	if alpha > 1 {
		alpha = 1
	}

	bounds := i.img.Bounds()
	dst := image.NewNRGBA(bounds)
	a := uint8(alpha * 255)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := i.img.At(x, y).RGBA()
			dst.SetNRGBA(x, y, color.NRGBA{
				R: uint8(r >> 8),
				G: uint8(g >> 8),
				B: uint8(b >> 8),
				A: a,
			})
		}
	}

	return &Image{img: dst, format: i.format}
}

// Watermark overlays another image at the specified position
func (i *Image) Watermark(mark *Image, x, y int) *Image {
	bounds := i.img.Bounds()
	dst := image.NewRGBA(bounds)
	draw.Draw(dst, bounds, i.img, bounds.Min, draw.Src)

	markBounds := mark.img.Bounds()
	offset := image.Pt(x, y)
	draw.Draw(dst, markBounds.Add(offset), mark.img, markBounds.Min, draw.Over)

	return &Image{img: dst, format: i.format}
}
