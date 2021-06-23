package imgresizer

import (
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"

	"golang.org/x/image/draw"
)

func saveImg(src image.Image, filepath string, Quality int, force ...bool) error {

	_, err := os.Stat(filepath)
	if os.IsExist(err) && len(force) > 0 && force[0] {
		return fmt.Errorf("该文件已存在")
	}
	dst, err := os.Create(filepath)
	defer func() {
		err = dst.Close()
	}()
	if err != nil {
		return err
	}
	if err = jpeg.Encode(dst, src, &jpeg.Options{Quality: 75}); err != nil {
		return err
	}
	return err
}
func scaleTo(src image.Image, scale draw.Scaler) image.Image {
	rect := image.Rect(0, 0, src.Bounds().Max.X, src.Bounds().Max.Y)
	dst := image.NewRGBA(rect)
	scale.Scale(dst, rect, src, src.Bounds(), draw.Over, nil)
	return dst
}

func openImg(filepath string) (image.Image, error) {
	fl, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer fl.Close()
	img, _, err := image.Decode(fl)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func Resize(filename, dstDir string, Quality int) error {
	img, err := openImg(filename)
	_, pureName := filepath.Split(filename)
	targetFile := filepath.Join(dstDir, pureName)
	if err != nil {
		return err
	}
	dst := scaleTo(img, draw.ApproxBiLinear)
	return saveImg(dst, targetFile, Quality, true)
}
