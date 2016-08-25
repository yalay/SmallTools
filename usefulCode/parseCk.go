package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"
	"util/web"
)

var filePath string
var reg *regexp.Regexp

func init() {
	flag.StringVar(&filePath, "log", "log", "log file")
	flag.Parse()

	reg = regexp.MustCompile(`\/c(lk1|lk2|lk3|4|5)?\?\S+`)
}

func getInvalidCkUrls() []string {
	f, err := os.Open(filePath)
	if err != nil {
		panic("open file fail: " + err.Error())
	}
	defer f.Close()

	invalidCkUrls := make([]string, 0)
	for sc := bufio.NewScanner(bufio.NewReader(f)); sc.Scan(); {
		ckUrl := reg.FindString(sc.Text())
		if ckUrl == "" {
			continue
		}

		if isValidTargetUrl(ckUrl) {
			continue
		}
		invalidCkUrls = append(invalidCkUrls, ckUrl)
	}
	return invalidCkUrls
}

func isValidTargetUrl(ckUrl string) bool {
	// 显式target，但对应的value被decode为原始链接
	if strings.Contains(ckUrl, web.Ktarget+"="+web.KHttp) ||
		strings.Contains(ckUrl, web.Ktarget+"="+web.KHttps) {
		return false
	}
	return true
}

func getSlotId(ckUrl string) string {
	infoValues := parseInfo(ckUrl)
	if infoValues == nil {
		return ""
	}

	return infoValues.Get(web.KslotId)
}

func parseInfo(ckUrl string) url.Values {
	ckUrl = ckUrl[strings.Index(ckUrl, web.Kinfo):]
	ckValues, err := url.ParseQuery(ckUrl)
	if err != nil {
		fmt.Printf("parse info err:%v", err)
		return nil
	}

	info := ckValues.Get(web.Kinfo)
	var infoValues url.Values
	if ckValues.Get(web.KisPb) == "" && ckValues.Get(web.Ksid) != web.Adx {
		infoValues, err = web.ParseB64Query(info)
		if err != nil || infoValues.Get(web.Ksid) == "" {
			infoValues, err = web.ParsePbQuery(info)
			if err != nil || infoValues.Get(web.Ksid) == "" {
				fmt.Printf("parse info err2:%v", err)
				return nil
			}
		}
	} else {
		infoValues, err = web.ParsePbQuery(info)
		if err != nil || infoValues.Get(web.Ksid) == "" {
			infoValues, err = web.ParseB64Query(info)
			if err != nil || infoValues.Get(web.Ksid) == "" {
				fmt.Printf("parse info err3:%v", err)
				return nil
			}
		}
	}

	return infoValues
}

func main() {
	invalidCkUrls := getInvalidCkUrls()
	for _, ckUrl := range invalidCkUrls {
		slotId := getSlotId(ckUrl)
		fmt.Println(slotId)
	}
}
