package main

import (
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/tealeg/xlsx"
	"os"
	"strings"
)

type AdviewSlot struct {
	appName string
	appId   string
	os      string
	pkgName string
}

var filePath = ""

func init() {
	flag.StringVar(&filePath, "file", "", "file=test.xlsx")
	flag.Parse()
}

func main() {
	if len(os.Args) < 2 {
		return
	}

	if filePath == "" {
		filePath = os.Args[1]
	}

	count := 0
	adviewSlots := readAdviewSlots(filePath)
	for key, info := range adviewSlots {
		count++
		if count > 10 {
			break
		}
		fmt.Printf("key:%s info:%+v\n", key, info)
	}
	excelFile, err := xlsx.OpenFile("adview未知广告位.xlsx")
	if err != nil {
		fmt.Printf("open err:%v\n", err)
		return
	}

	if len(excelFile.Sheets) == 0 {
		return
	}
	sheet := excelFile.Sheets[0]
	for _, row := range sheet.Rows {
		if len(row.Cells) < 4 {
			continue
		}
		appEncId, _ := row.Cells[1].String()
		if info, ok := adviewSlots[appEncId]; ok {
			row.Cells[3].Value = info.appName
			fmt.Printf("appEncId:%s\n", appEncId)
		}
	}
	excelFile.Save("MyXLSXFile.xlsx")
}

func readAdviewSlots(filePath string) map[string]*AdviewSlot {
	excelFile, err := xlsx.OpenFile(filePath)
	if err != nil {
		fmt.Printf("open err:%v\n", err)
		return nil
	}

	adviewSlots := make(map[string]*AdviewSlot)
	for _, sheet := range excelFile.Sheets {
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
				appName: appName,
				appId:   appId,
				os:      os,
				pkgName: pkgName,
			}

			adviewSlots[Md5Sum(appId)] = adviewSlot
		}
	}
	return adviewSlots
}

func Md5Sum(key string) string {
	h := md5.New()
	h.Write([]byte(key))
	return hex.EncodeToString(h.Sum(nil))
}
