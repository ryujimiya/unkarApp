package untidy

import (
	"./search"
	"./util"
	"./get2ch"
	"bufio"
	"bytes"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"time"
)

const (
	PEPSTR       = `<em>$1</em>`
	SEARCH_TITLE = "unkar＠スレッド検索"
)

type SearchData struct {
	Title          string
	Canonical      bool
	CanonicalQuery string
	Description    string
	Word           string
	Type           string
	Order          string
	Boardname      string
	Dataflag       bool
	SearchMax      int
	Min            int
	Max            int
	Time           float64
	SelectBoard    string
	SearchGutter   string
	SearchMain     string
	Year           int
}

type selectItem struct {
	name  string
	board string
}

type Search struct {
	s               *unsearch.Search
	canonical       bool
	boardNameList   map[string]string
	boardSelectList []selectItem
	boardGroup      map[string]bool
	data            *unsearch.SearchData
	boardName       string
	title           string
	description     string
	canonicalQuery  string
}

type SearchFooter struct {
	back int
	next int
	now  int
	data []int
}

var HtmlSearchTempl = template.Must(template.ParseFiles("template/search.templ"))

func newSearch(r *http.Request) *Search {
	v := r.URL.Query()
	page, perr := strconv.Atoi(v.Get("p"))
	if perr != nil {
		page = 1
	}

	q := &unsearch.Query{
		QueryStr: v.Get("q"),
		Page:     page,
		Board:    v.Get("board"),
		Stype:    v.Get("type"),
		Order:    v.Get("order"),
	}
	s := &Search{
		boardNameList:   make(map[string]string),
		boardSelectList: make([]selectItem, 0, 4),
		boardGroup:      make(map[string]bool),
	}
	// オブジェクト生成
	s.s = unsearch.NewSearch(20, r, q)
	// 接続を閉じる
	defer s.Close()
	// 検索＆検索結果を取得
	s.data = s.s.Fetch() // タイムアウトする恐れあり

	s.createNameList()
	s.createGroup()
	s.createBoardName()
	s.createDesc()
	s.createCanonicalQuery(v)
	return s
}

func (s *Search) Close() {
	if s.s != nil {
		s.s.Close()
	}
}

func (s *Search) createNameList() {
	// 板の分布を取得
	rc := get2ch.NewGet2ch("", "").GetBBSmenu(false)
	scanner := bufio.NewScanner(rc)
	for scanner.Scan() {
		it := strings.Trim(scanner.Text(), " 　\t")
		line := strings.Split(it+"<>", "<>")
		if line[1] != "" {
			data := strings.Split(line[0], "/")
			s.boardNameList[data[1]] = line[1]
			s.boardSelectList = append(s.boardSelectList, selectItem{
				name:  line[1],
				board: data[1],
			})
		} else {
			s.boardSelectList = append(s.boardSelectList, selectItem{
				name: it,
			})
		}
	}
	rc.Close()
}

func (s *Search) createGroup() {
	for _, it := range s.data.Data {
		s.boardGroup[it.Board] = true
	}
}

func (s *Search) createBoardName() {
	board := s.s.GetBoard()
	if board != "" {
		s.boardName = s.boardNameList[board]
	}
}

func (s *Search) createDesc() {
	// デフォルトタイトルに設定
	s.title = SEARCH_TITLE
	s.description = SEARCH_TITLE

	word := s.s.GetWord()
	if word != "" {
		var result string
		var option string
		var index string

		if name, ok := s.boardNameList[s.s.GetBoard()]; ok {
			result = "に一致する" + name + "の検索結果"
		} else {
			result = "の検索結果"
		}
		if s.s.GetType() == "score" {
			option += " 単語数"
		}
		if s.s.GetOrder() == "asc" {
			if option == "" {
				option += " "
			}
			option += "昇順"
		}
		if page := s.s.GetPage(); page > 1 {
			index = " " + strconv.Itoa(page) + "ページ目"
		}
		// タイトル合成
		s.title = fmt.Sprintf("「%s」%s%s%s - %s", word, result, option, index, SEARCH_TITLE)

		if s.data.Data != nil {
			s.description = fmt.Sprintf("「%s」%s%s %d件中 %d～%d件目 - %s", word, result, option, s.data.Searchmax, s.data.Min, s.data.Max, SEARCH_TITLE)
		} else {
			s.description = fmt.Sprintf("「%s」%s%s - %s", word, result, option, SEARCH_TITLE)
		}
	}
}

func (s *Search) createCanonicalQuery(v url.Values) {
	var q string
	cq := s.getCanonicalQuery([]string{"query", "board", "type", "order", "page"})
	if v.Encode() != cq {
		// リクエストと異なる
		// リクエストも一度分解して再構築しているので、順番も揃うはず
		s.canonical = true
	}
	if s.s.GetWord() != "" {
		// 検索単語が設定されている
		q = "?" + cq
	}
	s.canonicalQuery = `http://unkar.org/search` + q
}

