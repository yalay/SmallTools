package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	logFile = "log"
)

func init() {
	flag.StringVar(&logFile, "log", "./logFile", "log file")
	flag.Parse()
}

func atoi(s string) int64 {
	i, _ := strconv.Atoi(strings.TrimSpace(s))
	return int64(i)
}

func main() {
	f, err := os.Open(logFile)
	if err != nil {
		panic("open log config fail: " + err.Error())
	}
	defer f.Close()

	for sc := bufio.NewScanner(bufio.NewReader(f)); sc.Scan(); {
		timestamp := atoi(strings.TrimSpace(sc.Text()))
		tm := time.Unix(timestamp, 0)
		fmt.Println(tm.Format("2006-01-02 15:04:05"))
	}
}
