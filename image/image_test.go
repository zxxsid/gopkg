package image

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

// createTestImage 创建一个指定尺寸的测试图片，使用渐变色填充像素。
// RGB 值根据像素坐标计算：R=x%256, G=y%256, B=(x+y)%256，
// 这样可以在测试中通过像素颜色值验证裁剪、翻转等操作的正确性。
func createTestImage(width, height int) *Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{
				R: uint8(x % 256),
				G: uint8(y % 256),
				B: uint8((x + y) % 256),
				A: 255,
			})
		}
	}
	return FromImage(img, FormatPNG)
}

// encodeToPNG 将图片编码为 PNG 字节切片的测试辅助函数，编码失败时直接终止测试
func encodeToPNG(t *testing.T, img *Image) []byte {
	t.Helper()
	data, err := img.ToBytes(FormatPNG)
	if err != nil {
		t.Fatalf("编码为 PNG 失败: %v", err)
	}
	return data
}

// --- 构造函数测试 ---

// TestNewImage 测试从 io.Reader 解码图片，验证尺寸和格式是否正确
func TestNewImage(t *testing.T) {
	src := createTestImage(100, 80)
	data := encodeToPNG(t, src)

	got, err := NewImage(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("NewImage 失败: %v", err)
	}
	if got.Width() != 100 || got.Height() != 80 {
		t.Errorf("尺寸 = %dx%d, 期望 100x80", got.Width(), got.Height())
	}
	if got.Format() != FormatPNG {
		t.Errorf("格式 = %q, 期望 %q", got.Format(), FormatPNG)
	}
}

// TestOpen 测试从文件路径打开图片，验证尺寸是否正确
func TestOpen(t *testing.T) {
	src := createTestImage(60, 40)
	tmp := filepath.Join(t.TempDir(), "test.png")
	if err := src.Save(tmp); err != nil {
		t.Fatalf("保存失败: %v", err)
	}

	got, err := Open(tmp)
	if err != nil {
		t.Fatalf("Open 失败: %v", err)
	}
	if got.Width() != 60 || got.Height() != 40 {
		t.Errorf("尺寸 = %dx%d, 期望 60x40", got.Width(), got.Height())
	}
}

// TestFromImage 测试从标准库 image.Image 创建 Image 实例，
// 验证尺寸正确且 RawImage 返回的是同一个底层对象
func TestFromImage(t *testing.T) {
	raw := image.NewRGBA(image.Rect(0, 0, 10, 10))
	img := FromImage(raw, FormatPNG)
	if img.Width() != 10 || img.Height() != 10 {
		t.Errorf("尺寸 = %dx%d, 期望 10x10", img.Width(), img.Height())
	}
	if img.RawImage() != raw {
		t.Error("RawImage 应该返回原始的 image.Image 对象")
	}
}

// TestClone 测试深拷贝功能，验证副本与原图尺寸相同但底层对象不同
func TestClone(t *testing.T) {
	src := createTestImage(20, 20)
	cloned := src.Clone()
	if cloned.Width() != src.Width() || cloned.Height() != src.Height() {
		t.Errorf("副本尺寸与原图不匹配")
	}
	if cloned.RawImage() == src.RawImage() {
		t.Error("Clone 应该返回不同的底层 image.Image 对象")
	}
}

// --- 格式转换测试 ---

// TestConvertPNGToJPEG 测试 PNG 转 JPEG 格式转换，验证编码解码后格式正确
func TestConvertPNGToJPEG(t *testing.T) {
	src := createTestImage(100, 100)
	jpegData, err := src.ToBytes(FormatJPEG, Quality(90))
	if err != nil {
		t.Fatalf("编码为 JPEG 失败: %v", err)
	}

	got, err := NewImage(bytes.NewReader(jpegData))
	if err != nil {
		t.Fatalf("解码 JPEG 失败: %v", err)
	}
	if got.Format() != FormatJPEG {
		t.Errorf("格式 = %q, 期望 %q", got.Format(), FormatJPEG)
	}
}

