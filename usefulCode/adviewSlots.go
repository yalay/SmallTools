package main

import (
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/tealeg/xlsx"
	"os"
	"regexp"
	"strings"
	"time"
)

type AdviewSlot struct {
	slotType string
	appName  string
	appId    string
	os       string
	pkgName  string
}

var unknownSlotFile = ""
var adviewFilePath = ""
var sizeReg = regexp.MustCompile(`^\d+x\d+$`)

func init() {
	flag.StringVar(&unknownSlotFile, "u", "", "u=unknown.xlsx")
	flag.StringVar(&adviewFilePath, "a", "adview.xlsx", "a=adview.xlsx")
	flag.Parse()
}

func main() {
	if len(os.Args) < 2 {
		return
	}

	if unknownSlotFile == "" {
		unknownSlotFile = os.Args[1]
	}

	adviewSlots, err := readAdviewSlots(adviewFilePath)
	if err != nil {
		failedNotice(err.Error())
	}

	excelFile, err := xlsx.OpenFile(unknownSlotFile)
	if err != nil {
		fmt.Printf("open err:%v\n", err)
		return
	}

	if len(excelFile.Sheets) == 0 {
		return
	}
	sheet := excelFile.Sheets[0]
	for _, row := range sheet.Rows {
		curLen := len(row.Cells)
		if curLen < 2 {
			continue
		}

		if curLen < 8 {
			for i := 0; i < 8-curLen; i++ {
				row.AddCell()
			}
		}

		slotSize, _ := row.Cells[0].String()
		if !sizeReg.MatchString(slotSize) {
			continue
		}

		appEncId, _ := row.Cells[1].String()
		if info, ok := adviewSlots[appEncId]; ok {
			row.Cells[3].Value = slotSize + "_" + appEncId
			row.Cells[4].Value = info.slotType
			row.Cells[5].Value = info.appName
			row.Cells[6].Value = info.os
			row.Cells[7].Value = info.pkgName
		}
	}
	filePathPre := unknownSlotFile[:(len(unknownSlotFile) - 4)]
	excelFile.Save(filePathPre + "已匹配.xlsx")
}

func readAdviewSlots(filePath string) (map[string]*AdviewSlot, error) {
	excelFile, err := xlsx.OpenFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("adview提供广告位信息文件请放在当前目录并且重命名为" + adviewFilePath)
	}

	adviewSlots := make(map[string]*AdviewSlot)
	for _, sheet := range excelFile.Sheets {
		slotType := getAdType(sheet.Name)
		for _, row := range sheet.Rows {
			if len(row.Cells) < 4 {
				continue
			}

			appId, _ := row.Cells[1].String()
			if appId == "" || !strings.HasPrefix(appId, "SDK") {
				continue
			}

			appName, _ := row.Cells[0].String()
			os, _ := row.Cells[2].String()
			pkgName, _ := row.Cells[3].String()
			adviewSlot := &AdviewSlot{
				slotType: slotType,
				appName:  appName,
				appId:    appId,
				os:       os,
				pkgName:  pkgName,
			}

			adviewSlots[Md5Sum(appId)] = adviewSlot
		}
	}
	return adviewSlots, nil
}

func Md5Sum(key string) string {
	h := md5.New()
	h.Write([]byte(key))
	return hex.EncodeToString(h.Sum(nil))
}

func getAdType(sheet string) string {
	switch {
	case strings.Contains(sheet, "banner"):
		return "banner"
	case strings.Contains(sheet, "插屏"):
		return "插屏"
	case strings.Contains(sheet, "开屏"):
		return "开屏"
	case strings.Contains(sheet, "原生"):
		return "原生"
	}
	return "未知"
}

func failedNotice(info string) {
	fmt.Println("")
	fmt.Println("")
	fmt.Printf("注意：%s\n", info)
	fmt.Println("")
	fmt.Println("")
	fmt.Println("5秒后自动退出")
	time.Sleep(5 * time.Second)
	os.Exit(0)
}
