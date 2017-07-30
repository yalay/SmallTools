package main

import (
	"bufio"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/axgle/mahonia"
	"github.com/tealeg/xlsx"
)

type content struct {
	homeUrl  string
	title    string
	keywords string
	desc     string
	count    int
}

var reg, _ = regexp.Compile(`[0-9A-Za-z\p{Han}]+网`)

func getUrl() {
	f, err := os.Open("mm.log")
	if err != nil {
		panic("open file fail: " + err.Error())
	}
	defer f.Close()

	var slotContent = make(map[string]*content)
	for sc := bufio.NewScanner(bufio.NewReader(f)); sc.Scan(); {
		data := strings.Fields(sc.Text())
		if len(data) < 2 {
			continue
		}

		if content, ok := slotContent[data[0]]; ok {
			content.count++
			continue
		}

		if data[1] == "" {
			continue
		}

		website, err := url.QueryUnescape(data[1])
		if err != nil {
			fmt.Println(err)
			continue
		}

		if !strings.HasPrefix(website, "http") {
			continue
		}

		websiteUrl, err := url.Parse(website)
		if err != nil {
			fmt.Println(err)
			continue
		}

		slotContent[data[0]] = &content{
			homeUrl: websiteUrl.Scheme + "://" + websiteUrl.Host,
		}
	}

	var wg = sync.WaitGroup{}
	for _, curContent := range slotContent {
		go func(content *content) {
			wg.Add(1)
			newContent := parseUrl(content.homeUrl)
			if newContent != nil {
				content.title = newContent.title
				content.keywords = newContent.keywords
				content.desc = newContent.desc
			}
			wg.Done()
		}(curContent)
	}
	time.Sleep(time.Second)
	wg.Wait()
	save(slotContent)
}

func parseUrl(rawUrl string) *content {
	httpClient := http.Client{Timeout: 30 * time.Second}
	resp, err := httpClient.Get(rawUrl)
	if err != nil {
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil
	}

	document, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	var pageTitle, keywords, desc string
	var decode mahonia.Decoder
	var charset string
	document.Find("meta").Each(func(i int, s *goquery.Selection) {
		if html, ok := s.Attr("charset"); ok {
			lowerHtml := strings.ToLower(html)
			if strings.Contains(lowerHtml, "gb2312") || strings.Contains(lowerHtml, "gbk") {
				charset = "gbk"
			}
			return
		}

		if charsetHtml, ok := s.Attr("content"); ok {
			lowerHtml := strings.ToLower(charsetHtml)
			if strings.Contains(lowerHtml, "gb2312") || strings.Contains(lowerHtml, "gbk") {
				charset = "gbk"
				return
			}
		}
	})
	document.Find("meta").Each(func(i int, s *goquery.Selection) {
		if charset != "" {
			decode = mahonia.NewDecoder(charset)
		}

		if name, _ := s.Attr("name"); strings.EqualFold(name, "keywords") {
			keywords = s.AttrOr("content", "")
			if decode != nil {
				keywords = decode.ConvertString(keywords)
			}
			return
		}
		if name, _ := s.Attr("name"); strings.EqualFold(name, "description") {
			desc = s.AttrOr("content", "")
			if decode != nil {
				desc = decode.ConvertString(desc)
			}
			return
		}
	})
	pageTitle = document.Find("title").Text()
	if decode != nil {
		pageTitle = decode.ConvertString(pageTitle)
	}

	return &content{
		title:    parseTitle(pageTitle),
		keywords: keywords,
		desc:     desc,
	}
}

func save(slotContent map[string]*content) {
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Mediamax")
	if err != nil {
		fmt.Println(err.Error())
	}

	for slot, content := range slotContent {
		row := sheet.AddRow()
		row.AddCell().SetValue(slot)
		row.AddCell().SetValue(content.homeUrl)
		row.AddCell().SetValue(content.title)
		row.AddCell().SetValue(content.keywords)
		row.AddCell().SetValue(content.desc)
		row.AddCell().SetValue(content.count)
	}

	err = file.Save("mm.xlsx")
	if err != nil {
		fmt.Printf(err.Error())
	}
}

func parseTitle(oriText string) string {
	titles := strings.FieldsFunc(oriText, func(r rune) bool {
		return r == '-' || r == '_' || r == '|' || r == '—' ||
			r == '【' || r == '】' || r == ',' ||
			r == '，' || r == '–'
	})
	var shortestTitle string
	for i, title := range titles {
		title = strings.TrimSpace(title)
		if title == "" {
			continue
		}
		if strings.HasSuffix(title, "网") {
			return title
		}

		if shortestTitle == "" {
			shortestTitle = title
		} else {
			if i == len(titles)-1 && len(shortestTitle) > len(title) {
				shortestTitle = title
			}
		}
	}
	return shortestTitle
}

func main() {
	//fmt.Println(parseTitle("蜂鸟摄影论坛 – 极具人气的数码摄影论坛"))
	//fmt.Println(parseUrl("http://939636.cc"))
	getUrl()
}