func (s *Search) printSelectBoard() string {
	buf := bytes.Buffer{}
	q := s.s.GetQuery()
	board := q.Board
	buf.WriteString(`<select name="board" id="search-select-board">` + "\n")
	buf.WriteString(`<optgroup label="説明">` + "\n")
	buf.WriteString(`<option value="">板の絞り込み</option>` + "\n")
	for _, it := range s.boardSelectList {
		if it.board == "" {
			buf.WriteString(`</optgroup>` + "\n" + `<optgroup label="` + it.name + `">` + "\n")
		} else {
			if board == it.board {
				buf.WriteString(`<option value="` + it.board + `" selected>` + it.name + `</option>` + "\n")
			} else {
				buf.WriteString(`<option value="` + it.board + `">` + it.name + `</option>` + "\n")
			}
		}
	}
	buf.WriteString("</optgroup>\n</select>\n")
	return buf.String()
}

func (s *Search) printGutter() string {
	buf := bytes.Buffer{}
	q := s.s.GetQuery()
	url := "/search?" + s.getCanonicalQuery([]string{"query", "type", "order"})
	if len(s.boardGroup) > 1 {
		buf.WriteString(`<div class="gutter-area"><h3>板で絞り込む</h3>`)
		buf.WriteString(`<ul id="now-search">`)
		for key, _ := range s.boardGroup {
			if it, ok := s.boardNameList[key]; ok {
				buf.WriteString(`<li><a href="` + url + `&board=` + key + `" rel="nofollow">` + it + `</a></li>`)
			}
		}
		buf.WriteString(`</ul></div>`)
	} else if s.s.GetBoard() != "" {
		buf.WriteString(`<div class="gutter-area"><h3>板で絞り込む</h3>`)
		buf.WriteString(`<ul id="now-search">`)
		buf.WriteString(`<li><a href="` + url + `" rel="nofollow">絞り込みを解除</a></li>`)
		buf.WriteString(`</ul></div>`)
	}
	buf.WriteString(`<div class="gutter-area">`)
	buf.WriteString(`<h3>並び替え条件</h3>`)
	buf.WriteString(`<ul>`)
	url = "/search?" + s.getCanonicalQuery([]string{"query", "board", "order"})
	if q.Stype != "score" {
		buf.WriteString(`<li><a href="` + url + `" rel="nofollow"><em>時間順</em></a></li>`)
		buf.WriteString(`<li><a href="` + url + `&type=score" rel="nofollow">単語数順</a></li>`)
	} else {
		buf.WriteString(`<li><a href="` + url + `" rel="nofollow">時間順</a></li>`)
		buf.WriteString(`<li><a href="` + url + `&type=score" rel="nofollow"><em>単語数順</em></a></li>`)
	}
	buf.WriteString(`</ul>`)
	buf.WriteString(`</div>`)
	buf.WriteString(`<div class="gutter-area">`)
	buf.WriteString(`<ul>`)
	url = "/search?" + s.getCanonicalQuery([]string{"query", "board", "type"})
	if q.Order != "asc" {
		buf.WriteString(`<li><a href="` + url + `" rel="nofollow"><em>降順↓</em></a></li>`)
		buf.WriteString(`<li><a href="` + url + `&order=asc" rel="nofollow">昇順↑</a></li>`)
	} else {
		buf.WriteString(`<li><a href="` + url + `" rel="nofollow">降順↓</a></li>`)
		buf.WriteString(`<li><a href="` + url + `&order=asc" rel="nofollow"><em>昇順↑</em></a></li>`)
	}
	buf.WriteString(`</ul>`)
	buf.WriteString(`</div>`)
	return buf.String()
}

