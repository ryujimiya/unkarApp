package unsearch

/*
CREATE DATABASE unkar DEFAULT CHARACTER SET utf8;

CREATE TABLE thread_title(
id INT PRIMARY KEY AUTO_INCREMENT,
board VARCHAR(32),
number INT UNSIGNED,
title VARCHAR(100),
master VARCHAR(160),
resnum SMALLINT UNSIGNED,
INDEX board_index (board),
INDEX number_index (number),
FULLTEXT INDEX title_index (title) COMMENT 'parser "TokenMecab"'
) ENGINE = Mroonga DEFAULT CHARSET utf8;

EXPLAIN SELECT SQL_CALC_FOUND_ROWS board, number, title, master FROM thread_title WHERE MATCH(title) AGAINST('+あ' IN BOOLEAN MODE) ORDER BY number DESC LIMIT 0, 50;
EXPLAIN SELECT SQL_CALC_FOUND_ROWS board, number, title, master FROM thread_title WHERE MATCH(title) AGAINST('【PS3】Call of Duty:Black Ops' IN BOOLEAN MODE) ORDER BY number DESC;
EXPLAIN SELECT SQL_CALC_FOUND_ROWS board, number, title, master, MATCH(title) AGAINST('+け' IN BOOLEAN MODE) AS score FROM thread_title WHERE MATCH(title) AGAINST('+け' IN BOOLEAN MODE) ORDER BY score DESC;
SELECT SQL_CALC_FOUND_ROWS board, number, title FROM thread_title WHERE MATCH(title) AGAINST('ダウンロード 刑事罰' IN BOOLEAN MODE) ORDER BY number DESC;
SELECT SQL_CALC_FOUND_ROWS board, number, title, MATCH(title) AGAINST('+京都' IN BOOLEAN MODE) AS score FROM thread_title WHERE MATCH(title) AGAINST('+京都' IN BOOLEAN MODE) ORDER BY score DESC;
SELECT SQL_CALC_FOUND_ROWS board, number, title, MATCH(title) AGAINST('+京都' IN BOOLEAN MODE) AS score FROM thread_title WHERE MATCH(title) AGAINST('+京都' IN BOOLEAN MODE) ORDER BY score DESC;

ALTER TABLE thread_title ADD FULLTEXT INDEX (title, number);

SQLモードを変更して動作させる必要がある
SET @@GLOBAL.sql_mode='NO_ENGINE_SUBSTITUTION,STRICT_TRANS_TABLES,NO_BACKSLASH_ESCAPES';

*/

import (
	"../util"
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"time"
	//"os" //DEBUG
)

type Query struct {
	QueryStr string
	Page     int
	Board    string
	Stype    string
	Order    string
}

type Search struct {
	db               *sql.DB
	search_flag      bool
	searchMax        int
	param_page_value int
	ql               *Query
	req              *http.Request
}

type DBItem struct {
	Board  string `json:"Board"`
	Number int    `json:"Number"`
	Title  string `json:"Title"`
	Master string `json:"-"`
	Resnum int    `json:"Resnum"`
	Score  int    `json:"-"`
}

type SearchData struct {
	SearchFlag bool     `json:"SearchFlag"`
	Word       string   `json:"Word"`
	Max        int      `json:"Max"`
	Min        int      `json:"Min"`
	Searchmax  int      `json:"Searchmax"`
	Data       []DBItem `json:"Data"`
}

const (
	DB_USER                    = "unkar"
	DB_NAME                    = "unkar"
	DB_PASS                    = "unkokkounkar"
	DB_HOST                    = "localhost:3306"
	DB_TABLE                   = "thread_title"
	PARAM_PAGE_VALUE           = 100
	PARAM_CATEGORY_ALL         = "all"
	INSERT_MASTER_TEXT_LEN     = 160
	THREAD_NUMBER_MIN          = 927990000 // 1999/05/30
	THREAD_NUMBER_MAX_ADD_TIME = time.Hour * 24 * 60
	INSERT_CONN_SIZE           = 2
	INSERT_IDLE_SIZE           = 1
	UPDATE_CONN_SIZE           = 2
	UPDATE_IDLE_SIZE           = 1
	SEARCH_CONN_SIZE           = 2
	SEARCH_IDLE_SIZE           = 2
)

var query_kill_list = map[string]bool{
	"":     true,
	"-1\"": true,
}
var RegOption = regexp.MustCompile(`^([\-\+])?(.+)$`)
var RegSpace = regexp.MustCompile(`[　\s\t\(\)]+`)
var regTag = regexp.MustCompile("<\\/?[^>]*>")
var regUrl = regexp.MustCompile("(?:s?h?ttps?|sssp):\\/\\/[-_.!~*'()\\w;\\/?:\\@&=+\\$,%#\\|]+")

