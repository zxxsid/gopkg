package image

import (
	"image"

	xdraw "golang.org/x/image/draw" // 使用 x/image/draw 提供高质量的图片缩放算法
)

// Resize 将图片缩放到指定的宽度和高度。
// 使用 CatmullRom 插值算法（三次插值），在缩放质量和性能之间取得良好平衡。
// 如果 width 或 height 为 0，则根据另一个维度按原始宽高比自动计算。
// 如果两个维度都 <= 0，则返回原图的深拷贝。
func (i *Image) Resize(width, height int) *Image {
	if width <= 0 && height <= 0 {
		return i.Clone()
	}

	srcBounds := i.img.Bounds()
	srcW := srcBounds.Dx()
	srcH := srcBounds.Dy()

	// 当 width 为 0 时，根据目标高度和原始宽高比计算宽度
	if width <= 0 {
		width = height * srcW / srcH
	}
	// 当 height 为 0 时，根据目标宽度和原始宽高比计算高度
	if height <= 0 {
		height = width * srcH / srcW
	}
	// 确保最终尺寸至少为 1 像素
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

// ResizeByWidth 按指定宽度等比缩放图片。
// 高度根据原始宽高比自动计算，保持图片不变形。
func (i *Image) ResizeByWidth(width int) *Image {
	return i.Resize(width, 0)
}

// ResizeByHeight 按指定高度等比缩放图片。
// 宽度根据原始宽高比自动计算，保持图片不变形。
func (i *Image) ResizeByHeight(height int) *Image {
	return i.Resize(0, height)
}

// Thumbnail 生成缩略图，使图片完整地包含在 maxWidth x maxHeight 的范围内。
// 保持原始宽高比，且不会将小图放大（使用 ImageFitInside 模式）。
// 适用于生成列表页缩略图、头像预览等场景。
func (i *Image) Thumbnail(maxWidth, maxHeight int) *Image {
	return i.ResizeWithFit(maxWidth, maxHeight, ImageFitInside)
}

// ResizeWithFit 按指定的适配模式缩放图片到目标尺寸。
// 支持五种适配模式：
//   - ImageFitContain: 保持宽高比，图片完全包含在目标尺寸内（可能有留白区域）
//   - ImageFitCover:   保持宽高比，图片完全覆盖目标尺寸（可能超出目标区域）
//   - ImageFitFill:    不保持宽高比，强制拉伸到目标尺寸
//   - ImageFitInside:  类似 contain，但不会放大小于目标尺寸的图片
//   - ImageFitOutside: 类似 cover，但不会缩小大于目标尺寸的图片
func (i *Image) ResizeWithFit(width, height int, fit string) *Image {
	srcBounds := i.img.Bounds()
	srcW := float64(srcBounds.Dx())
	srcH := float64(srcBounds.Dy())
	targetW := float64(width)
	targetH := float64(height)

	// 分别计算宽度和高度的缩放比例
	ratioW := targetW / srcW
	ratioH := targetH / srcH

	var newW, newH float64

	switch fit {
	case ImageFitContain:
		// 取较小的缩放比例，确保图片完全包含在目标区域内
		ratio := min(ratioW, ratioH)
		newW = srcW * ratio
		newH = srcH * ratio

	case ImageFitCover:
		// 取较大的缩放比例，确保图片完全覆盖目标区域
		ratio := max(ratioW, ratioH)
		newW = srcW * ratio
		newH = srcH * ratio

	case ImageFitFill:
		// 直接使用目标尺寸，不保持宽高比
		newW = targetW
		newH = targetH

	case ImageFitInside:
		// 类似 contain，但缩放比例不超过 1（不会放大图片）
		ratio := min(ratioW, ratioH)
		if ratio > 1 {
			ratio = 1
		}
		newW = srcW * ratio
		newH = srcH * ratio

	case ImageFitOutside:
		// 类似 cover，但缩放比例不小于 1（不会缩小图片）
		ratio := max(ratioW, ratioH)
		if ratio < 1 {
			ratio = 1
		}
		newW = srcW * ratio
		newH = srcH * ratio

	default:
		// 未知模式时默认使用 contain 行为
		ratio := min(ratioW, ratioH)
		newW = srcW * ratio
		newH = srcH * ratio
	}

	// 将浮点数转换为整数，并确保最小值为 1
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
