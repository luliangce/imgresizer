/*
Copyright © 2021 luliangce@gmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/luliangce/imgresizer"
	"github.com/spf13/cobra"
	"golang.org/x/image/draw"
	"golang.org/x/sync/semaphore"
)

var (
	dstDir  string
	quality int
	ratio   int
	scaler  string
)

func findImg(files ...string) []string {
	imgs := []string{}
	cpattern := regexp.MustCompile(`.+\.(?:png|jpg|jpeg)`)
	for i := 0; i < len(files); i++ {
		if cpattern.Match([]byte(files[i])) {
			imgs = append(imgs, files[i])
		}
	}
	return imgs
}

func stat(done chan string) {
	go func() {
		startTime := time.Now()
		finishCount := 0
		for range done {
			finishCount++
			if finishCount%10 == 0 && finishCount > 0 {
				log.Printf("%d images resized in %v,%.2f/s",
					finishCount,
					time.Since(startTime),
					float64(finishCount)/time.Since(startTime).Seconds())
			}
		}
	}()
}

var rootCmd = &cobra.Command{
	Use:   "imgresizer img1 [img2] [img3]",
	Short: "将输入图片压缩为指定尺寸比例的jpg",
	Long:  `将输入图片压缩为指定尺寸比例的jpg`,
	Args:  cobra.MinimumNArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		log.Printf("image's width and height will scale to [%d%%] with [%d%%] of quality", ratio, quality)

		interpolator := map[string]draw.Interpolator{
			"A": draw.ApproxBiLinear,
			"N": draw.NearestNeighbor,
			"B": draw.BiLinear,
			"C": draw.CatmullRom,
		}[strings.ToUpper(scaler)]
		if interpolator == nil {
			log.Fatal("wrong scaler,only N/A/B/C is available")
		}

		if runtime.GOOS == "windows" {
			//windows glob compatible
			m, err := filepath.Glob(args[0])
			if err != nil {
				log.Fatal(err)
			}
			args = m
		}
		imgs := findImg(args...)
		if len(imgs) == 0 {
			log.Printf("img not found")
			return
		}
		log.Printf("%d image(s) will be resized and save to directory [%s]", len(imgs), dstDir)
		err := os.MkdirAll(dstDir, os.ModePerm)
		if err != nil && !os.IsExist(err) {
			log.Fatal(err)
		}

		wg := new(sync.WaitGroup)

		done := make(chan string, len(imgs))
		stat(done) // print process info
		sem := semaphore.NewWeighted(10)
		handle := func(img string) {
			err = imgresizer.Resize(img, dstDir, quality, ratio, interpolator)
			if err != nil {
				log.Fatal(err)
			}
			done <- img
			sem.Release(1)
			wg.Done()
		}

		for i := 0; i < len(imgs); i++ {
			wg.Add(1)
			sem.Acquire(context.Background(), 1)
			handle(imgs[i])
		}
		wg.Wait()
		close(done)
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {

	rootCmd.Flags().StringVarP(&dstDir, "destination", "d", "resized", "目标文件夹")
	rootCmd.Flags().IntVarP(&quality, "quality", "q", 75, "图片质量")
	rootCmd.Flags().IntVarP(&ratio, "ratio", "r", 100, "相对于原图的尺寸比例,0~100")

	scaleDesc := `需要使用的缩放算法,可以使用N/A/B/C 三种，
N -NearestNeighbor
	NearestNeighbor is the nearest neighbor interpolator. It is very fast,
	but usually gives very low quality results. When scaling up, the result
	will look 'blocky'.

A -ApproxBiLinear
	ApproxBiLinear is a mixture of the nearest neighbor and bi-linear
	interpolators. It is fast, but usually gives medium quality results.
	
	It implements bi-linear interpolation when upscaling and a bi-linear
	blend of the 4 nearest neighbor pixels when downscaling. This yields
	nicer quality than nearest neighbor interpolation when upscaling, but
	the time taken is independent of the number of source pixels, unlike the
	bi-linear interpolator. When downscaling a large image, the performance
	difference can be significant.

B -BiLinear
	BiLinear is the tent kernel. It is slow, but usually gives high quality results.

C -CatmullRom
	CatmullRom is the Catmull-Rom kernel. It is very slow, but usually gives
	very high quality results.

	It is an instance of the more general cubic BC-spline kernel with parameters
	B=0 and C=0.5. See Mitchell and Netravali, "Reconstruction Filters in
	Computer Graphics", Computer Graphics, Vol. 22, No. 4, pp. 221-228.
	`
	rootCmd.Flags().StringVarP(&scaler, "scaler", "s", "A", scaleDesc)

}
