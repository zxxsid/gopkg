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

// createTestImage builds an RGBA image with a gradient pattern
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

func encodeToPNG(t *testing.T, img *Image) []byte {
	t.Helper()
	data, err := img.ToBytes(FormatPNG)
	if err != nil {
		t.Fatalf("encode to PNG: %v", err)
	}
	return data
}

// --- Constructor tests ---

func TestNewImage(t *testing.T) {
	src := createTestImage(100, 80)
	data := encodeToPNG(t, src)

	got, err := NewImage(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("NewImage: %v", err)
	}
	if got.Width() != 100 || got.Height() != 80 {
		t.Errorf("dimensions = %dx%d, want 100x80", got.Width(), got.Height())
	}
	if got.Format() != FormatPNG {
		t.Errorf("format = %q, want %q", got.Format(), FormatPNG)
	}
}

func TestOpen(t *testing.T) {
	src := createTestImage(60, 40)
	tmp := filepath.Join(t.TempDir(), "test.png")
	if err := src.Save(tmp); err != nil {
		t.Fatalf("save: %v", err)
	}

	got, err := Open(tmp)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	if got.Width() != 60 || got.Height() != 40 {
		t.Errorf("dimensions = %dx%d, want 60x40", got.Width(), got.Height())
	}
}

func TestFromImage(t *testing.T) {
	raw := image.NewRGBA(image.Rect(0, 0, 10, 10))
	img := FromImage(raw, FormatPNG)
	if img.Width() != 10 || img.Height() != 10 {
		t.Errorf("dimensions = %dx%d, want 10x10", img.Width(), img.Height())
	}
	if img.RawImage() != raw {
		t.Error("RawImage should return the original image")
	}
}

func TestClone(t *testing.T) {
	src := createTestImage(20, 20)
	cloned := src.Clone()
	if cloned.Width() != src.Width() || cloned.Height() != src.Height() {
		t.Errorf("cloned dimensions mismatch")
	}
	if cloned.RawImage() == src.RawImage() {
		t.Error("clone should return a different underlying image")
	}
}

// --- Format conversion tests ---

func TestConvertPNGToJPEG(t *testing.T) {
	src := createTestImage(100, 100)
	jpegData, err := src.ToBytes(FormatJPEG, Quality(90))
	if err != nil {
		t.Fatalf("encode JPEG: %v", err)
	}

	got, err := NewImage(bytes.NewReader(jpegData))
	if err != nil {
		t.Fatalf("decode JPEG: %v", err)
	}
	if got.Format() != FormatJPEG {
		t.Errorf("format = %q, want %q", got.Format(), FormatJPEG)
	}
}

func TestConvertPNGToGIF(t *testing.T) {
	src := createTestImage(50, 50)
	gifData, err := src.ToBytes(FormatGIF)
	if err != nil {
		t.Fatalf("encode GIF: %v", err)
	}

	got, err := NewImage(bytes.NewReader(gifData))
	if err != nil {
		t.Fatalf("decode GIF: %v", err)
	}
	if got.Format() != FormatGIF {
		t.Errorf("format = %q, want %q", got.Format(), FormatGIF)
	}
}

func TestConvertPNGToBMP(t *testing.T) {
	src := createTestImage(50, 50)
	bmpData, err := src.ToBytes(FormatBMP)
	if err != nil {
		t.Fatalf("encode BMP: %v", err)
	}

	got, err := NewImage(bytes.NewReader(bmpData))
	if err != nil {
		t.Fatalf("decode BMP: %v", err)
	}
	if got.Format() != FormatBMP {
		t.Errorf("format = %q, want %q", got.Format(), FormatBMP)
	}
}

func TestConvertPNGToTIFF(t *testing.T) {
	src := createTestImage(50, 50)
	tiffData, err := src.ToBytes(FormatTIFF)
	if err != nil {
		t.Fatalf("encode TIFF: %v", err)
	}

	got, err := NewImage(bytes.NewReader(tiffData))
	if err != nil {
		t.Fatalf("decode TIFF: %v", err)
	}
	if got.Format() != FormatTIFF {
		t.Errorf("format = %q, want %q", got.Format(), FormatTIFF)
	}
}

