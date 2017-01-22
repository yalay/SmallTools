package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/jhoonb/archivex"
)

var (
	folderPath   string
	outputFolder string
)

func init() {
	flag.StringVar(&folderPath, "f", "", "f=test")
	flag.StringVar(&outputFolder, "o", "zip", "o=zip")
	flag.Parse()
}

// Example using only func zip
func zip(srcFolder, dstFolder string) {
	imgPaths := getAllImgPath(srcFolder)
	if len(imgPaths) == 0 {
		return
	}

	baseName := path.Base(srcFolder)
	nameFields := strings.SplitN(baseName, " ", 3)
	if len(nameFields) < 3 {
		return
	}

	os.MkdirAll(dstFolder, 0755)
	zip := new(archivex.ZipFile)
	newName := fmt.Sprintf("%s%s[%dP].zip", nameFields[0], nameFields[1], len(imgPaths))
	zip.Create(path.Join(dstFolder, newName))
	for _, imgPath := range imgPaths {
		data, err := ioutil.ReadFile(imgPath)
		if err != nil {
			log.Printf("read img err:%v\n", err)
			continue
		}
		zip.Add(path.Base(imgPath), data)
	}
	zip.Add("tengmm.com.txt", nil)
	zip.Close()
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
			paths = append(paths, path.Join(folder, fileName))
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
	zip(folderPath, outputFolder)
}