func NewSearch(pagevalue int, req *http.Request, q *Query) *Search {
	s := &Search{}
	s.param_page_value = unutil.MinInt(pagevalue, PARAM_PAGE_VALUE)
	s.req = req
	s.ql = q

	s.search_flag = s.checkParam()
	if s.search_flag {
		var err error
		s.db, err = connect(SEARCH_CONN_SIZE, SEARCH_IDLE_SIZE)
		s.search_flag = err == nil
	}
	return s
}

func (s *Search) Close() {
	if s.db != nil {
		s.db.Close()
	}
}

func (s *Search) GetSearchFlag() bool {
	return s.search_flag
}

func (s *Search) GetPageValue() int {
	return s.param_page_value
}

func (s *Search) GetQuery() *Query {
	return s.ql
}

func (s *Search) GetWord() string {
	return template.HTMLEscapeString(s.ql.QueryStr)
}

func (s *Search) GetPage() int {
	return s.ql.Page
}

func (s *Search) GetBoard() string {
	return template.HTMLEscapeString(s.ql.Board)
}

func (s *Search) GetType() string {
	return template.HTMLEscapeString(s.ql.Stype)
}

func (s *Search) GetOrder() string {
	return template.HTMLEscapeString(s.ql.Order)
}

func (s *Search) Fetch() *SearchData {
	var list []DBItem
	var searchmax int
	word := s.GetWord()
	retdata := &SearchData{
		SearchFlag: s.search_flag,
		Word:       word,
	}
	if s.search_flag {
		var err error
		list, searchmax, err = s.searchExec(s.createQuery())
		if err != nil {
			return retdata
		}
	} else {
		return retdata
	}
	p := s.ql.Page
	retdata.Max = unutil.MinInt(s.param_page_value*p, searchmax)
	retdata.Min = s.param_page_value*(p-1) + 1
	s.searchMax = searchmax
	retdata.Searchmax = s.searchMax
	retdata.Data = list
	return retdata
}

func (s *Search) checkParam() (retflag bool) {
	if s.ql == nil {
		s.ql = &Query{}
	}

	if _, ok := query_kill_list[s.ql.QueryStr]; ok {
		// 無効な文字列
		retflag = false
	} else {
		// 文字コードの判定が必要かも
		// 最近のページは全てUTF-8なので多分UTF-8だろう
		s.ql.QueryStr = s.ql.QueryStr
		retflag = true
	}
	if s.ql.Page <= 0 {
		s.ql.Page = 1
	}
	if s.ql.Board != "" {
		// 特になし
	}
	if s.ql.Stype == "score" {
		s.ql.Stype = "score"
	} else {
		s.ql.Stype = "number"
	}
	if s.ql.Order == "asc" {
		s.ql.Order = "asc"
	} else {
		s.ql.Order = "desc"
	}
	return
}

func (s *Search) createQuery() string {
	fmt.Printf("Search::createQuery\r\n") //DEBUG
	matchtext := s.getMatchFullText()
	boardtext := ""
	if s.ql.Board != "" {
		boardtext = "AND board = '" + sqlEscape(s.ql.Board) + "'"
	}
	page_start := (s.ql.Page - 1) * s.param_page_value

	sql := fmt.Sprintf(
		`SELECT SQL_CALC_FOUND_ROWS board, number, title, master, resnum, %s AS score FROM %s WHERE %s %s ORDER BY %s %s LIMIT %d, %d`,
		matchtext,
		DB_TABLE,
		matchtext,
		boardtext,
		s.ql.Stype,
		s.ql.Order,
		page_start,
		s.param_page_value)
	fmt.Print(sql) //DEBUG
	fmt.Printf("\r\n") //DEBUG
	return sql
}

func (s *Search) searchExec(sql string) ([]DBItem, int, error) {
	if s.db == nil {
		return nil, 0, errors.New("s.db nil")
	}
	rows, err := s.db.Query(sql)
	if rows == nil || err != nil {
		return nil, 0, errors.New("query fail")
	}

	list := make([]DBItem, 0, PARAM_PAGE_VALUE)
	for rows.Next() {
		item := DBItem{}
		rows.Scan(&item.Board, &item.Number, &item.Title, &item.Master, &item.Resnum, &item.Score)
		list = append(list, item)
	}
	if len(list) <= 0 {
		return nil, 0, errors.New("list <= 0")
	}

	var searchmax int
	if result, err := s.db.Query("SELECT FOUND_ROWS()"); err == nil {
		result.Next()
		err = result.Scan(&searchmax)
		if err != nil {
			searchmax = 0
		}
	}

	return list, searchmax, nil
}

func (s *Search) getMatchFullText() string {
	list := []string{}
	for _, value := range SplitSpace(s.ql.QueryStr) {
		if match := RegOption.FindStringSubmatch(value); match != nil {
			if match[1] == "" {
				match[1] = "+"
			}
			list = append(list, match[1]+sqlEscape(match[2]))
		}
	}
	return "MATCH(title) AGAINST('" + strings.Join(list, " ") + "' IN BOOLEAN MODE)"
}

