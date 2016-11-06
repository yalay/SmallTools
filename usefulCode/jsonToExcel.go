package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/tealeg/xlsx"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type JsonData struct {
	Title        string
	Type         string
	Media_id     int
	AdFormat     []int
	AdForm       int
	Adspace_id   int
	Screen_level int
	Vertical     string
	Width        int
	Height       int
}

type JsonDatas struct {
	Datas []JsonData
}

var adFormatMap = map[int]string{
	1: "STATIC_PIC",
	2: "DYNAMIC_PIC",
	3: "SWF",
	4: "TXT",
}

var adFormMap = map[int]string{
	1: "web硬广",
	2: "app硬广",
	3: "开屏",
	4: "插屏",
	5: "信息流",
}

var adTypeMap = map[string]string{
	"1": "web",
	"2": "app-android",
	"3": "app-ios",
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

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("read err:%v\n", err)
		return
	}

	var jsonDatas = &JsonDatas{}
	err = json.Unmarshal(data, jsonDatas)
	if err != nil {
		fmt.Printf("json err:%v\n", err)
		return
	}

	xlsxfile := xlsx.NewFile()
	sheet, err := xlsxfile.AddSheet("Sheet1")
	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	writeExcel(sheet, jsonDatas)
	filePathPre := filePath[:(len(filePath) - 3)]
	err = xlsxfile.Save(filePathPre + "xlsx")
	if err != nil {
		fmt.Printf(err.Error())
	}
}

func writeExcel(sheet *xlsx.Sheet, jsonDatas *JsonDatas) {
	if sheet == nil || len(jsonDatas.Datas) == 0 {
		return
	}
	for _, data := range jsonDatas.Datas {
		row := sheet.AddRow()
		row.AddCell().Value = data.Title
		row.AddCell().Value = adTypeMap[data.Type]
		row.AddCell().Value = strconv.Itoa(data.Media_id)
		row.AddCell().Value = transAdFormat(data.AdFormat)
		row.AddCell().Value = adFormMap[data.AdForm]
		row.AddCell().Value = strconv.Itoa(data.Adspace_id)
		row.AddCell().Value = strconv.Itoa(data.Screen_level)
		row.AddCell().Value = data.Vertical
		row.AddCell().Value = strconv.Itoa(data.Width) + "x" + strconv.Itoa(data.Height)
	}
}

func transAdFormat(adFormats []int) string {
	result := ""
	for _, adFormat := range adFormats {
		result += adFormatMap[adFormat]
	}
	strings.Trim(result, ",")
	return result
}
