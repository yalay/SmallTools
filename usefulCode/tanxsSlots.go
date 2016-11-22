package main

import (
	"flag"
	"fmt"
	"github.com/tealeg/xlsx"
	"os"
	"strings"
	"time"
)

type MogoSlot struct {
	slotType string
	appName  string
	os       string
	pkgName  string
}

type TanxSlot struct {
	publishName string
	slotName    string
	pid         string
	domain      string
	publishCat  string
	slotType    string
	devType     string
}

var filePath = ""

func init() {
	flag.StringVar(&filePath, "file", "", "file=unknown.xlsx")
	flag.Parse()
}

func main() {
	if len(os.Args) < 2 {
		return
	}

	if filePath == "" {
		filePath = os.Args[1]
	}

	mogoSlots, err := readMogoSlots("mogo.xlsx")
	if err != nil {
		failedNotice(err)
	}

	tanxSlots, err := readTanxSlots("tanx.xlsx")
	if err != nil {
		failedNotice(err)
	}

	excelFile, err := xlsx.OpenFile(filePath)
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
		if curLen < 4 {
			continue
		}

		slotId, _ := row.Cells[0].String()
		slotIdFields := strings.Split(strings.TrimSpace(slotId), "_")
		fieldsLen := len(slotIdFields)
		if fieldsLen > 4 {
			pkgName := strings.ToLower(slotIdFields[fieldsLen-1])
			mogoInfo, ok := mogoSlots[pkgName]
			if ok {
				if curLen < 7 {
					for i := 0; i < 7-curLen; i++ {
						row.AddCell()
					}
				}
				row.Cells[4].Value = mogoInfo.appName
				row.Cells[5].Value = mogoInfo.os
				row.Cells[6].Value = mogoInfo.slotType
			}
		} else if fieldsLen == 4 {
			tanxInfo, ok := tanxSlots[slotId]
			if ok {
				if curLen < 10 {
					for i := 0; i < 10-curLen; i++ {
						row.AddCell()
					}
				}
				row.Cells[4].Value = tanxInfo.publishName
				row.Cells[6].Value = tanxInfo.slotType
				row.Cells[7].Value = tanxInfo.domain
				row.Cells[8].Value = tanxInfo.publishCat
				row.Cells[9].Value = tanxInfo.devType
			}
		}
	}
	filePathPre := filePath[:(len(filePath) - 4)]
	excelFile.Save(filePathPre + "已匹配.xlsx")
}

func readMogoSlots(filePath string) (map[string]*MogoSlot, error) {
	excelFile, err := xlsx.OpenFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("%s 没有mogo.xlsx文件")
	}

	mogoSlots := make(map[string]*MogoSlot)
	for _, sheet := range excelFile.Sheets {
		for i, row := range sheet.Rows {
			if i == 0 || len(row.Cells) < 4 {
				continue
			}

			appName, _ := row.Cells[0].String()
			os, _ := row.Cells[1].String()
			slotType, _ := row.Cells[2].String()
			pkgName, _ := row.Cells[3].String()

			pkgName = strings.TrimSpace(pkgName)
			if pkgName == "" {
				continue
			}

			mogoSlot := &MogoSlot{
				slotType: getAdType(slotType),
				appName:  appName,
				os:       os,
				pkgName:  pkgName,
			}
			mogoSlots[strings.ToLower(pkgName)] = mogoSlot
		}
	}
	return mogoSlots, nil
}

func readTanxSlots(filePath string) (map[string]*TanxSlot, error) {
	excelFile, err := xlsx.OpenFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("%s 没有tanx.xlsx文件")
	}

	tanxSlots := make(map[string]*TanxSlot)
	for _, sheet := range excelFile.Sheets {
		for i, row := range sheet.Rows {
			if i == 0 || len(row.Cells) < 8 {
				continue
			}

			publishName, _ := row.Cells[0].String()
			slotName, _ := row.Cells[1].String()
			pid, _ := row.Cells[2].String()
			domain, _ := row.Cells[4].String()
			publishCat, _ := row.Cells[5].String()
			slotType, _ := row.Cells[7].String()
			devType, _ := row.Cells[8].String()

			pid = strings.TrimSpace(pid)
			if pid == "" {
				continue
			}

			tanxSlot := &TanxSlot{
				publishName: publishName,
				slotName:    slotName,
				pid:         pid,
				domain:      domain,
				publishCat:  publishCat,
				slotType:    getAdType(slotType),
				devType:     devType,
			}
			tanxSlots[pid] = tanxSlot
		}
	}
	return tanxSlots, nil
}

func getAdType(slotType string) string {
	switch {
	case strings.Contains(slotType, "横幅"),
		strings.Contains(slotType, "对联"),
		strings.Contains(slotType, "浮窗"),
		strings.Contains(slotType, "固定"),
		strings.Contains(slotType, "悬停"),
		strings.Contains(slotType, "折叠"):
		return "banner"
	case strings.Contains(slotType, "插屏"):
		return "插屏"
	case strings.Contains(slotType, "native"):
		return "feeds"
	}
	return "slotType"
}

func failedNotice(err error) {
	fmt.Printf("%s\n", err)
	fmt.Printf("5秒后自动退出\n")
	time.Sleep(5 * time.Second)
	os.Exit(0)
}
