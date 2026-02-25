package image

import (
	"image"
	"image/color"
	"image/draw"
)

// Rotate90 将图片顺时针旋转 90 度。
// 旋转后图片的宽高互换：原始的宽变为新的高，原始的高变为新的宽。
// 通过逐像素映射实现：原坐标 (x, y) 映射到新坐标 (h-1-y, x)。
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

// Rotate180 将图片旋转 180 度（等同于同时水平和垂直翻转）。
// 旋转后图片的宽高保持不变。
// 通过逐像素映射实现：原坐标 (x, y) 映射到新坐标 (w-1-x, h-1-y)。
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

// Rotate270 将图片顺时针旋转 270 度（等同于逆时针旋转 90 度）。
// 旋转后图片的宽高互换：原始的宽变为新的高，原始的高变为新的宽。
// 通过逐像素映射实现：原坐标 (x, y) 映射到新坐标 (y, w-1-x)。
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

// FlipH 水平翻转图片（左右镜像）。
// 翻转后图片尺寸不变，每一行的像素顺序被反转。
// 通过逐像素映射实现：原坐标 (x, y) 映射到新坐标 (w-1-x, y)。
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

// FlipV 垂直翻转图片（上下镜像）。
// 翻转后图片尺寸不变，每一列的像素顺序被反转。
// 通过逐像素映射实现：原坐标 (x, y) 映射到新坐标 (x, h-1-y)。
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

// Grayscale 将图片转换为灰度图。
// 使用标准库的 image.Gray 类型，内部采用 ITU-R BT.601 标准的亮度公式：
// Y = 0.299*R + 0.587*G + 0.114*B
// 转换后的灰度图每个像素的 R、G、B 值相同，等于计算出的亮度值。
func (i *Image) Grayscale() *Image {
	bounds := i.img.Bounds()
	gray := image.NewGray(bounds)
	draw.Draw(gray, bounds, i.img, bounds.Min, draw.Src)
	return &Image{img: gray, format: i.format}
}

// Opacity 设置图片的整体不透明度，返回一个应用了新透明度的副本。
// alpha 参数取值范围 0.0（完全透明）到 1.0（完全不透明）。
// 超出范围的值会被自动裁剪到 [0, 1] 区间。
// 常用于制作半透明水印、淡入淡出效果等场景。
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

	// 遍历每个像素，保留 RGB 值但将 Alpha 通道设置为指定值
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

// Watermark 将另一张图片作为水印叠加到当前图片的指定位置 (x, y)。
// 水印图片按原始尺寸叠加，使用 draw.Over 模式进行混合（支持透明度）。
// 如果水印图片带有 Alpha 通道，叠加时会正确处理半透明效果。
// 可以先对水印调用 Opacity() 设置透明度，再传入此方法实现半透明水印效果。
func (i *Image) Watermark(mark *Image, x, y int) *Image {
	bounds := i.img.Bounds()
	dst := image.NewRGBA(bounds)
	// 先将原图绘制到目标画布
	draw.Draw(dst, bounds, i.img, bounds.Min, draw.Src)

	// 再将水印图片以 Over 模式叠加到指定位置
	markBounds := mark.img.Bounds()
	offset := image.Pt(x, y)
	draw.Draw(dst, markBounds.Add(offset), mark.img, markBounds.Min, draw.Over)

	return &Image{img: dst, format: i.format}
}