type Insert struct {
	db *sql.DB
	bs int
	ql []string
}

var insertKillList = map[string]bool{
	"bbylive": true,
	"bbynews": true,
	"bbypink": true,
}

func NewInsert(bufsize int) *Insert {
	fmt.Printf("NewInsert\r\n") // DEBUG
	return &Insert{
		bs: bufsize,
		ql: make([]string, 0, bufsize),
	}
}

func (in *Insert) Push(it *DBItem) {
	fmt.Printf("Insert::Push\r\n") // DEBUG
	if it == nil {
		return
	}
	if _, ok := insertKillList[it.Board]; !ok {
		query := fmt.Sprintf(
			"('%s',%d,'%s','%s',%d)",
			it.Board,
			it.Number,
			sqlEscape(it.Title),
			sqlEscape(it.Master),
			it.Resnum)
		fmt.Print(query) // DEBUG
		fmt.Printf("\r\n") // DEBUG
		in.ql = append(in.ql, query)
		if len(in.ql) >= in.bs {
			in.Exec()
		}
	}
}

func (in *Insert) Exec() {
	fmt.Printf("Insert::Exec\r\n") // DEBUG
	if len(in.ql) > 0 {
		if in.db == nil {
			var err error
			in.db, err = connect(INSERT_CONN_SIZE, INSERT_IDLE_SIZE)
			if err != nil {
				in.db = nil
			}
		}
		if in.db != nil {
			query := "INSERT INTO " + DB_TABLE + " (board,number,title,master,resnum) VALUES" + strings.Join(in.ql, ", ")
			fmt.Print(query) // DEBUG
			fmt.Printf("\r\n") // DEBUG
			in.db.Exec(query)
		}
		in.ql = make([]string, 0, in.bs)
	}
}

type Update struct {
	db *sql.DB
}

func NewUpdate() *Update {
	fmt.Printf("NewUpdate\r\n") // DEBUG
	return &Update{}
}

func (u *Update) Update(board, thread string, resnum int) {
	fmt.Printf("Update::Update\r\n") // DEBUG
	if u.db == nil {
		var err error
		u.db, err = connect(UPDATE_CONN_SIZE, UPDATE_IDLE_SIZE)
		if err != nil {
			u.db = nil
		}
	}
	if u.db != nil {
		query := fmt.Sprintf(
			"UPDATE %s SET resnum=%d WHERE board='%s' AND number=%s",
			DB_TABLE,
			resnum,
			board,
			thread)
		fmt.Print(query) // DEBUG
		fmt.Printf("\r\n") // DEBUG
		u.db.Exec(query)
	}
}

func CreateDBItem(data []byte, board, thread string, linecount int) (item *DBItem) {
	i := bytes.IndexByte(data, '\n')
	if i <= 0 {
		return
	}
	resu := unutil.ShiftJISToUtf8String(string(data[:i]))
	res := strings.Split(resu, "<>")
	if len(res) > 4 {
		res[3] = regTag.ReplaceAllString(res[3], "") // tag
		res[3] = regUrl.ReplaceAllString(res[3], "") // url
		res[3] = strings.Trim(res[3], " 　")
		if linecount <= 0 {
			linecount = bytes.Count(data, []byte{'\n'})
		}
		max := time.Now().Add(THREAD_NUMBER_MAX_ADD_TIME).Unix()
		num, err := strconv.ParseInt(thread, 10, 64)
		if err == nil && num > THREAD_NUMBER_MIN && num <= max {
			item = &DBItem{
				Board:  board,
				Number: int(num),
				Title:  res[4],
				Master: unutil.Utf8Substr(res[3], INSERT_MASTER_TEXT_LEN),
				Resnum: linecount,
			}
		}
	}
	return
}

func SplitSpace(text string) []string {
	// 全角スペースを探索する
	text = RegSpace.ReplaceAllString(text, " ")
	return strings.SplitN(text, " ", 20)
}

func sqlEscape(s string) string {
	return strings.Replace(s, "'", "''", -1)
}

func connect(conn, idle int) (*sql.DB, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", DB_USER, DB_PASS, DB_HOST, DB_NAME))
	fmt.Printf("connect\r\n") //DEBUG
	fmt.Print(db) // DEBUG
	fmt.Printf("\r\n") // DEBUG
	fmt.Print(err) // DEBUG
	fmt.Printf("\r\n") // DEBUG
	if err == nil {
		db.SetMaxOpenConns(conn)
		db.SetMaxIdleConns(idle)
		fmt.Printf("success\r\n") // DEBUG
	}
	//os.Exit(-1) // DEBUG
	return db, err
}
