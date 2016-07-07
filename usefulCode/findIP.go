package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
)

var filePath string
var reg *regexp.Regexp

func init() {
	flag.StringVar(&filePath, "log", "log", "log file")
	flag.Parse()

	reg = regexp.MustCompile(`\d+\.\d+\.\d+\.\d+`)
}

func findIp() {
	f, err := os.Open(filePath)
	if err != nil {
		panic("open file fail: " + err.Error())
	}
	defer f.Close()

	for sc := bufio.NewScanner(bufio.NewReader(f)); sc.Scan(); {
		ip := reg.FindString(sc.Text())
		if ip != "" {
			fmt.Printf("%s\n", ip)
		}

	}
}

func main() {
	findIp()
}
