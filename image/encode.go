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

// EncodeOptions 编码选项，用于控制图片编码时的参数
type EncodeOptions struct {
	Quality int // JPEG 编码质量，取值范围 1-100，值越大质量越高、文件越大
}

// EncodeOption 编码选项的函数式参数类型，采用 functional options 模式，
// 使用者可以通过链式调用灵活组合编码参数
type EncodeOption func(*EncodeOptions)

// Quality 创建一个设置 JPEG 编码质量的选项。
// 参数 q 的取值范围为 1-100，默认值为 85。
// 仅在编码为 JPEG 格式时有效，其他格式会忽略此选项。
func Quality(q int) EncodeOption {
	return func(opts *EncodeOptions) {
		opts.Quality = q
	}
}

// defaultEncodeOptions 返回默认的编码选项，JPEG 质量默认为 85
func defaultEncodeOptions() *EncodeOptions {
	return &EncodeOptions{Quality: 85}
}

// Encode 将图片编码为指定格式并写入 io.Writer。
// 支持的格式包括：jpeg/jpg、png、gif、bmp、tiff。
// 可以通过 opts 参数指定编码选项（如 JPEG 质量）。
// WebP 格式仅支持解码不支持编码，传入 webp 会返回错误。
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
		return fmt.Errorf("不支持的编码格式: %s", format)
	}
}

// Save 将图片保存到指定路径的文件中。
// 编码格式根据文件扩展名自动推断（如 .jpg → JPEG，.png → PNG）。
// 支持的扩展名：.jpg/.jpeg/.png/.gif/.bmp/.tiff/.tif
// 可以通过 opts 参数指定编码选项（如 Quality(90) 设置 JPEG 质量）。
func (i *Image) Save(path string, opts ...EncodeOption) error {
	format := formatFromPath(path)
	if format == "" {
		return fmt.Errorf("无法从文件扩展名推断格式: %s", path)
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return i.Encode(f, format, opts...)
}

// ToBytes 将图片编码为指定格式的字节切片。
// 适用于需要将图片数据存储到内存、数据库或通过网络传输的场景。
// 内部使用 bytes.Buffer 缓存编码数据后返回字节切片。
func (i *Image) ToBytes(format string, opts ...EncodeOption) ([]byte, error) {
	var buf bytes.Buffer
	if err := i.Encode(&buf, format, opts...); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// formatFromPath 根据文件路径的扩展名推断图片格式。
// 返回对应的格式常量字符串，如果扩展名无法识别则返回空字符串。
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
