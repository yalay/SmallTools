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
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/nfnt/resize"
)

var (
	folderPath       string
	parentFolderPath string
	outputFolder     string
)

func init() {
	flag.StringVar(&folderPath, "f", "", "f=test")
	flag.StringVar(&outputFolder, "o", "thumbnail", "o=thumbnail")
	flag.StringVar(&parentFolderPath, "p", "", "p=parent")
	flag.Parse()
}

func batchThumbnail(srcFolders []string, dstFolder string) {
	var wg sync.WaitGroup
	for _, srcFolder := range srcFolders {
		wg.Add(1)
		go func(folder string) {
			thumbnail(folder, dstFolder)
			log.Println("thumbnail end:" + folder)
			wg.Done()
		}(srcFolder)
	}
	wg.Wait()
}

// [LeYuan]Vol.003[50P] xxxx
func thumbnail(srcFolder, dstFolder string) {
	imgPaths := getAllImgPath(srcFolder)
	if len(imgPaths) == 0 {
		return
	}

	baseName := filepath.Base(srcFolder)
	nameFields := strings.SplitN(baseName, " ", 3)
	if len(nameFields) < 3 {
		return
	}

	catName := nameFields[0]
	catName = strings.Trim(catName, "[]")
	newName := fmt.Sprintf("%s-%s-%dP", catName, nameFields[1], len(imgPaths))
	newFolder := filepath.Join(dstFolder, newName)
	os.MkdirAll(newFolder, 0755)
	for _, imgPath := range imgPaths {
		img, err := loadImage(imgPath)
		if err != nil {
			log.Printf("loadImage err:%v\n", err)
			continue
		}

		imgName := filepath.Base(imgPath)
		newFile, err := os.Create(filepath.Join(newFolder, imgName))
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
				continue
			}
			paths = append(paths, filepath.Join(folder, fileName))
		}
	}
	return paths
}

func getChildrenFolderPath(parentPath string) []string {
	folders, err := ioutil.ReadDir(parentPath)
	if err != nil {
		time.Sleep(5 * time.Second)
		return nil
	}

	rspFolders := make([]string, 0, len(folders))
	for _, folder := range folders {
		if !folder.IsDir() {
			continue
		}
		rspFolders = append(rspFolders, filepath.Join(parentPath, folder.Name()))
	}
	return rspFolders
}

func main() {
	if len(os.Args) < 2 {
		return
	}

	if folderPath == "" {
		folderPath = os.Args[1]
	}

	if parentFolderPath != "" {
		foldersPath := getChildrenFolderPath(parentFolderPath)
		batchThumbnail(foldersPath, outputFolder)
	} else {
		thumbnail(folderPath, outputFolder)
	}
}
