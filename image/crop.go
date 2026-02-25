package image

import (
	"image"
	"image/draw"
)

// Anchor 裁剪锚点类型，用于指定 CropAnchor 方法的裁剪起始位置。
// 支持九宫格的九个锚点位置（左上、上、右上、左、中、右、左下、下、右下）。
type Anchor int

// 九宫格锚点位置常量
const (
	AnchorCenter      Anchor = iota // 中心点
	AnchorTopLeft                   // 左上角
	AnchorTop                       // 上方居中
	AnchorTopRight                  // 右上角
	AnchorLeft                      // 左侧居中
	AnchorRight                     // 右侧居中
	AnchorBottomLeft                // 左下角
	AnchorBottom                    // 下方居中
	AnchorBottomRight               // 右下角
)

// Crop 从图片的指定起点 (x, y) 开始，裁剪出 width x height 大小的矩形区域。
// 如果裁剪区域超出图片边界，会自动裁切到图片范围内（取交集）。
// 如果裁剪区域与图片完全不重叠，则返回原图的深拷贝。
func (i *Image) Crop(x, y, width, height int) *Image {
	// 构造裁剪矩形并与图片边界取交集，防止越界
	cropRect := image.Rect(x, y, x+width, y+height)
	cropRect = cropRect.Intersect(i.img.Bounds())
	if cropRect.Empty() {
		return i.Clone()
	}

	// 创建新图片并将裁剪区域的像素复制过去
	dst := image.NewRGBA(image.Rect(0, 0, cropRect.Dx(), cropRect.Dy()))
	draw.Draw(dst, dst.Bounds(), i.img, cropRect.Min, draw.Src)
	return &Image{img: dst, format: i.format}
}

// CropCenter 以图片的中心点为基准，向四周各扩展一半的宽高来裁剪出指定大小的区域。
// 适用于需要从图片正中央提取内容的场景，如头像裁剪、焦点区域提取等。
// 如果指定的宽高大于图片尺寸，会自动裁切到图片范围内。
func (i *Image) CropCenter(width, height int) *Image {
	bounds := i.img.Bounds()
	// 计算图片中心点坐标
	cx := (bounds.Min.X + bounds.Max.X) / 2
	cy := (bounds.Min.Y + bounds.Max.Y) / 2
	// 从中心点向左上方偏移，计算裁剪起点
	x := cx - width/2
	y := cy - height/2
	return i.Crop(x, y, width, height)
}

// CropAnchor 根据指定的锚点位置，裁剪出 width x height 大小的矩形区域。
// 锚点决定了裁剪区域在原图中的对齐方式，支持九宫格的九个位置。
// 例如 AnchorTopLeft 从左上角开始裁剪，AnchorCenter 从中心裁剪，
// AnchorBottomRight 从右下角开始裁剪。
func (i *Image) CropAnchor(width, height int, anchor Anchor) *Image {
	bounds := i.img.Bounds()
	srcW := bounds.Dx()
	srcH := bounds.Dy()

	var x, y int

	// 根据锚点类型计算裁剪区域的左上角坐标
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
		// 未知锚点默认使用中心点
		x = bounds.Min.X + (srcW-width)/2
		y = bounds.Min.Y + (srcH-height)/2
	}

	return i.Crop(x, y, width, height)
}
