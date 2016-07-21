package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	filePath = "./ip.conf"
)

type ipRecord struct {
	low  int
	high int
	area string
}

var ipRecords []ipRecord
var ipRecordsNum int
var ipArea string
var ipConf string

func init() {
	flag.StringVar(&ipConf, "ipConf", "./ip.conf", "ipConf file")
	flag.StringVar(&ipArea, "ipArea", "./ipArea.log", "ipArea file")
	flag.Parse()

	f, err := os.Open(ipConf)
	if err != nil {
		panic("open ip config fail: " + err.Error())
	}
	defer f.Close()

	for sc := bufio.NewScanner(bufio.NewReader(f)); sc.Scan(); {
		line := strings.TrimSuffix(sc.Text(), ";")
		sections := strings.Split(line, "\t")
		if len(sections) != 2 {
			fmt.Printf("format error :%s\n", line)
			continue
		}
		ips := strings.Split(sections[0], "-")
		if len(ips) != 2 {
			fmt.Printf("format error :%s\n", line)
			continue
		}
		low := IpToInt(ips[0])
		high := IpToInt(ips[1])
		record := ipRecord{low, high, sections[1]}
		if ipRecordsNum == 0 {
			ipRecords = append(ipRecords, record)
			ipRecordsNum++
			continue
		}
		preRecord := ipRecords[ipRecordsNum-1]
		if record.area == preRecord.area && preRecord.high+1 == record.low {
			ipRecords[ipRecordsNum-1].high = record.high
		} else {
			ipRecords = append(ipRecords, record)
			ipRecordsNum++
		}
	}
}

// Atoi convert string to int
func atoi(s string) int {
	i, _ := strconv.Atoi(strings.TrimSpace(s))
	return i
}

func IpToInt(ip string) int {
	ids := strings.Split(ip, ".")
	if len(ids) != 4 {
		fmt.Printf("converet ip to int error")
		return -1
	}
	ret := 0
	for i := range ids {
		ret = ret<<8 + atoi(ids[i])
	}
	return ret
}

func GetArea(ip string) string {
	ipInt := IpToInt(ip)
	for left, right := 0, ipRecordsNum-1; left <= right; {
		mid := (left + right) >> 1
		record := ipRecords[mid]
		if ipInt < record.low {
			right = mid - 1
		} else if ipInt > record.high {
			left = mid + 1
		} else {
			return record.area
		}
	}
	return ""
}

func main() {
	f, err := os.Open(ipArea)
	if err != nil {
		panic("open ip config fail: " + err.Error())
	}
	defer f.Close()

	for sc := bufio.NewScanner(bufio.NewReader(f)); sc.Scan(); {
		ipArea := strings.Fields(sc.Text())
		if len(ipArea) != 2 {
			fmt.Printf("ipArea format wrong:%s\n", sc.Text())
			continue
		}

		area := GetArea(ipArea[0])
		if !strings.HasPrefix(area, "CN") {
			fmt.Printf("ip:%s area:%s\n", ipArea[0], area)
		}
	}
}