func TestEncodeUnsupportedFormat(t *testing.T) {
	src := createTestImage(10, 10)
	_, err := src.ToBytes("webp")
	if err == nil {
		t.Error("expected error for unsupported encoding format")
	}
}

func TestSaveWithExtension(t *testing.T) {
	src := createTestImage(30, 30)
	dir := t.TempDir()

	for _, ext := range []string{".png", ".jpg", ".gif", ".bmp", ".tiff"} {
		path := filepath.Join(dir, "test"+ext)
		if err := src.Save(path); err != nil {
			t.Errorf("Save %s: %v", ext, err)
		}
		if _, err := os.Stat(path); err != nil {
			t.Errorf("file %s not created", path)
		}
	}
}

// --- Resize tests ---

func TestResize(t *testing.T) {
	src := createTestImage(200, 100)
	resized := src.Resize(100, 50)
	if resized.Width() != 100 || resized.Height() != 50 {
		t.Errorf("dimensions = %dx%d, want 100x50", resized.Width(), resized.Height())
	}
}

func TestResizeByWidth(t *testing.T) {
	src := createTestImage(200, 100)
	resized := src.ResizeByWidth(100)
	if resized.Width() != 100 {
		t.Errorf("width = %d, want 100", resized.Width())
	}
	if resized.Height() != 50 {
		t.Errorf("height = %d, want 50", resized.Height())
	}
}

func TestResizeByHeight(t *testing.T) {
	src := createTestImage(200, 100)
	resized := src.ResizeByHeight(50)
	if resized.Height() != 50 {
		t.Errorf("height = %d, want 50", resized.Height())
	}
	if resized.Width() != 100 {
		t.Errorf("width = %d, want 100", resized.Width())
	}
}

func TestThumbnail(t *testing.T) {
	src := createTestImage(800, 600)
	thumb := src.Thumbnail(200, 200)
	if thumb.Width() > 200 || thumb.Height() > 200 {
		t.Errorf("thumbnail %dx%d exceeds 200x200", thumb.Width(), thumb.Height())
	}
	if thumb.Width() != 200 && thumb.Height() != 200 {
		t.Errorf("thumbnail should touch at least one bound")
	}
}

func TestThumbnailNoUpscale(t *testing.T) {
	src := createTestImage(100, 80)
	thumb := src.Thumbnail(200, 200)
	if thumb.Width() != 100 || thumb.Height() != 80 {
		t.Errorf("thumbnail should not upscale: got %dx%d", thumb.Width(), thumb.Height())
	}
}

func TestResizeWithFitContain(t *testing.T) {
	src := createTestImage(400, 200)
	resized := src.ResizeWithFit(100, 100, ImageFitContain)
	if resized.Width() != 100 || resized.Height() != 50 {
		t.Errorf("contain: %dx%d, want 100x50", resized.Width(), resized.Height())
	}
}

func TestResizeWithFitCover(t *testing.T) {
	src := createTestImage(400, 200)
	resized := src.ResizeWithFit(100, 100, ImageFitCover)
	if resized.Width() != 200 || resized.Height() != 100 {
		t.Errorf("cover: %dx%d, want 200x100", resized.Width(), resized.Height())
	}
}

func TestResizeWithFitFill(t *testing.T) {
	src := createTestImage(400, 200)
	resized := src.ResizeWithFit(100, 100, ImageFitFill)
	if resized.Width() != 100 || resized.Height() != 100 {
		t.Errorf("fill: %dx%d, want 100x100", resized.Width(), resized.Height())
	}
}

// --- Crop tests ---

func TestCrop(t *testing.T) {
	src := createTestImage(200, 200)
	cropped := src.Crop(50, 50, 100, 80)
	if cropped.Width() != 100 || cropped.Height() != 80 {
		t.Errorf("dimensions = %dx%d, want 100x80", cropped.Width(), cropped.Height())
	}
}

func TestCropCenter(t *testing.T) {
	src := createTestImage(200, 200)
	cropped := src.CropCenter(100, 100)
	if cropped.Width() != 100 || cropped.Height() != 100 {
		t.Errorf("dimensions = %dx%d, want 100x100", cropped.Width(), cropped.Height())
	}
}

