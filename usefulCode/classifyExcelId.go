package main

import (
	"flag"
	"fmt"
	"github.com/tealeg/xlsx"
	"os"
	"strings"
)

const (
	kEncTypeUnknown = "未知"
	kEncTypeOri     = "原始"
	kEncTypeMd5     = "MD5"
	kEncTypeSha1    = "SHA1"
)

const (
	kIdTypeUnknown   = "未知"
	kIdTypeIdfa      = "IDFA"
	kIdTypeImei      = "IMEI"
	kIdTypeAndroidId = "AndroidId"
)

var filePath = ""

var fileOpenMap = make(map[string]*os.File)

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

	xlFile, err := xlsx.OpenFile(filePath)
	if err != nil {
		fmt.Printf("open err:%v", err)
		return
	}

	filePathPre := filePath[:(len(filePath) - 4)]
	defer closeAllOpenTxt()
	for _, sheet := range xlFile.Sheets {
		for _, row := range sheet.Rows {
			if len(row.Cells) < 2 {
				continue
			}
			idTypeStr, _ := row.Cells[0].String()
			idValue, _ := row.Cells[1].String()

			idType, encType := getIdTypeAndEncType(idTypeStr, idValue)
			txtFileName := filePathPre + idType + "." + encType + ".txt"
			writeTxt(txtFileName, idValue)
		}
	}
}

func getIdTypeAndEncType(idTypeStr, idValue string) (string, string) {
	var idType = kIdTypeUnknown
	var encType = kEncTypeUnknown
	if len(idValue) == 32 {
		encType = kEncTypeMd5
	}

	if len(idValue) == 40 {
		encType = kEncTypeSha1
	}

	idTypeLower := strings.ToLower(idTypeStr)
	switch {
	case strings.Contains(idTypeLower, "idfa"):
		idType = kIdTypeIdfa
		if len(idValue) == 36 {
			encType = kEncTypeOri
		}
	case strings.Contains(idTypeLower, "android"):
		idType = kIdTypeAndroidId
		if len(idValue) == 16 {
			encType = kEncTypeOri
		}
	case strings.Contains(idTypeLower, "imei"):
		idType = kIdTypeImei
		if len(idValue) == 15 || len(idValue) == 14 {
			encType = kEncTypeOri
		}
	}
	return idType, encType
}

func writeTxt(fileName string, idValue string) {
	if fileName == "" {
		return
	}

	txtFile, ok := fileOpenMap[fileName]
	if !ok || txtFile == nil {
		txtFile, _ = os.Create(fileName)
		fileOpenMap[fileName] = txtFile
	}

	txtFile.WriteString(idValue + "\r\n")
}

func closeAllOpenTxt() {
	for _, txtFile := range fileOpenMap {
		if txtFile == nil {
			continue
		}
		txtFile.Close()
	}
}