func (s *Search) printMain() string {
	buf := bytes.Buffer{}
	if s.data.Data != nil {
		wordlist := unsearch.SplitSpace(s.s.GetWord())
		for i, it := range wordlist {
			wordlist[i] = regexp.QuoteMeta(it)
		}
		reg, err := regexp.Compile("(?i:(" + strings.Join(wordlist, "|") + "))")
		if err != nil {
			reg = unsearch.RegSpace
		}

		for _, it := range s.data.Data {
			var master string
			if it.Master != "" {
				master = it.Master
				if len([]rune(master)) >= 100 {
					master += "..."
				}
			} else {
				master = ""
			}
			url := "/r/" + it.Board + "/" + strconv.Itoa(it.Number)
			bbs := "/r/" + it.Board
			url_min := "/r/" + it.Board + "/" + strconv.Itoa(it.Number)
			buf.WriteString(`<article><div class="thread-area">`)
			buf.WriteString(`<span class="thread-title"><a href="` + url + `">` + reg.ReplaceAllString(it.Title, PEPSTR) + `</a></span>`)
			buf.WriteString(`<div class="thread-info">`)
			buf.WriteString(`<p><span class="path">` + url_min + `</span>`)
			if b, ok := s.boardNameList[it.Board]; ok {
				buf.WriteString(` - <a href="` + bbs + `" class="boardname">` + b + `</a>`)
			}
			if it.Resnum != 0 {
				buf.WriteString(` : <span class="resnum">` + strconv.Itoa(it.Resnum) + `レス</span>`)
			}
			buf.WriteString("</p>\n")
			buf.WriteString(`<p class="thread-info-description">` + unutil.CreateDateString(time.Unix(int64(it.Number), 0)) + ` - ` + reg.ReplaceAllString(master, PEPSTR) + "</p>\n")
			buf.WriteString("</div>\n</div>\n</article>\n")
		}
		if list := s.footerNavi(); list.data != nil {
			path := "/search?" + s.getCanonicalQuery([]string{"query", "board", "type", "order"})
			buf.WriteString("<section>\n<div id=\"page-nav\">\n<ul>\n")
			if list.back > 0 {
				buf.WriteString(`<li><a href="` + path + `&p=` + strconv.Itoa(list.back) + `">前</a></li>`)
			}
			for _, it := range list.data {
				sit := strconv.Itoa(it)
				if it == list.now {
					buf.WriteString(`<li class="page-now"><a href="` + path + `&p=` + sit + `">` + sit + `</a></li>`)
				} else {
					buf.WriteString(`<li><a href="` + path + `&p=` + sit + `">` + sit + `</a></li>`)
				}
			}
			if list.next > 0 {
				buf.WriteString(`<li><a href="` + path + `&p=` + strconv.Itoa(list.next) + `">次</a></li>`)
			}
			buf.WriteString("</ul>\n</div>\n</section>")
		}
	}
	return buf.String()
}

func (s *Search) footerNavi() (ret SearchFooter) {
	if s.data.Searchmax == 0 {
		return
	}
	q := s.s.GetQuery()
	pagemax := int(math.Ceil(float64(s.data.Searchmax) / float64(s.s.GetPageValue())))
	pagedata := pageIndex(1, q.Page, pagemax, 10)
	tmp := q.Page - 1
	if tmp > 0 {
		ret.back = tmp
	}
	tmp = q.Page + 1
	if tmp <= pagemax {
		ret.next = tmp
	}
	ret.now = q.Page
	ret.data = unutil.Range(pagedata[0], pagedata[1], 1)
	return
}

func (s *Search) getCanonicalQuery(list []string) string {
	v := url.Values{}
	q := s.s.GetQuery()
	for _, it := range list {
		switch it {
		case "query":
			if q.QueryStr != "" {
				v.Set("q", s.s.GetWord())
			} else {
				return ""
			}
		case "board":
			if q.Board != "" {
				v.Set("board", s.s.GetBoard())
			}
		case "type":
			if q.Stype == "score" {
				v.Set("type", "score")
			}
		case "order":
			if q.Order == "asc" {
				v.Set("order", "asc")
			}
		case "page":
			if q.Page > 1 {
				v.Set("p", strconv.Itoa(q.Page))
			}
		}
	}
	return v.Encode()
}

func pageIndex(min, now, max, hex int) [2]int {
	var tmp int
	if hex <= 1 {
		hex = 10
	}
	tmpmin := now - (hex / 2)
	tmpmax := now + int(math.Ceil(float64(hex)/float64(2)))
	if tmpmin >= min { // 最小と比べて同じか大きい
		if tmpmax <= max { // 最大と比べて同じか小さい
			min = tmpmin
			max = tmpmax
		} else { // あふれた場合
			tmp = tmpmax - max // 差を計算
			tmpmin -= tmp      // 下を詰める
			if tmpmin >= min { // 最小と比べて同じか大きい
				min = tmpmin
			}
		}
	} else {
		tmp = min - tmpmin // 差を計算
		tmpmax += tmp      // 上を広げる
		if tmpmax <= max { // 最大と比べて同じか小さい
			max = tmpmax
		}
	}
	return [2]int{min, max}
}

func SearchMain(r *http.Request) (out unutil.Output) {
	start := time.Now().UnixNano()
	search := newSearch(r)
	diff := float64(time.Now().UnixNano()-start) / float64(time.Second)
	buf := &bytes.Buffer{}
	HtmlSearchTempl.Execute(buf, &SearchData{
		Title:          search.title,
		Canonical:      search.canonical,
		CanonicalQuery: search.canonicalQuery,
		Description:    search.description,
		Word:           search.s.GetWord(),
		Type:           search.s.GetType(),
		Order:          search.s.GetOrder(),
		Boardname:      search.boardName,
		Dataflag:       search.data.SearchFlag,
		SearchMax:      search.data.Searchmax,
		Min:            search.data.Min,
		Max:            search.data.Max,
		Time:           diff,
		SelectBoard:    search.printSelectBoard(),
		SearchGutter:   search.printGutter(),
		SearchMain:     search.printMain(),
		Year:           time.Now().Year(),
	})

	// 戻り値に設定
	out.Code = http.StatusOK
	out.Header = http.Header{}
	out.Header.Set("Content-type", "text/html; charset=utf-8")
	out.Reader = buf
	out.ZFlag = true
	return
}