func TestCropAnchorTopLeft(t *testing.T) {
	src := createTestImage(200, 200)
	cropped := src.CropAnchor(100, 100, AnchorTopLeft)
	if cropped.Width() != 100 || cropped.Height() != 100 {
		t.Errorf("dimensions = %dx%d, want 100x100", cropped.Width(), cropped.Height())
	}
}

func TestCropAnchorBottomRight(t *testing.T) {
	src := createTestImage(200, 200)
	cropped := src.CropAnchor(100, 100, AnchorBottomRight)
	if cropped.Width() != 100 || cropped.Height() != 100 {
		t.Errorf("dimensions = %dx%d, want 100x100", cropped.Width(), cropped.Height())
	}
}

func TestCropOutOfBounds(t *testing.T) {
	src := createTestImage(100, 100)
	cropped := src.Crop(50, 50, 200, 200)
	if cropped.Width() != 50 || cropped.Height() != 50 {
		t.Errorf("should clamp to image bounds, got %dx%d", cropped.Width(), cropped.Height())
	}
}

// --- Compress tests ---

func TestCompress(t *testing.T) {
	src := createTestImage(500, 500)
	maxBytes := int64(50_000)
	data, err := src.Compress(maxBytes)
	if err != nil {
		t.Fatalf("Compress: %v", err)
	}
	if int64(len(data)) > maxBytes {
		t.Errorf("size = %d, want <= %d", len(data), maxBytes)
	}
}

func TestCompressToFile(t *testing.T) {
	src := createTestImage(300, 300)
	path := filepath.Join(t.TempDir(), "compressed.jpg")
	maxBytes := int64(30_000)
	if err := src.CompressToFile(maxBytes, path); err != nil {
		t.Fatalf("CompressToFile: %v", err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if info.Size() > maxBytes {
		t.Errorf("file size = %d, want <= %d", info.Size(), maxBytes)
	}
}

func TestCompressAuto(t *testing.T) {
	src := createTestImage(1000, 1000)
	maxBytes := int64(10_000)
	data, err := src.CompressAuto(maxBytes)
	if err != nil {
		t.Fatalf("CompressAuto: %v", err)
	}
	if int64(len(data)) > maxBytes {
		t.Errorf("size = %d, want <= %d", len(data), maxBytes)
	}
}

// --- Transform tests ---

func TestRotate90(t *testing.T) {
	src := createTestImage(200, 100)
	rotated := src.Rotate90()
	if rotated.Width() != 100 || rotated.Height() != 200 {
		t.Errorf("dimensions = %dx%d, want 100x200", rotated.Width(), rotated.Height())
	}
}

func TestRotate180(t *testing.T) {
	src := createTestImage(200, 100)
	rotated := src.Rotate180()
	if rotated.Width() != 200 || rotated.Height() != 100 {
		t.Errorf("dimensions = %dx%d, want 200x100", rotated.Width(), rotated.Height())
	}
}

func TestRotate270(t *testing.T) {
	src := createTestImage(200, 100)
	rotated := src.Rotate270()
	if rotated.Width() != 100 || rotated.Height() != 200 {
		t.Errorf("dimensions = %dx%d, want 100x200", rotated.Width(), rotated.Height())
	}
}

func TestFlipH(t *testing.T) {
	src := createTestImage(100, 100)
	flipped := src.FlipH()
	if flipped.Width() != 100 || flipped.Height() != 100 {
		t.Errorf("dimensions should be unchanged")
	}

	srcR, _, _, _ := src.RawImage().At(0, 0).RGBA()
	flipR, _, _, _ := flipped.RawImage().At(99, 0).RGBA()
	if srcR != flipR {
		t.Errorf("pixel (0,0) should match flipped pixel (99,0)")
	}
}

func TestFlipV(t *testing.T) {
	src := createTestImage(100, 100)
	flipped := src.FlipV()
	if flipped.Width() != 100 || flipped.Height() != 100 {
		t.Errorf("dimensions should be unchanged")
	}

	srcR, _, _, _ := src.RawImage().At(0, 0).RGBA()
	flipR, _, _, _ := flipped.RawImage().At(0, 99).RGBA()
	if srcR != flipR {
		t.Errorf("pixel (0,0) should match flipped pixel (0,99)")
	}
}

func TestGrayscale(t *testing.T) {
	src := createTestImage(100, 100)
	gray := src.Grayscale()
	if gray.Width() != 100 || gray.Height() != 100 {
		t.Errorf("dimensions should be unchanged")
	}

	r, g, b, _ := gray.RawImage().At(50, 50).RGBA()
	if r != g || g != b {
		t.Errorf("grayscale pixel should have equal RGB values, got r=%d g=%d b=%d", r, g, b)
	}
}

func TestOpacity(t *testing.T) {
	src := createTestImage(100, 100)
	halfAlpha := src.Opacity(0.5)
	_, _, _, a := halfAlpha.RawImage().At(50, 50).RGBA()
	expected := uint32(127) << 8
	diff := int64(a) - int64(expected)
	if diff < -512 || diff > 512 {
		t.Errorf("alpha = %d, want ~%d", a, expected)
	}
}

func TestWatermark(t *testing.T) {
	base := createTestImage(200, 200)
	mark := createTestImage(50, 50)
	result := base.Watermark(mark, 10, 10)
	if result.Width() != 200 || result.Height() != 200 {
		t.Errorf("dimensions should match base image")
	}
}

// --- Integration: full pipeline ---

func TestFullPipeline(t *testing.T) {
	src := createTestImage(800, 600)

	result := src.
		Resize(400, 0).
		CropCenter(300, 200).
		Grayscale().
		Rotate90()

	if result.Width() != 200 || result.Height() != 300 {
		t.Errorf("pipeline dimensions = %dx%d, want 200x300", result.Width(), result.Height())
	}

	dir := t.TempDir()
	jpegPath := filepath.Join(dir, "pipeline.jpg")
	if err := result.Save(jpegPath, Quality(75)); err != nil {
		t.Fatalf("save: %v", err)
	}

	reloaded, err := Open(jpegPath)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	if reloaded.Format() != FormatJPEG {
		t.Errorf("format = %q, want jpeg", reloaded.Format())
	}
}

// --- Roundtrip: save and reload preserves dimensions ---

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
			t.Errorf("Save(%s): %v", ext, err)
			continue
		}
		got, err := Open(path)
		if err != nil {
			t.Errorf("Open(%s): %v", ext, err)
			continue
		}
		if got.Format() != wantFmt {
			t.Errorf("%s format = %q, want %q", ext, got.Format(), wantFmt)
		}
		if got.Width() != 150 || got.Height() != 100 {
			t.Errorf("%s dimensions = %dx%d, want 150x100", ext, got.Width(), got.Height())
		}
	}
}

