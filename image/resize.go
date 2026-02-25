package image

import (
	"image"

	xdraw "golang.org/x/image/draw"
)

// Resize scales the image to the given dimensions.
// If width or height is 0, it is calculated to preserve the aspect ratio.
func (i *Image) Resize(width, height int) *Image {
	if width <= 0 && height <= 0 {
		return i.Clone()
	}

	srcBounds := i.img.Bounds()
	srcW := srcBounds.Dx()
	srcH := srcBounds.Dy()

	if width <= 0 {
		width = height * srcW / srcH
	}
	if height <= 0 {
		height = width * srcH / srcW
	}
	if width < 1 {
		width = 1
	}
	if height < 1 {
		height = 1
	}

	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	xdraw.CatmullRom.Scale(dst, dst.Bounds(), i.img, srcBounds, xdraw.Over, nil)
	return &Image{img: dst, format: i.format}
}

// ResizeByWidth resizes the image to the given width, preserving the aspect ratio
func (i *Image) ResizeByWidth(width int) *Image {
	return i.Resize(width, 0)
}

// ResizeByHeight resizes the image to the given height, preserving the aspect ratio
func (i *Image) ResizeByHeight(height int) *Image {
	return i.Resize(0, height)
}

// Thumbnail generates a thumbnail that fits within maxWidth x maxHeight,
// preserving the aspect ratio and never upscaling
func (i *Image) Thumbnail(maxWidth, maxHeight int) *Image {
	return i.ResizeWithFit(maxWidth, maxHeight, ImageFitInside)
}

// ResizeWithFit resizes the image according to the specified fit mode
func (i *Image) ResizeWithFit(width, height int, fit string) *Image {
	srcBounds := i.img.Bounds()
	srcW := float64(srcBounds.Dx())
	srcH := float64(srcBounds.Dy())
	targetW := float64(width)
	targetH := float64(height)

	ratioW := targetW / srcW
	ratioH := targetH / srcH

	var newW, newH float64

	switch fit {
	case ImageFitContain:
		ratio := min(ratioW, ratioH)
		newW = srcW * ratio
		newH = srcH * ratio

	case ImageFitCover:
		ratio := max(ratioW, ratioH)
		newW = srcW * ratio
		newH = srcH * ratio

	case ImageFitFill:
		newW = targetW
		newH = targetH

	case ImageFitInside:
		ratio := min(ratioW, ratioH)
		if ratio > 1 {
			ratio = 1
		}
		newW = srcW * ratio
		newH = srcH * ratio

	case ImageFitOutside:
		ratio := max(ratioW, ratioH)
		if ratio < 1 {
			ratio = 1
		}
		newW = srcW * ratio
		newH = srcH * ratio

	default:
		ratio := min(ratioW, ratioH)
		newW = srcW * ratio
		newH = srcH * ratio
	}

	finalW := int(newW)
	finalH := int(newH)
	if finalW < 1 {
		finalW = 1
	}
	if finalH < 1 {
		finalH = 1
	}

	dst := image.NewRGBA(image.Rect(0, 0, finalW, finalH))
	xdraw.CatmullRom.Scale(dst, dst.Bounds(), i.img, srcBounds, xdraw.Over, nil)
	return &Image{img: dst, format: i.format}
}
