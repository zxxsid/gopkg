package image

import (
	"image"
	"image/draw"
	_ "image/gif"  // 注册 GIF 格式的解码器，使 image.Decode 能识别 GIF 格式
	_ "image/jpeg" // 注册 JPEG 格式的解码器，使 image.Decode 能识别 JPEG 格式
	_ "image/png"  // 注册 PNG 格式的解码器，使 image.Decode 能识别 PNG 格式
	"io"
	"os"

	_ "golang.org/x/image/bmp"  // 注册 BMP 格式的解码器，使 image.Decode 能识别 BMP 格式
	_ "golang.org/x/image/tiff" // 注册 TIFF 格式的解码器，使 image.Decode 能识别 TIFF 格式
	_ "golang.org/x/image/webp" // 注册 WebP 格式的解码器，使 image.Decode 能识别 WebP 格式
)

// 图片适配模式常量，用于控制缩放时图片与目标尺寸之间的适配方式
const (
	ImageFitContain = "contain" // 保持宽高比缩放，使图片完全包含在目标尺寸内（可能有留白）
	ImageFitCover   = "cover"   // 保持宽高比缩放，使图片完全覆盖目标尺寸（可能有裁切）
	ImageFitFill    = "fill"    // 不保持宽高比，拉伸图片以完全填充目标尺寸
	ImageFitInside  = "inside"  // 类似 contain，但不会放大比目标尺寸小的图片
	ImageFitOutside = "outside" // 类似 cover，但不会缩小比目标尺寸大的图片
)

// 支持的图片格式常量，用于编码和解码时指定格式
const (
	FormatJPEG = "jpeg" // JPEG 格式，有损压缩，适合照片类图片
	FormatPNG  = "png"  // PNG 格式，无损压缩，支持透明通道
	FormatGIF  = "gif"  // GIF 格式，支持动画和透明，最多 256 色
	FormatBMP  = "bmp"  // BMP 格式，无压缩的位图格式
	FormatTIFF = "tiff" // TIFF 格式，支持无损压缩，常用于印刷和扫描
	FormatWEBP = "webp" // WebP 格式，仅支持解码，不支持编码
)

// Image 是图片处理的核心结构体，封装了标准库的 image.Image，
// 提供格式转换、缩放、裁剪、压缩、旋转等一系列图片处理方法。
// 所有变换操作都返回新的 Image 实例，不会修改原始图片（不可变设计）。
type Image struct {
	img    image.Image // 底层的标准库图片对象
	format string      // 图片的原始格式名称（如 "png"、"jpeg"）
}

// NewImage 从 io.Reader 中读取并解码图片数据。
// 支持 PNG、JPEG、GIF、BMP、TIFF、WebP 格式（通过 init 注册的解码器自动识别）。
// 返回解码后的 Image 实例和可能的错误。
func NewImage(r io.Reader) (*Image, error) {
	img, format, err := image.Decode(r)
	if err != nil {
		return nil, err
	}
	return &Image{img: img, format: format}, nil
}

// Open 根据文件路径打开并解码图片文件。
// 内部调用 NewImage 进行解码，格式由文件内容自动识别（而非文件扩展名）。
func Open(path string) (*Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return NewImage(f)
}

// FromImage 从标准库的 image.Image 对象创建 Image 实例。
// format 参数用于指定图片格式名称，后续编码时会用到。
// 适用于从其他来源获取的 image.Image 对象需要进行二次处理的场景。
func FromImage(img image.Image, format string) *Image {
	return &Image{img: img, format: format}
}

// Format 返回图片的格式名称（如 "png"、"jpeg"、"gif" 等）。
// 该值在解码时由 image.Decode 自动识别并设置。
func (i *Image) Format() string {
	return i.format
}

// Width 返回图片的宽度（单位：像素）
func (i *Image) Width() int {
	return i.img.Bounds().Dx()
}

// Height 返回图片的高度（单位：像素）
func (i *Image) Height() int {
	return i.img.Bounds().Dy()
}

// Bounds 返回图片的边界矩形，包含图片的起始坐标和结束坐标。
// 大多数情况下起始坐标为 (0,0)，但某些图片格式可能有非零的起始坐标。
func (i *Image) Bounds() image.Rectangle {
	return i.img.Bounds()
}

// RawImage 返回底层的标准库 image.Image 对象。
// 用于需要直接操作标准库图片接口的场景，如与第三方库的互操作。
func (i *Image) RawImage() image.Image {
	return i.img
}

// Clone 创建当前图片的深拷贝副本。
// 返回一个新的 Image 实例，其像素数据独立于原始图片，
// 修改副本不会影响原始图片。内部使用 draw.Draw 进行高效的像素复制。
func (i *Image) Clone() *Image {
	bounds := i.img.Bounds()
	dst := image.NewRGBA(bounds)
	draw.Draw(dst, bounds, i.img, bounds.Min, draw.Src)
	return &Image{img: dst, format: i.format}
}