// TestConvertPNGToGIF 测试 PNG 转 GIF 格式转换
func TestConvertPNGToGIF(t *testing.T) {
	src := createTestImage(50, 50)
	gifData, err := src.ToBytes(FormatGIF)
	if err != nil {
		t.Fatalf("编码为 GIF 失败: %v", err)
	}

	got, err := NewImage(bytes.NewReader(gifData))
	if err != nil {
		t.Fatalf("解码 GIF 失败: %v", err)
	}
	if got.Format() != FormatGIF {
		t.Errorf("格式 = %q, 期望 %q", got.Format(), FormatGIF)
	}
}

// TestConvertPNGToBMP 测试 PNG 转 BMP 格式转换
func TestConvertPNGToBMP(t *testing.T) {
	src := createTestImage(50, 50)
	bmpData, err := src.ToBytes(FormatBMP)
	if err != nil {
		t.Fatalf("编码为 BMP 失败: %v", err)
	}

	got, err := NewImage(bytes.NewReader(bmpData))
	if err != nil {
		t.Fatalf("解码 BMP 失败: %v", err)
	}
	if got.Format() != FormatBMP {
		t.Errorf("格式 = %q, 期望 %q", got.Format(), FormatBMP)
	}
}

// TestConvertPNGToTIFF 测试 PNG 转 TIFF 格式转换
func TestConvertPNGToTIFF(t *testing.T) {
	src := createTestImage(50, 50)
	tiffData, err := src.ToBytes(FormatTIFF)
	if err != nil {
		t.Fatalf("编码为 TIFF 失败: %v", err)
	}

	got, err := NewImage(bytes.NewReader(tiffData))
	if err != nil {
		t.Fatalf("解码 TIFF 失败: %v", err)
	}
	if got.Format() != FormatTIFF {
		t.Errorf("格式 = %q, 期望 %q", got.Format(), FormatTIFF)
	}
}

// TestEncodeUnsupportedFormat 测试不支持的编码格式（如 webp）应返回错误
func TestEncodeUnsupportedFormat(t *testing.T) {
	src := createTestImage(10, 10)
	_, err := src.ToBytes("webp")
	if err == nil {
		t.Error("对不支持的编码格式应返回错误")
	}
}

// TestSaveWithExtension 测试通过不同扩展名保存文件，验证文件是否成功创建
func TestSaveWithExtension(t *testing.T) {
	src := createTestImage(30, 30)
	dir := t.TempDir()

	for _, ext := range []string{".png", ".jpg", ".gif", ".bmp", ".tiff"} {
		path := filepath.Join(dir, "test"+ext)
		if err := src.Save(path); err != nil {
			t.Errorf("保存 %s 失败: %v", ext, err)
		}
		if _, err := os.Stat(path); err != nil {
			t.Errorf("文件 %s 未创建", path)
		}
	}
}

// --- 缩放测试 ---

// TestResize 测试指定宽高的缩放，验证缩放后的尺寸是否正确
func TestResize(t *testing.T) {
	src := createTestImage(200, 100)
	resized := src.Resize(100, 50)
	if resized.Width() != 100 || resized.Height() != 50 {
		t.Errorf("尺寸 = %dx%d, 期望 100x50", resized.Width(), resized.Height())
	}
}

// TestResizeByWidth 测试按宽度等比缩放，验证高度是否按比例自动计算
func TestResizeByWidth(t *testing.T) {
	src := createTestImage(200, 100)
	resized := src.ResizeByWidth(100)
	if resized.Width() != 100 {
		t.Errorf("宽度 = %d, 期望 100", resized.Width())
	}
	if resized.Height() != 50 {
		t.Errorf("高度 = %d, 期望 50（按宽高比 2:1 自动计算）", resized.Height())
	}
}

// TestResizeByHeight 测试按高度等比缩放，验证宽度是否按比例自动计算
func TestResizeByHeight(t *testing.T) {
	src := createTestImage(200, 100)
	resized := src.ResizeByHeight(50)
	if resized.Height() != 50 {
		t.Errorf("高度 = %d, 期望 50", resized.Height())
	}
	if resized.Width() != 100 {
		t.Errorf("宽度 = %d, 期望 100（按宽高比 2:1 自动计算）", resized.Width())
	}
}

