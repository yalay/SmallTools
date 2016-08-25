package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"
)

var (
	logFile = "log"
)

func init() {
	flag.StringVar(&logFile, "log", "./logFile", "log file")
	flag.Parse()
}

func main() {
	f, err := os.Open(logFile)
	if err != nil {
		panic("open log config fail: " + err.Error())
	}
	defer f.Close()

	for sc := bufio.NewScanner(bufio.NewReader(f)); sc.Scan(); {
		oriUrl := strings.TrimSpace(sc.Text())
		decodeUrl, _ := url.QueryUnescape(oriUrl)
		fmt.Println(decodeUrl)
	}
}
