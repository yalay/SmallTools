package main

import (
	"controllers"
)

func main() {
	w:=controllers.NewWindows("views/demo1.html", "test")
	//获取根元素
	root, _ := w.GetRootElement()
	//设置元素处理程序
	controllers.ThumbHandlers(root)
	//显示窗口
	w.Show()
	//运行窗口，进入消息循环
	w.Run()
}
