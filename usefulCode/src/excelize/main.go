package main

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"flag"

	_ "github.com/go-sql-driver/mysql"
	"github.com/xuri/excelize"
)

var catIds = map[string]int{
	"推薦": 1,
	"國際": 5,
	"財經": 7,
	"生活": 8,
	"娛樂": 9,
	"汽車": 11,
	"體育": 13,
	"健康": 14,
	"旅遊": 15,
	"科技": 16,
	"電玩": 18,
	"動漫": 19,
	"數碼": 20,
	"情感": 21,
	"女生": 22,
	"吃喝": 23,
	"港聞": 24,
	"萌寵": 25,
	"寫真": 26,
	"親子": 27,
	"社會": 28,
	"美圖": 29,
	"歷史": 30,
	"軍事": 31,
	"養生": 32,
	"文化": 33,
	"搞笑": 34,
	"電影": 35,
	"绅士": 36,
	"地產": 37,
	"小說": 38,
	"星座": 39,
	"心測": 40,
	"韓流": 41,
	"視頻": 42,
	"段子": 43,
	"藝術": 44,
	"动漫": 49,
	"爆笑": 50,
	"体育": 52,
	"娱乐": 53,
	"美食": 54,
}

var db *sql.DB

type NewsSource struct {
	Name       string
	WebId      int
	Website    string
	SourceName string
	ChannelId  int
}

var debug bool
var fileName string

func init() {
	flag.BoolVar(&debug, "d", false, "is debug")
	flag.StringVar(&fileName, "f", "./test.xlsx", "excel file")
	flag.Parse()

	var err error
	db, err = sql.Open("mysql", "root:@10.8.54.136/app_news?charset=utf8")
	if err != nil {
		panic(err)
	}
}

func (n *NewsSource) Insert() {
	webSql := fmt.Sprintf("select id from zyz_web_source where name like '%s' LIMIT 1", n.Name)
	rows, err := db.Query(webSql)
	if err != nil {
		fmt.Println(err)
		return
	}

	var webId int
	for rows.Next() {
		rows.Scan(&webId)
	}
	rows.Close()
	if webId == 0 {
		fmt.Println("not found. " + n.Website)
		return
	}

	sourceQuerySql := fmt.Sprintf("select id from zyz_article_source_grab where website='%s'", n.Website)
	rows2, err := db.Query(sourceQuerySql)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer rows2.Close()
	if rows2.Next() {
		fmt.Println(n.Website + "exist")
		return
	}

	sourceInsertSql := fmt.Sprintf("INSERT zyz_article_source_grab (webid, website, channelid) VAlUES (%d, '%s', %d)", webId, n.Website, n.ChannelId)
	if debug {
		fmt.Println(sourceInsertSql)
	} else {
		_, err := db.Exec(sourceInsertSql)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

}

func main() {
	xlsx, err := excelize.OpenFile(fileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	index := xlsx.GetSheetIndex("新增url分析结果")

	// Get all the rows in a sheet.
	var lastRows2 string
	var websites = []string{}
	var newss = make([]*NewsSource, 0)
	rows := xlsx.GetRows("sheet" + strconv.Itoa(index))
	for _, row := range rows {
		if len(rows) < 6 {
			continue
		}
		if row[2] != "" {
			lastRows2 = row[2]
			websites = append(websites, row[2])
		}

		if row[5] == "可以" {
			news := &NewsSource{
				Name:      lastRows2,
				Website:   row[3],
				ChannelId: catIds[row[4]],
			}
			newss = append(newss, news)
		}
	}

	insertSource(websites)
	for _, news := range newss {
		news.Insert()
	}
}

func insertSource(websites []string) {
	for _, website := range websites {
		if website == "" {
			continue
		}

		sourceSql := fmt.Sprintf("select name from zyz_web_source where name like '%s'", website)
		rows, err := db.Query(sourceSql)
		if err != nil {
			fmt.Println(err)
			return
		}
		var sourceName string
		for rows.Next() {
			rows.Scan(&sourceName)
		}
		rows.Close()

		if sourceName != "" {
			fmt.Println(sourceName + " exist")
			continue
		}

		sourceInsertSql := fmt.Sprintf("INSERT zyz_web_source (name) VAlUES ('%s')", website)
		if debug {
			fmt.Println(sourceInsertSql)
		} else {
			_, err := db.Exec(sourceInsertSql)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
