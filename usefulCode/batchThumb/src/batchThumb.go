package main

import (
	"bytes"
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/nfnt/resize"
)

var (
	folderPath   string
	outputFolder string
)

func init() {
	flag.StringVar(&folderPath, "f", "", "f=test")
	flag.StringVar(&outputFolder, "o", "thumbnail", "o=thumbnail")
	flag.Parse()
}

// [LeYuan]Vol.003[50P] xxxx
func batchThumbnail(srcFolder, dstFolder string) {
	imgPaths := getAllImgPath(srcFolder)
	if len(imgPaths) == 0 {
		return
	}

	baseName := path.Base(srcFolder)
	nameFields := strings.SplitN(baseName, " ", 3)
	if len(nameFields) < 3 {
		return
	}

	newName := fmt.Sprintf("%s%s[%dP]", nameFields[0], nameFields[1], len(imgPaths))
	newFolder := path.Join(dstFolder, newName)
	os.MkdirAll(newFolder, 0755)
	for _, imgPath := range imgPaths {
		img, err := loadImage(imgPath)
		if err != nil {
			log.Printf("loadImage err:%v\n", err)
			continue
		}

		imgName := path.Base(imgPath)
		newFile, err := os.Create(path.Join(newFolder, imgName))
		if err != nil {
			log.Printf("new file err:%v\n", err)
			continue
		}

		imgThumb := thumbnailSimple(800, 0, img)
		buff := &bytes.Buffer{}
		jpeg.Encode(buff, imgThumb, nil)
		buff.WriteTo(newFile)
		newFile.Close()
	}
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
	batchThumbnail(folderPath, outputFolder)
}