// TestThumbnail 测试缩略图生成，验证结果不超过目标尺寸且至少一边贴合目标
func TestThumbnail(t *testing.T) {
	src := createTestImage(800, 600)
	thumb := src.Thumbnail(200, 200)
	if thumb.Width() > 200 || thumb.Height() > 200 {
		t.Errorf("缩略图 %dx%d 超出 200x200 的限制", thumb.Width(), thumb.Height())
	}
	if thumb.Width() != 200 && thumb.Height() != 200 {
		t.Errorf("缩略图应至少有一边等于目标尺寸")
	}
}

// TestThumbnailNoUpscale 测试缩略图不放大小图：
// 当原图小于目标尺寸时，应保持原始尺寸不变
func TestThumbnailNoUpscale(t *testing.T) {
	src := createTestImage(100, 80)
	thumb := src.Thumbnail(200, 200)
	if thumb.Width() != 100 || thumb.Height() != 80 {
		t.Errorf("缩略图不应放大小图: 得到 %dx%d", thumb.Width(), thumb.Height())
	}
}

// TestResizeWithFitContain 测试 contain 模式：
// 400x200 的图片适配到 100x100，应得到 100x50（取较小缩放比 0.25）
func TestResizeWithFitContain(t *testing.T) {
	src := createTestImage(400, 200)
	resized := src.ResizeWithFit(100, 100, ImageFitContain)
	if resized.Width() != 100 || resized.Height() != 50 {
		t.Errorf("contain 模式: %dx%d, 期望 100x50", resized.Width(), resized.Height())
	}
}

// TestResizeWithFitCover 测试 cover 模式：
// 400x200 的图片适配到 100x100，应得到 200x100（取较大缩放比 0.5）
func TestResizeWithFitCover(t *testing.T) {
	src := createTestImage(400, 200)
	resized := src.ResizeWithFit(100, 100, ImageFitCover)
	if resized.Width() != 200 || resized.Height() != 100 {
		t.Errorf("cover 模式: %dx%d, 期望 200x100", resized.Width(), resized.Height())
	}
}

// TestResizeWithFitFill 测试 fill 模式：
// 不保持宽高比，强制拉伸到目标尺寸 100x100
func TestResizeWithFitFill(t *testing.T) {
	src := createTestImage(400, 200)
	resized := src.ResizeWithFit(100, 100, ImageFitFill)
	if resized.Width() != 100 || resized.Height() != 100 {
		t.Errorf("fill 模式: %dx%d, 期望 100x100", resized.Width(), resized.Height())
	}
}

// --- 裁剪测试 ---

// TestCrop 测试从指定坐标 (50,50) 裁剪 100x80 的区域
func TestCrop(t *testing.T) {
	src := createTestImage(200, 200)
	cropped := src.Crop(50, 50, 100, 80)
	if cropped.Width() != 100 || cropped.Height() != 80 {
		t.Errorf("尺寸 = %dx%d, 期望 100x80", cropped.Width(), cropped.Height())
	}
}

// TestCropCenter 测试从中心裁剪 100x100 的区域
func TestCropCenter(t *testing.T) {
	src := createTestImage(200, 200)
	cropped := src.CropCenter(100, 100)
	if cropped.Width() != 100 || cropped.Height() != 100 {
		t.Errorf("尺寸 = %dx%d, 期望 100x100", cropped.Width(), cropped.Height())
	}
}

// TestCropAnchorTopLeft 测试从左上角锚点裁剪 100x100 的区域
func TestCropAnchorTopLeft(t *testing.T) {
	src := createTestImage(200, 200)
	cropped := src.CropAnchor(100, 100, AnchorTopLeft)
	if cropped.Width() != 100 || cropped.Height() != 100 {
		t.Errorf("尺寸 = %dx%d, 期望 100x100", cropped.Width(), cropped.Height())
	}
}

