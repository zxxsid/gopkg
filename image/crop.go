package image

import (
	"image"
	"image/draw"
)

// Anchor specifies the anchor point for cropping
type Anchor int

const (
	AnchorCenter Anchor = iota
	AnchorTopLeft
	AnchorTop
	AnchorTopRight
	AnchorLeft
	AnchorRight
	AnchorBottomLeft
	AnchorBottom
	AnchorBottomRight
)

// Crop extracts a rectangular region starting at (x, y) with given width and height
func (i *Image) Crop(x, y, width, height int) *Image {
	cropRect := image.Rect(x, y, x+width, y+height)
	cropRect = cropRect.Intersect(i.img.Bounds())
	if cropRect.Empty() {
		return i.Clone()
	}

	dst := image.NewRGBA(image.Rect(0, 0, cropRect.Dx(), cropRect.Dy()))
	draw.Draw(dst, dst.Bounds(), i.img, cropRect.Min, draw.Src)
	return &Image{img: dst, format: i.format}
}

// CropCenter crops a region of the given size from the center of the image
func (i *Image) CropCenter(width, height int) *Image {
	bounds := i.img.Bounds()
	cx := (bounds.Min.X + bounds.Max.X) / 2
	cy := (bounds.Min.Y + bounds.Max.Y) / 2
	x := cx - width/2
	y := cy - height/2
	return i.Crop(x, y, width, height)
}

// CropAnchor crops a region of the given size from the specified anchor point
func (i *Image) CropAnchor(width, height int, anchor Anchor) *Image {
	bounds := i.img.Bounds()
	srcW := bounds.Dx()
	srcH := bounds.Dy()

	var x, y int

	switch anchor {
	case AnchorTopLeft:
		x, y = bounds.Min.X, bounds.Min.Y
	case AnchorTop:
		x = bounds.Min.X + (srcW-width)/2
		y = bounds.Min.Y
	case AnchorTopRight:
		x = bounds.Max.X - width
		y = bounds.Min.Y
	case AnchorLeft:
		x = bounds.Min.X
		y = bounds.Min.Y + (srcH-height)/2
	case AnchorCenter:
		x = bounds.Min.X + (srcW-width)/2
		y = bounds.Min.Y + (srcH-height)/2
	case AnchorRight:
		x = bounds.Max.X - width
		y = bounds.Min.Y + (srcH-height)/2
	case AnchorBottomLeft:
		x = bounds.Min.X
		y = bounds.Max.Y - height
	case AnchorBottom:
		x = bounds.Min.X + (srcW-width)/2
		y = bounds.Max.Y - height
	case AnchorBottomRight:
		x = bounds.Max.X - width
		y = bounds.Max.Y - height
	default:
		x = bounds.Min.X + (srcW-width)/2
		y = bounds.Min.Y + (srcH-height)/2
	}

	return i.Crop(x, y, width, height)
}
