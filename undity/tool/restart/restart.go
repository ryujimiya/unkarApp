package main

import (
	"bufio"
	"bytes"
	"code.google.com/p/mahonia"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"
)

const DAT_MAX_SIZE = 512 * 1024

var RegsLine = regexp.MustCompile(`(\d+)\.dat<>(?:.*\s\((\d+)\))`)

func main() {
	path := "/2ch/dat"
	dir, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, it := range dir {
		if it.IsDir() {
			// ディレクトリの場合
			index(path + "/" + it.Name())
		}
	}
}

func index(path string) {
	sub, err := ioutil.ReadFile(path + "/subject.txt")
	if err != nil {
		return
	}
	count := 0
	linecount := 0
	now := time.Now().Unix()
	restart := time.Date(2013, time.January, 1, 0, 0, 0, 0, time.UTC)
	scanner := bufio.NewScanner(mahonia.NewDecoder("cp932").NewReader(bytes.NewReader(sub)))
	for scanner.Scan() {
		m := RegsLine.FindStringSubmatch(scanner.Text())
		if m == nil {
			continue
		}
		linecount++
		rescount, _ := strconv.ParseInt(m[2], 10, 64)
		l := len(m[1])
		if l >= 9 && l <= 10 && rescount < 1000 {
			p := path + "/" + m[1][:4] + "/" + m[1] + ".dat"
			stat, err := os.Stat(p)
			if err == nil {
				if (stat.Size() < DAT_MAX_SIZE) && (stat.ModTime().Unix() > now) {
					// 時間改変
					os.Chtimes(p, restart, restart)
					count++
				}
			}
		}
	}
	log.Printf("Line:%d\tCount:%d\tPath:%s", linecount, count, path)
}
