package controllers

import (
	"github.com/sciter-sdk/go-sciter"
	"github.com/sciter-sdk/go-sciter/window"
	"log"
)

func NewWindows(file, title string) *window.Window {
	w, err := window.New(
		sciter.SW_TITLEBAR | sciter.SW_MAIN | sciter.SW_CONTROLS | sciter.SW_ENABLE_DEBUG,
		&sciter.Rect{200, 100, 726, 600})
	if err != nil {
		log.Fatal(err)
	}
	//加载文件
	w.LoadFile(file)
	//设置标题
	w.SetTitle(title)
	return w
}
