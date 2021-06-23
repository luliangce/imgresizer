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
	if err = jpeg.Encode(dst, src, &jpeg.Options{Quality: Quality}); err != nil {
		return err
	}
	return err
}
func scaleTo(src image.Image, scale draw.Scaler, ratio int) image.Image {
	if ratio > 100 {
		ratio = 100
	} else if ratio < 1 {
		ratio = 1
	}

	rect := image.Rect(0, 0, src.Bounds().Max.X*ratio/100, src.Bounds().Max.Y*ratio/100)
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

func Resize(filename, dstDir string, Quality int, ratio int, scaler draw.Interpolator) error {
	img, err := openImg(filename)
	_, pureName := filepath.Split(filename)
	targetFile := filepath.Join(dstDir, pureName)
	if err != nil {
		return err
	}
	dst := scaleTo(img, scaler, ratio)
	return saveImg(dst, targetFile, Quality, true)
}