// TestCropAnchorBottomRight 测试从右下角锚点裁剪 100x100 的区域
func TestCropAnchorBottomRight(t *testing.T) {
	src := createTestImage(200, 200)
	cropped := src.CropAnchor(100, 100, AnchorBottomRight)
	if cropped.Width() != 100 || cropped.Height() != 100 {
		t.Errorf("尺寸 = %dx%d, 期望 100x100", cropped.Width(), cropped.Height())
	}
}

// TestCropOutOfBounds 测试裁剪区域超出图片边界时的自动裁切行为：
// 从 (50,50) 裁剪 200x200，实际可用区域只有 50x50
func TestCropOutOfBounds(t *testing.T) {
	src := createTestImage(100, 100)
	cropped := src.Crop(50, 50, 200, 200)
	if cropped.Width() != 50 || cropped.Height() != 50 {
		t.Errorf("应自动裁切到图片范围内, 得到 %dx%d", cropped.Width(), cropped.Height())
	}
}

// --- 压缩测试 ---

// TestCompress 测试将 500x500 图片压缩到 50KB 以内
func TestCompress(t *testing.T) {
	src := createTestImage(500, 500)
	maxBytes := int64(50_000)
	data, err := src.Compress(maxBytes)
	if err != nil {
		t.Fatalf("Compress 失败: %v", err)
	}
	if int64(len(data)) > maxBytes {
		t.Errorf("压缩后大小 = %d 字节, 期望 <= %d", len(data), maxBytes)
	}
}

