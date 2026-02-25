package image

import (
	"bytes"
	"fmt"
	"image/jpeg"
	"os"
)

// Compress 将图片压缩到指定的文件大小以内（单位：字节）。
// 采用二分查找算法在 JPEG 质量 1-100 之间寻找满足大小限制的最高质量值，
// 以在文件大小和画质之间取得最佳平衡。
// 例如 Compress(200*1024) 将图片压缩到 200KB 以内。
// 如果在当前尺寸下即使质量为 1 也无法满足大小要求，则返回错误。
// 此时可以使用 CompressAuto 方法，它会自动缩小尺寸再尝试压缩。
func (i *Image) Compress(maxBytes int64) ([]byte, error) {
	if maxBytes <= 0 {
		return nil, fmt.Errorf("maxBytes 必须为正数")
	}

	// 使用二分查找寻找满足大小限制的最高 JPEG 质量
	lo, hi := 1, 100
	var best []byte

	for lo <= hi {
		mid := (lo + hi) / 2
		var buf bytes.Buffer
		if err := jpeg.Encode(&buf, i.img, &jpeg.Options{Quality: mid}); err != nil {
			return nil, err
		}
		if int64(buf.Len()) <= maxBytes {
			// 当前质量满足大小要求，记录结果并尝试更高质量
			best = buf.Bytes()
			lo = mid + 1
		} else {
			// 当前质量超出大小限制，尝试更低质量
			hi = mid - 1
		}
	}

	// 如果所有质量等级都不满足要求，报告最小可能的文件大小
	if best == nil {
		var buf bytes.Buffer
		if err := jpeg.Encode(&buf, i.img, &jpeg.Options{Quality: 1}); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("无法压缩到 %d 字节以内，当前尺寸下最小为 %d 字节", maxBytes, buf.Len())
	}

	return best, nil
}

// CompressToFile 将图片压缩到指定大小以内并保存到文件。
// 内部调用 Compress 方法进行压缩，然后将结果写入指定路径。
// 输出格式固定为 JPEG，文件权限为 0644。
func (i *Image) CompressToFile(maxBytes int64, path string) error {
	data, err := i.Compress(maxBytes)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// CompressAuto 智能压缩图片到目标大小以内（单位：字节）。
// 首先尝试仅通过降低 JPEG 质量来达到目标大小；
// 如果质量降到最低仍然超出限制，则逐步缩小图片尺寸（每次缩小 10%），
// 再配合质量调节来满足大小要求。
// 适用于需要严格控制文件大小的场景，如上传限制、缩略图生成等。
func (i *Image) CompressAuto(maxBytes int64) ([]byte, error) {
	if maxBytes <= 0 {
		return nil, fmt.Errorf("maxBytes 必须为正数")
	}

	// 第一步：尝试仅通过降低质量来满足大小要求
	data, err := i.Compress(maxBytes)
	if err == nil {
		return data, nil
	}

	// 第二步：逐步缩小图片尺寸（从 90% 到 10%），每次缩小后再尝试压缩
	for scale := 0.9; scale >= 0.1; scale -= 0.1 {
		newW := int(float64(i.Width()) * scale)
		newH := int(float64(i.Height()) * scale)
		if newW < 1 || newH < 1 {
			continue
		}
		resized := i.Resize(newW, newH)
		data, err = resized.Compress(maxBytes)
		if err == nil {
			return data, nil
		}
	}

	return nil, fmt.Errorf("即使缩小到最小尺寸也无法压缩到 %d 字节以内", maxBytes)
}