// --- Benchmark ---

func BenchmarkResize(b *testing.B) {
	src := createTestImage(1000, 1000)
	b.ResetTimer()
	for b.Loop() {
		src.Resize(500, 500)
	}
}

func BenchmarkCrop(b *testing.B) {
	src := createTestImage(1000, 1000)
	b.ResetTimer()
	for b.Loop() {
		src.Crop(100, 100, 500, 500)
	}
}

func BenchmarkGrayscale(b *testing.B) {
	src := createTestImage(1000, 1000)
	b.ResetTimer()
	for b.Loop() {
		src.Grayscale()
	}
}

func BenchmarkEncodePNG(b *testing.B) {
	src := createTestImage(500, 500)
	b.ResetTimer()
	for b.Loop() {
		src.ToBytes(FormatPNG)
	}
}

func BenchmarkEncodeJPEG(b *testing.B) {
	src := createTestImage(500, 500)
	b.ResetTimer()
	for b.Loop() {
		src.ToBytes(FormatJPEG, Quality(85))
	}
}

func BenchmarkCompress(b *testing.B) {
	src := createTestImage(500, 500)
	b.ResetTimer()
	for b.Loop() {
		src.Compress(50_000)
	}
}

// --- Decode from actual PNG bytes ---

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
		t.Fatalf("NewImage: %v", err)
	}
	if img.Width() != 4 || img.Height() != 4 {
		t.Errorf("dimensions = %dx%d, want 4x4", img.Width(), img.Height())
	}
	if img.Format() != FormatPNG {
		t.Errorf("format = %q, want %q", img.Format(), FormatPNG)
	}
}
