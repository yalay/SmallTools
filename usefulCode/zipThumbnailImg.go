package main

import (
	"bytes"
	"flag"
	"image"
	"image/jpeg"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/jhoonb/archivex"
	"github.com/nfnt/resize"
)

var (
	folderPath string
)

func init() {
	flag.StringVar(&folderPath, "f", "", "f=test")
	flag.Parse()
}

// Example using only func zip
func zip(folder string) {
	zip := new(archivex.ZipFile)
	zip.Create(path.Base(folder) + ".zip")
	imgPaths := getAllImgPath(folder)
	for _, imgPath := range imgPaths {
		img, err := loadImage(imgPath)
		if err != nil {
			log.Printf("loadImage err:%v\n", err)
			continue
		}
		imgThumb := thumbnailSimple(800, 0, img)

		buff := &bytes.Buffer{}
		jpeg.Encode(buff, imgThumb, nil)
		zip.Add(path.Base(imgPath), buff.Bytes())
	}
	zip.Add("tengmm.com.txt", nil)
	zip.Close()
}

func loadImage(imgPath string) (img image.Image, err error) {
	file, err := os.Open(imgPath)
	if err != nil {
		return
	}
	defer file.Close()
	img, _, err = image.Decode(file)
	return
}

// 简单的缩放,指定最大宽和高
func thumbnailSimple(maxWidth, maxHeight uint, img image.Image) image.Image {
	oriBounds := img.Bounds()
	oriWidth := uint(oriBounds.Dx())
	oriHeight := uint(oriBounds.Dy())

	if maxWidth == 0 {
		maxWidth = oriWidth
	}

	if maxHeight == 0 {
		maxHeight = oriHeight
	}
	return resize.Thumbnail(maxWidth, maxHeight, img, resize.Lanczos3)
}

func getAllImgPath(folder string) []string {
	files, err := ioutil.ReadDir(folder)
	if err != nil {
		return nil
	}
	paths := make([]string, 0, len(files))
	for _, file := range files {
		if file.IsDir() {
			continue
		} else {
			fileName := file.Name()
			lowerFileName := strings.ToLower(fileName)
			if !strings.HasSuffix(lowerFileName, ".jpg") &&
				!strings.HasSuffix(lowerFileName, ".jpeg") &&
				!strings.HasSuffix(lowerFileName, ".png") {
				log.Printf("not img:%s\n", fileName)
			}
			paths = append(paths, folder+"/"+fileName)
		}
	}
	return paths
}

func main() {
	if len(os.Args) < 2 {
		return
	}

	if folderPath == "" {
		folderPath = os.Args[1]
	}
	zip(folderPath)
}
