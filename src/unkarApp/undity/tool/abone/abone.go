package main

import (
	"bufio"
	"bytes"
	"code.google.com/p/mahonia"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	ROOT_DIR        = "/2ch/dat"
	ABONE_TEXT      = "unkarで削除<>あぼーん<>あぼーん<>あぼーん<>あぼーん"
	PARALLEL_DELETE = 16

	DB_USER  = "unkar"
	DB_NAME  = "unkar"
	DB_PASS  = "unkokkounkar"
	DB_HOST  = "133.242.1.134:4223"
	DB_TABLE = "thread_title"
)

type DBCmd struct {
	board  string
	number uint64
}

var regLineSplit = regexp.MustCompile(`\/(\w+)\/(\d{9,10})\/?$`)
var aboneText = utf8ToShiftJIS(ABONE_TEXT)
var aboneTime = time.Now().Add(time.Hour * 24 * 365 * 5).UTC()
var dbCommand chan<- *DBCmd

func main() {
	var endch <-chan bool
	dbCommand, endch = delDBProc()
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		analyzeLine(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "なんかエラー:", err)
	}

	dbCommand <- nil
	<-endch
}

func analyzeLine(line string) {
	list := strings.Split(line, " ")
	if m := regLineSplit.FindStringSubmatch(list[0]); m != nil {
		path := ROOT_DIR + "/" + m[1] + "/" + m[2][:4] + "/" + m[2] + ".dat"
		if len(list) > 1 {
			reslist := []int{}
			res := strings.Split(list[1], ",")
			for _, it := range res {
				num, err := strconv.Atoi(it)
				if err == nil && num > 0 && num < 1010 {
					// 正しく変換できた場合
					reslist = append(reslist, num)
				}
			}
			delRes(path, reslist)
		} else {
			delThread(path, m[1], m[2])
		}
	}
}

func delRes(path string, reslist []int) {
	if data, err := ioutil.ReadFile(path); err == nil {
		list := strings.Split(string(data), "\n")
		l := len(list)
		for _, it := range reslist {
			if it > 0 && l >= it {
				list[it-1] = aboneText
			}
		}
		err := ioutil.WriteFile(path, []byte(strings.Join(list, "\n")), 0666)
		if err == nil {
			os.Chtimes(path, aboneTime, aboneTime)
			oklog(path)
		} else {
			errlog(path)
		}
	} else {
		errlog(path)
	}
}

func delThread(path, board, number string) {
	err := os.Remove(path)
	if err == nil {
		num, err := strconv.ParseUint(number, 10, 64)
		if err == nil {
			// DB削除
			dbCommand <- &DBCmd{
				board:  board,
				number: num,
			}
		}
		oklog(path)
	} else {
		errlog(path)
	}
}

func oklog(msg string) {
	fmt.Fprintln(os.Stdout, "ok - "+msg)
}

func errlog(msg string) {
	fmt.Fprintln(os.Stdout, "miss - "+msg)
}

func utf8ToShiftJIS(data string) string {
	buf := bytes.Buffer{}
	enc := mahonia.NewEncoder("cp932")
	enc.NewWriter(&buf).Write([]byte(data))
	return buf.String()
}

func sqlEscape(s string) string {
	return strings.Replace(s, "'", "''", -1)
}

func delDBProc() (chan<- *DBCmd, <-chan bool) {
	ch := make(chan *DBCmd, 16)
	endch := make(chan bool)
	go func(ch <-chan *DBCmd, endch chan<- bool) {
		db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", DB_USER, DB_PASS, DB_HOST, DB_NAME))
		if err != nil {
			panic(err)
		}
		cmdlist := []*DBCmd{}
		for cmd := range ch {
			if cmd != nil {
				cmdlist = append(cmdlist, cmd)
				if len(cmdlist) > PARALLEL_DELETE {
					exeQuery(db, cmdlist)
					cmdlist = []*DBCmd{}
				}
			} else {
				// 処理を抜ける
				break
			}
		}
		exeQuery(db, cmdlist)
		db.Close()
		endch <- true
	}(ch, endch)
	return ch, endch
}

func exeQuery(db *sql.DB, cmdlist []*DBCmd) {
	if db != nil && len(cmdlist) > 0 {
		qlist := []string{}
		for _, cmd := range cmdlist {
			qlist = append(qlist, fmt.Sprintf(`(board='%s' AND number=%d)`, sqlEscape(cmd.board), cmd.number))
		}
		query := fmt.Sprintf(`DELETE FROM %s WHERE %s`, DB_TABLE, strings.Join(qlist, " OR "))
		res, err := db.Exec(query)
		if err != nil {
			return
		}
		num, _ := res.RowsAffected()
		fmt.Fprintln(os.Stdout, "db delete - ", num)
	}
}