// TestCompressToFile 测试压缩并保存到文件，验证文件大小不超过限制
func TestCompressToFile(t *testing.T) {
	src := createTestImage(300, 300)
	path := filepath.Join(t.TempDir(), "compressed.jpg")
	maxBytes := int64(30_000)
	if err := src.CompressToFile(maxBytes, path); err != nil {
		t.Fatalf("CompressToFile 失败: %v", err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("读取文件信息失败: %v", err)
	}
	if info.Size() > maxBytes {
		t.Errorf("文件大小 = %d 字节, 期望 <= %d", info.Size(), maxBytes)
	}
}

// TestCompressAuto 测试智能压缩：将 1000x1000 的大图压缩到 10KB 以内，
// 需要同时缩小尺寸和降低质量才能达到目标
func TestCompressAuto(t *testing.T) {
	src := createTestImage(1000, 1000)
	maxBytes := int64(10_000)
	data, err := src.CompressAuto(maxBytes)
	if err != nil {
		t.Fatalf("CompressAuto 失败: %v", err)
	}
	if int64(len(data)) > maxBytes {
		t.Errorf("压缩后大小 = %d 字节, 期望 <= %d", len(data), maxBytes)
	}
}

// --- 变换操作测试 ---

// TestRotate90 测试顺时针旋转 90 度后宽高互换：200x100 → 100x200
func TestRotate90(t *testing.T) {
	src := createTestImage(200, 100)
	rotated := src.Rotate90()
	if rotated.Width() != 100 || rotated.Height() != 200 {
		t.Errorf("尺寸 = %dx%d, 期望 100x200", rotated.Width(), rotated.Height())
	}
}

// TestRotate180 测试旋转 180 度后尺寸不变
func TestRotate180(t *testing.T) {
	src := createTestImage(200, 100)
	rotated := src.Rotate180()
	if rotated.Width() != 200 || rotated.Height() != 100 {
		t.Errorf("尺寸 = %dx%d, 期望 200x100", rotated.Width(), rotated.Height())
	}
}

// TestRotate270 测试顺时针旋转 270 度后宽高互换：200x100 → 100x200
func TestRotate270(t *testing.T) {
	src := createTestImage(200, 100)
	rotated := src.Rotate270()
	if rotated.Width() != 100 || rotated.Height() != 200 {
		t.Errorf("尺寸 = %dx%d, 期望 100x200", rotated.Width(), rotated.Height())
	}
}

// TestFlipH 测试水平翻转：左上角像素应与右上角像素对应
func TestFlipH(t *testing.T) {
	src := createTestImage(100, 100)
	flipped := src.FlipH()
	if flipped.Width() != 100 || flipped.Height() != 100 {
		t.Errorf("翻转后尺寸不应改变")
	}

	// 验证水平翻转的像素对应关系：(0,0) 的 R 值应等于翻转后 (99,0) 的 R 值
	srcR, _, _, _ := src.RawImage().At(0, 0).RGBA()
	flipR, _, _, _ := flipped.RawImage().At(99, 0).RGBA()
	if srcR != flipR {
		t.Errorf("原图 (0,0) 的 R 值应等于翻转后 (99,0) 的 R 值")
	}
}

// TestFlipV 测试垂直翻转：左上角像素应与左下角像素对应
func TestFlipV(t *testing.T) {
	src := createTestImage(100, 100)
	flipped := src.FlipV()
	if flipped.Width() != 100 || flipped.Height() != 100 {
		t.Errorf("翻转后尺寸不应改变")
	}

	// 验证垂直翻转的像素对应关系：(0,0) 的 R 值应等于翻转后 (0,99) 的 R 值
	srcR, _, _, _ := src.RawImage().At(0, 0).RGBA()
	flipR, _, _, _ := flipped.RawImage().At(0, 99).RGBA()
	if srcR != flipR {
		t.Errorf("原图 (0,0) 的 R 值应等于翻转后 (0,99) 的 R 值")
	}
}

// TestGrayscale 测试灰度转换：转换后每个像素的 R、G、B 值应相同
func TestGrayscale(t *testing.T) {
	src := createTestImage(100, 100)
	gray := src.Grayscale()
	if gray.Width() != 100 || gray.Height() != 100 {
		t.Errorf("灰度转换后尺寸不应改变")
	}

	r, g, b, _ := gray.RawImage().At(50, 50).RGBA()
	if r != g || g != b {
		t.Errorf("灰度像素的 RGB 值应相同, 得到 r=%d g=%d b=%d", r, g, b)
	}
}

// TestOpacity 测试不透明度设置：设为 0.5 后 Alpha 值应约为 127
func TestOpacity(t *testing.T) {
	src := createTestImage(100, 100)
	halfAlpha := src.Opacity(0.5)
	_, _, _, a := halfAlpha.RawImage().At(50, 50).RGBA()
	expected := uint32(127) << 8
	diff := int64(a) - int64(expected)
	if diff < -512 || diff > 512 {
		t.Errorf("Alpha = %d, 期望约 %d（允许 ±512 的误差）", a, expected)
	}
}

// TestWatermark 测试水印叠加：叠加后尺寸应与底图一致
func TestWatermark(t *testing.T) {
	base := createTestImage(200, 200)
	mark := createTestImage(50, 50)
	result := base.Watermark(mark, 10, 10)
	if result.Width() != 200 || result.Height() != 200 {
		t.Errorf("水印叠加后尺寸应与底图一致")
	}
}

// --- 集成测试：完整处理流水线 ---

// TestFullPipeline 测试链式调用多个处理操作的完整流水线：
// 缩放 → 中心裁剪 → 灰度转换 → 旋转 → 保存为 JPEG → 重新加载验证
func TestFullPipeline(t *testing.T) {
	src := createTestImage(800, 600)

	// 链式处理：等比缩放到宽400 → 中心裁剪300x200 → 灰度化 → 顺时针旋转90度
	result := src.
		Resize(400, 0).
		CropCenter(300, 200).
		Grayscale().
		Rotate90()

	// 旋转90度后宽高互换：300x200 → 200x300
	if result.Width() != 200 || result.Height() != 300 {
		t.Errorf("流水线处理后尺寸 = %dx%d, 期望 200x300", result.Width(), result.Height())
	}

	// 保存为 JPEG 并重新加载验证格式
	dir := t.TempDir()
	jpegPath := filepath.Join(dir, "pipeline.jpg")
	if err := result.Save(jpegPath, Quality(75)); err != nil {
		t.Fatalf("保存失败: %v", err)
	}

	reloaded, err := Open(jpegPath)
	if err != nil {
		t.Fatalf("重新加载失败: %v", err)
	}
	if reloaded.Format() != FormatJPEG {
		t.Errorf("格式 = %q, 期望 jpeg", reloaded.Format())
	}
}

// --- 往返测试：保存并重新加载后尺寸应保持一致 ---

// TestSaveReloadRoundtrip 测试各种格式的保存-加载往返一致性，
// 验证每种格式保存后重新打开，尺寸和格式都应正确
func TestSaveReloadRoundtrip(t *testing.T) {
	src := createTestImage(150, 100)
	dir := t.TempDir()

	formats := map[string]string{
		".png":  FormatPNG,
		".jpg":  FormatJPEG,
		".gif":  FormatGIF,
		".bmp":  FormatBMP,
		".tiff": FormatTIFF,
	}

	for ext, wantFmt := range formats {
		path := filepath.Join(dir, "roundtrip"+ext)
		if err := src.Save(path); err != nil {
			t.Errorf("Save(%s) 失败: %v", ext, err)
			continue
		}
		got, err := Open(path)
		if err != nil {
			t.Errorf("Open(%s) 失败: %v", ext, err)
			continue
		}
		if got.Format() != wantFmt {
			t.Errorf("%s 格式 = %q, 期望 %q", ext, got.Format(), wantFmt)
		}
		if got.Width() != 150 || got.Height() != 100 {
			t.Errorf("%s 尺寸 = %dx%d, 期望 150x100", ext, got.Width(), got.Height())
		}
	}
}

// --- 性能基准测试 ---

// BenchmarkResize 基准测试：将 1000x1000 图片缩放到 500x500
func BenchmarkResize(b *testing.B) {
	src := createTestImage(1000, 1000)
	b.ResetTimer()
	for b.Loop() {
		src.Resize(500, 500)
	}
}

// BenchmarkCrop 基准测试：从 1000x1000 图片裁剪 500x500 区域
func BenchmarkCrop(b *testing.B) {
	src := createTestImage(1000, 1000)
	b.ResetTimer()
	for b.Loop() {
		src.Crop(100, 100, 500, 500)
	}
}

// BenchmarkGrayscale 基准测试：将 1000x1000 图片转为灰度
func BenchmarkGrayscale(b *testing.B) {
	src := createTestImage(1000, 1000)
	b.ResetTimer()
	for b.Loop() {
		src.Grayscale()
	}
}

// BenchmarkEncodePNG 基准测试：将 500x500 图片编码为 PNG
func BenchmarkEncodePNG(b *testing.B) {
	src := createTestImage(500, 500)
	b.ResetTimer()
	for b.Loop() {
		src.ToBytes(FormatPNG)
	}
}

// BenchmarkEncodeJPEG 基准测试：将 500x500 图片编码为 JPEG（质量 85）
func BenchmarkEncodeJPEG(b *testing.B) {
	src := createTestImage(500, 500)
	b.ResetTimer()
	for b.Loop() {
		src.ToBytes(FormatJPEG, Quality(85))
	}
}

// BenchmarkCompress 基准测试：将 500x500 图片压缩到 50KB 以内
func BenchmarkCompress(b *testing.B) {
	src := createTestImage(500, 500)
	b.ResetTimer()
	for b.Loop() {
		src.Compress(50_000)
	}
}

// --- 从原始 PNG 字节解码测试 ---

// TestNewImageFromPNGBytes 测试从手动构建的 PNG 字节数据解码图片，
// 验证 NewImage 能正确处理标准的 PNG 编码数据
func TestNewImageFromPNGBytes(t *testing.T) {
	raw := image.NewRGBA(image.Rect(0, 0, 4, 4))
	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			raw.Set(x, y, color.RGBA{R: 255, A: 255})
		}
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, raw); err != nil {
		t.Fatal(err)
	}

	img, err := NewImage(&buf)
	if err != nil {
		t.Fatalf("NewImage 失败: %v", err)
	}
	if img.Width() != 4 || img.Height() != 4 {
		t.Errorf("尺寸 = %dx%d, 期望 4x4", img.Width(), img.Height())
	}
	if img.Format() != FormatPNG {
		t.Errorf("格式 = %q, 期望 %q", img.Format(), FormatPNG)
	}
}
