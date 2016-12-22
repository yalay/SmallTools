package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	sjson "github.com/bitly/go-simplejson"
)

var fileName, keys string

func init() {
	flag.StringVar(&fileName, "f", "", "info file")
	flag.StringVar(&keys, "k", "", "ip,rip")
	flag.Parse()
}

func main() {
	f, err := os.Open(fileName)
	if err != nil {
		panic("open ip config fail: " + err.Error())
	}
	defer f.Close()

	for sc := bufio.NewScanner(bufio.NewReader(f)); sc.Scan(); {
		fmt.Println(getJsonValues(sc.Text(), strings.Split(keys, ",")))
	}
}

func getJsonValues(info string, keys []string) string {
	result := ""
	simpleJson, err := sjson.NewJson([]byte(info))
	if err != nil {
		fmt.Println(err.Error())
	}
	for _, key := range keys {
		value := simpleJson.Get(key).Interface()
		valueStr := fmt.Sprintf("%v", value)
		if result == "" {
			result = valueStr
		} else {
			result += " " + valueStr
		}
	}
	return result
}
