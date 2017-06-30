package controllers

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/jhoonb/archivex"
	"github.com/sciter-sdk/go-sciter"
)

//设置元素的处理程序
func SetElementHandlers(root *sciter.Element) {
	btn1, _ := root.SelectById("btn1")
	//处理元素简单点击事件
	btn1.OnClick(func() {
		srcDir, err := root.SelectById("srcdir")
		if err != nil {
			fmt.Println(err)
			return
		}

		srcText, err := srcDir.Text()
		if err != nil {
			fmt.Println(err)
			return
		}

		srcText = strings.TrimSpace(srcText)
		if srcText == "" {
			fmt.Println("input empty")
			return
		}

		dstDir, err := root.SelectById("dstdir")
		if err != nil {
			fmt.Println(err)
			return
		}

		dstText, err := dstDir.Text()
		if err != nil {
			fmt.Println(err)
			return
		}
		dstText = strings.TrimSpace(dstText)
		if dstText == "" {
			dstText = filepath.Dir(srcText)
			dstDir.SetText(dstText)
		}

		var includeCurrentFolder bool
		includeDir, err := root.SelectById("includecur")
		if err == nil {
			state, err := includeDir.State()
			if err == nil && state&sciter.STATE_CHECKED == sciter.STATE_CHECKED {
				fmt.Println("includeCurrentFolder")
				includeCurrentFolder = true
			}
		}

		btn1.SetState(sciter.STATE_DISABLED, 0, true)
		zipDir(srcText, dstText, includeCurrentFolder)
		btn1.SetState(0, sciter.STATE_DISABLED, true)
	})
}

func zipDir(srcDir, dstDir string, includeCurrentFolder bool) {
	dirName := filepath.Base(srcDir)
	zip := new(archivex.ZipFile)
	zip.Create(filepath.Join(dstDir, dirName+".zip"))
	zip.AddAll(srcDir, includeCurrentFolder)
	zip.Close()
	time.Sleep(5*time.Second)
}
