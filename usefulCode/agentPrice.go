package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
)

var filePath string

func init() {
	flag.StringVar(&filePath, "log", "log", "log file")
	flag.Parse()
}

func countPrice() {
	f, err := os.Open(filePath)
	if err != nil {
		panic("open file fail: " + err.Error())
	}
	defer f.Close()

	count := 0
	settlePriceTotal := 0.0
	aPriceTotal := 0.0
	for sc := bufio.NewScanner(bufio.NewReader(f)); sc.Scan(); {
		values, err := url.ParseQuery(sc.Text())
		if err != nil {
			fmt.Printf("ParseQuery %s err:%v", sc.Text(), err)
			break
		}
		aPriceStr := values.Get("aprice")
		settlePriceStr := values.Get("settlePrice")
		aPrice, err := Atof(aPriceStr)
		if err != nil {
			fmt.Printf("aPriceStr %s err:%v", aPriceStr, err)
			break
		}
		settlePrice, err := Atof(settlePriceStr)
		if err != nil {
			fmt.Printf("settlePriceStr %s err:%v", settlePriceStr, err)
			break
		}

		count++
		settlePriceTotal += settlePrice
		aPriceTotal += aPrice
	}
	fmt.Printf("count %d aPriceTotal:%f settlePriceTotal:%f", count, aPriceTotal, settlePriceTotal)
}

func Atof(s string) (float64, error) {
	return strconv.ParseFloat(strings.TrimSpace(s), 64)
}

func main() {
	countPrice()
}
