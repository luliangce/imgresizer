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
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"

	"github.com/luliangce/imgresizer"
	"github.com/spf13/cobra"
)

var (
	dstDir  string
	quality int
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

var rootCmd = &cobra.Command{
	Use:   "imgresizer img1 [img2] [img3]",
	Short: "将输入图片压缩为同尺寸的低质量jpg",
	Long:  `将输入图片压缩为同尺寸的低质量jpg`,
	Args:  cobra.MinimumNArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		if runtime.GOOS == "windows" {
			//windows glob compatible
			m, err := filepath.Glob(args[0])
			if err != nil {
				log.Fatal(err)
			}
			args = m
		}
		imgs := findImg(args...)

		err := os.MkdirAll(dstDir, os.ModePerm)
		if err != nil && !os.IsExist(err) {
			log.Fatal(err)
		}
		for i := 0; i < len(imgs); i++ {
			err = imgresizer.Resize(imgs[i], dstDir, quality)
			if err != nil {
				log.Fatal(err)
			}
		}
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {

	rootCmd.Flags().StringVarP(&dstDir, "destination", "d", "resized", "目标文件夹")
	rootCmd.Flags().IntVarP(&quality, "quality", "q", 75, "图片质量")

}
