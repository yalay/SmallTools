package controllers

import (
	"encoding/json"
	"fmt"
	"image"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/disintegration/imaging"
	"github.com/sciter-sdk/go-sciter"
	"strconv"
)

const (
	kTypeFile = iota
	kTypeDir
)

const (
	kUnitPx = iota
	kUnitPercent
)

type DragMsg struct {
	Type     int32  `json:"type"`
	FullPath string `json:"fullpath"`
	FileNum  int32  `json:"fileNum"`
	Size     int64  `json:"size"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	Err      error  `json:"err"`
}

type ExecReqMsg struct {
	Path    string `json:"path"`
	Unit    int    `json:"unit"`
	Format  int    `json:"format"`
	Width   string `json:"width"`
	Height  string `json:"height"`
	Publish string `json:"publish"`
}

type ExecRspMsg struct {
	Err string `json:"err"`
}

func ThumbHandlers(root *sciter.Element) {
	dropZone, _ := root.SelectById("drop-zone")
	dropZone.DefineMethod("showAttr", func(args ...*sciter.Value) *sciter.Value {
		if len(args) == 0 {
			return nil
		}

		dragData := args[0].String()
		dragData = strings.TrimSpace(dragData)
		if dragData == "" {
			return nil
		}

		dragDataPath := strings.TrimPrefix(dragData, "file://")
		dragDataInfo, err := os.Stat(dragDataPath)
		if err != nil {
			fmt.Println(err)
			return nil
		}

		msg := &DragMsg{
			FullPath: dragDataPath,
		}
		if dragDataInfo.IsDir() {
			var fileNum int32
			var firstImgInfo os.FileInfo
			fileInfos, _ := ioutil.ReadDir(dragDataPath)
			for _, fileInfo := range fileInfos {
				if isImgFile(fileInfo.Name()) {
					if firstImgInfo == nil {
						firstImgInfo = fileInfo
					}
					fileNum++
				}
			}
			fmt.Println(firstImgInfo)
			msg.Type = kTypeDir
			msg.FileNum = fileNum
			if firstImgInfo != nil {
				msg.Size = firstImgInfo.Size()
				file, err := os.Open(filepath.Join(dragDataPath, firstImgInfo.Name()))
				if err == nil {
					img, _, err := image.Decode(file)
					if err == nil {
						msg.Width = img.Bounds().Dx()
						msg.Height = img.Bounds().Dy()
					} else {
						msg.Err = err
					}
				} else {
					msg.Err = err
				}
			}
		} else {
			msg.Type = kTypeFile
			msg.Size = dragDataInfo.Size()
			file, err := os.Open(dragDataPath)
			if err == nil {
				img, _, err := image.Decode(file)
				if err == nil {
					msg.Width = img.Bounds().Dx()
					msg.Height = img.Bounds().Dy()
				} else {
					msg.Err = err
				}
			} else {
				msg.Err = err
			}
		}
		msgData, _ := json.MarshalIndent(msg, "", "  ")
		fmt.Println(string(msgData))
		return sciter.NewValue(string(msgData))
	})

	execButton, _ := root.SelectById("exec")
	execButton.DefineMethod("exec", func(args ...*sciter.Value) *sciter.Value {
		if len(args) == 0 {
			return nil
		}
		reqMsg := &ExecReqMsg{}
		rspMsg := ExecRspMsg{}
		execData := args[0].String()
		err := json.Unmarshal([]byte(execData), reqMsg)
		if err != nil {
			rspMsg.Err = err.Error()
		} else {
			if reqMsg.Path == "" {
				rspMsg.Err = "没有选择文件"
			} else {
				switch reqMsg.Unit {
				case kUnitPx:
					srcImage, err := imaging.Open(reqMsg.Path)
					if err != nil {
						rspMsg.Err = err.Error()
					} else {
						fmt.Println(reqMsg)
						width, _ := strconv.Atoi(reqMsg.Width)
						height, _ := strconv.Atoi(reqMsg.Height)
						if width == 0 || height == 0 {
							rspMsg.Err = "宽高选择错误"
						} else {
							lastDot := strings.LastIndex(reqMsg.Path, ".")
							dstImage := imaging.Resize(srcImage, width, height, imaging.Lanczos)
							imaging.Save(dstImage, reqMsg.Path[:lastDot]+".thumb"+reqMsg.Path[lastDot:])
						}
					}
				}
			}
		}
		msgData, _ := json.MarshalIndent(rspMsg, "", "  ")
		return sciter.NewValue(string(msgData))
	})
}

func isImgFile(fileName string) bool {
	fileExt := strings.ToLower(filepath.Ext(fileName))
	if fileExt == ".jpg" || fileExt == ".jpeg" ||
		fileExt == ".png" || fileExt == ".gif" ||
		fileExt == ".webp" {
		return true
	}
	return false
}
