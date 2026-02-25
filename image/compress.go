package image

import (
	"bytes"
	"fmt"
	"image/jpeg"
	"os"
)

// Compress encodes the image as JPEG and uses binary search to find
// the highest quality that keeps the output under maxBytes
func (i *Image) Compress(maxBytes int64) ([]byte, error) {
	if maxBytes <= 0 {
		return nil, fmt.Errorf("maxBytes must be positive")
	}

	lo, hi := 1, 100
	var best []byte

	for lo <= hi {
		mid := (lo + hi) / 2
		var buf bytes.Buffer
		if err := jpeg.Encode(&buf, i.img, &jpeg.Options{Quality: mid}); err != nil {
			return nil, err
		}
		if int64(buf.Len()) <= maxBytes {
			best = buf.Bytes()
			lo = mid + 1
		} else {
			hi = mid - 1
		}
	}

	if best == nil {
		var buf bytes.Buffer
		if err := jpeg.Encode(&buf, i.img, &jpeg.Options{Quality: 1}); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("cannot compress to %d bytes, minimum size at current dimensions is %d bytes", maxBytes, buf.Len())
	}

	return best, nil
}

// CompressToFile compresses the image and saves the result to a file
func (i *Image) CompressToFile(maxBytes int64, path string) error {
	data, err := i.Compress(maxBytes)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// CompressAuto compresses the image to the target size,
// automatically reducing dimensions if quality reduction alone is not enough
func (i *Image) CompressAuto(maxBytes int64) ([]byte, error) {
	if maxBytes <= 0 {
		return nil, fmt.Errorf("maxBytes must be positive")
	}

	data, err := i.Compress(maxBytes)
	if err == nil {
		return data, nil
	}

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

	return nil, fmt.Errorf("cannot compress to %d bytes even at minimum dimensions", maxBytes)
}
