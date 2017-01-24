package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/jhoonb/archivex"
)

var (
	folderPath       string
	parentFolderPath string
	outputFolder     string
)

func init() {
	flag.StringVar(&folderPath, "f", "", "f=test")
	flag.StringVar(&parentFolderPath, "p", "", "f=parent")
	flag.StringVar(&outputFolder, "o", "zip", "o=zip")
	flag.Parse()
}

func batchZipAd(srcFolders []string, dstFolder string) {
	var wg sync.WaitGroup
	for _, srcFolder := range srcFolders {
		wg.Add(1)
		go func(folder string) {
			zipAd(folder, dstFolder)
			log.Println("zip end:" + folder)
			wg.Done()
		}(srcFolder)
	}
	wg.Wait()
}

func zipAd(srcFolder, dstFolder string) {
	imgPaths := getAllImgPath(srcFolder)
	if len(imgPaths) == 0 {
		return
	}

	baseName := filepath.Base(srcFolder)
	nameFields := strings.SplitN(baseName, " ", 3)
	if len(nameFields) < 3 {
		return
	}

	os.MkdirAll(dstFolder, 0755)
	zipFile := &archivex.ZipFile{}
	newName := fmt.Sprintf("%s%s[%dP].zip", nameFields[0], nameFields[1], len(imgPaths))
	err := zipFile.Create(filepath.Join(dstFolder, newName))
	if err != nil {
		log.Printf("create err:%v\n", err)
		time.Sleep(5 * time.Second)
		return
	}

	for _, imgPath := range imgPaths {
		data, err := ioutil.ReadFile(imgPath)
		if err != nil {
			log.Printf("read img err:%v\n", err)
			continue
		}
		zipFile.Add(filepath.Base(imgPath), data)
		time.Sleep(5 * time.Second)
	}
	zipFile.Add("tengmm.com.txt", []byte("资源来自http://tengmm.com，欢迎回访。"))
	zipFile.Close()
}

func getAllImgPath(folder string) []string {
	files, err := ioutil.ReadDir(folder)
	if err != nil {
		time.Sleep(5 * time.Second)
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
				!strings.HasSuffix(lowerFileName, ".png") &&
				!strings.HasSuffix(lowerFileName, ".gif") {
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
		batchZipAd(foldersPath, outputFolder)
	} else {
		zipAd(folderPath, outputFolder)
	}
}
