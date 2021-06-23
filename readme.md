# WIP:imgresizer

## description

一个将图片质量降低的工具

## install

clone and run 
```shell
$ go get
$ go install bin/imgresizer.go
```

this cmd will install a binary to your GOBIN(e.g `~/go/bin`),then

```shell
➜  imgresizer : master ✔ : ᐅ  imgresizer --help
将输入图片压缩为同尺寸的低质量jpg

Usage:
  imgresizer img1 [img2] [img3] [flags]

Flags:
  -d, --destination string   目标文件夹 (default "resized")
  -h, --help                 help for imgresizer
  -q, --quality int          图片质量 (default 75)
```

## TODO

- [ ] 指定算法
- [ ] 指定尺寸缩放比例
- [ ] 测试