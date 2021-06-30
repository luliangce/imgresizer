# imgresizer

## description

将`golang.org/x/image/draw`中的四种缩放算法封装而来的命令行工具，初衷是用于批量减小图片体积。

## install

clone and run 
```shell
$ go get
$ go install bin/imgresizer.go
```

this cmd will install a binary to your GOBIN(e.g `~/go/bin`),then

```shell
➜  imgresizer : master ✔ : ᐅ  imgresizer --help
输入图片压缩为指定尺寸比例的jpg

Usage:
  imgresizer img1 [img2] [img3] [flags]

Flags:
  -d, --destination string   目标文件夹 (default "resized")
  -h, --help                 help for imgresizer
  -q, --quality int          图片质量 (default 75)
  -r, --ratio int            相对于原图的尺寸比例,0~100 (default 100)
  -s, --scaler string        需要使用的缩放算法,可以使用N/A/B/C 三种，
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
                                 (default "A")
```

## TODO

- [X] 指定算法
- [X] 指定尺寸缩放比例
- [ ] ~~测试~~