package untidy

import (
	"./conf"
	"./get2ch"
	"./search"
	"./util"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

const (
	EXPIRES_HYSTERESIS_TIME = 24 * 2 * time.Hour
)

type ConvSearch struct {
	search *unsearch.Search
}

func NewConvSearch(word string, pageadd, page int, r *http.Request) *ConvSearch {
	q := &unsearch.Query{
		QueryStr: word,
		Page:     page,
	}
	cs := &ConvSearch{
		search: unsearch.NewSearch(pageadd, r, q),
	}
	return cs
}

func (cs *ConvSearch) Close() {
	if cs.search != nil {
		cs.search.Close()
	}
}

func (cs *ConvSearch) pageCreate() []byte {
	data, err := json.Marshal(cs.search.Fetch())
	if err != nil {
		data = []byte(`{}`)
	}
	return data
}

type Conv struct {
	mod        time.Time
	print_data bytes.Buffer
	code404    bool
	path       string
	r          *http.Request
	Server     string
	Code       int
	Header     http.Header
}

var RegThread = regexp.MustCompile(`^/(?:\w+\.(?:2ch\.net|bbspink\.com)/)?(\w+)/(\d{9,10})(.+)?`)
var RegSearch = regexp.MustCompile(`^/\*/search/(.*)/(\d+)/(\d+)$`)
var RegBoard = regexp.MustCompile(`^/(?:\w+\.(?:2ch\.net|bbspink\.com)/)?(\w+)`)
var RegRange = regexp.MustCompile(`^bytes=(\d+)-$`)
var dummyDat = unutil.Utf8ToShiftJIS([]byte(`以下、unkarがお送りします<><>2038/01/19(火) 03:19:08.32 ID:UnKouNkO0<> は？ <>403 - Forbidden　←は？
以下、unkarがお送りします<><>2038/01/19(火) 03:19:54.87 ID:uNkOUnTi0<> <a href="../test/read.cgi/poverty/9246666279/1">&gt;&gt;1</a> <br> 4xx は、お前が言ってることとかがおかしいエラー <br> 400 お前が言ってることが俺は理解できない <br> 401 お前は認証しないと見れない <br> 403 お前が言ったことは理解したし認証もしたが、俺は拒否する。 <br> 404 お前が欲しがってるものは、無い。 <br>  <br> 5xx は、俺がおかしくてお前に返せないエラー。 <br> 500 なんか、処理してたらエラー起こった。返せない。 <br> 503 ちょっと忙しいorメンテナンスしてるんで、返せない。 <>
以下、unkarがお送りします<><>2038/01/19(火) 03:19:61.32 ID:UnTIUnKo0<> このページは403のつもりです。 <br> http://www.amazon.co.jp/gp/bestsellers/dvd/562020/ref=zg_bs_nav_d_1_d?tag=unkar-22&tag=unkokko--22 <>
`))

func NewConv(path string, r *http.Request) *Conv {
	cv := &Conv{
		path:   path,
		r:      r,
		Code:   http.StatusOK,
		Header: http.Header{},
	}
	if match := RegThread.FindStringSubmatch(path); match != nil {
		// スレッド取得＆表示
		if match[3] != "" {
			cv.dat_main_dummy()
		} else {
			cv.dat_main(match[1], match[2])
		}
	} else if path == "/server" {
		// 板一覧取得＆表示
		cv.ita_print()
	} else if path == "/*/now" {
		// 直近のアクセス
		cv.now_access(100)
	} else if match := RegSearch.FindStringSubmatch(path); match != nil {
		cv.search(r, match)
	} else if match := RegBoard.FindStringSubmatch(path); match != nil {
		// スレッド一覧取得＆表示
		cv.sure_main(match[1])
	} else {
		// 板一覧取得＆表示
		cv.ita_print()
	}
	return cv
}

func (cv *Conv) GetMod() time.Time {
	return cv.mod
}

func (cv *Conv) GetData() []byte {
	return cv.print_data.Bytes()
}

func (cv *Conv) Is404() bool {
	return cv.code404
}

// 板一覧表示
func (cv *Conv) ita_print() {
	g := get2ch.NewGet2ch("", "")
	rc := g.GetBBSmenu(true)
	cv.print_data.ReadFrom(rc)
	rc.Close()
	cv.mod = g.GetModified()
	unutil.CheckNotModified(cv.r, cv.mod)
	// ヘッダー送信
	cv.Header.Set("Content-type", "text/plain; charset=utf-8")
}

// 直近のアクセス
func (cv *Conv) now_access(num int) {
	data := bytes.Buffer{}
	list := get2ch.GetViewThreadList(0, num)
	for _, it := range list {
		url := it.Board + "/" + it.Thread
		data.WriteString(fmt.Sprintf("%s<>%s<>%d<>%d<>%s<>%d", url, it.Title, it.Lastdate.Unix(), it.Res, it.Boardname, it.Addtime))
		data.WriteString("\n")
	}
	cv.print_data.ReadFrom(&data)
	// ヘッダー送信
	cv.Header.Set("Content-type", "text/plain; charset=utf-8")
}

// 直近のアクセス
func (cv *Conv) search(r *http.Request, match []string) {
	num, _ := strconv.Atoi(match[2])
	page, _ := strconv.Atoi(match[3])
	word := match[1]
	cs := NewConvSearch(word, num, page, r)
	cv.print_data.Write(cs.pageCreate())
	cs.Close()
	// ヘッダー送信
	cv.Header.Set("Content-type", "text/plain; charset=utf-8")
}

// スレッド一覧取得
func (cv *Conv) sure_main(board string) {
	var title string
	// データの取得
	g := get2ch.NewGet2ch(board, "")
	q := cv.r.URL.Query()
	cv.Server = g.GetServer("")
	if title = g.GetBoardName(); title == "" {
		title = "不明な板名"
	}
	if q.Get("name") != "" {
		cv.print_data.WriteString(cv.path[1:] + "\n" + title + "\n" + cv.Server + "\n")
	} else {
		cv.print_data.ReadFrom(unutil.ShiftJISToUtf8Reader(bytes.NewReader(g.GetData())))
		cv.print_data.WriteString(title + "\n")
	}
	cv.mod = g.GetModified()
	unutil.CheckNotModified(cv.r, cv.mod)
	// ヘッダー送信
	cv.Header.Set("Content-type", "text/plain; charset=utf-8")
}

// .datファイル解析
func (cv *Conv) dat_main(board, thread_number string) {
	// データの取得
	g := get2ch.NewGet2ch(board, thread_number)
	data := g.GetData()
	cv.mod = g.GetModified()
	cv.code404 = g.Is404()
	cv.Server = g.GetServer("")
	// キャッシュの確認
	unutil.CheckNotModified(cv.r, cv.mod)
	encoding := "utf-8"
	q := cv.r.URL.Query()
	name := q.Get("name") != ""

	if cv.code404 {
		cv.mod = time.Time{}
		if name {
			// UTF-8
			cv.print_data.WriteString(cv.path[1:] + "\n不明\n" + cv.Server + "\n")
		} else {
			data = []byte("unkar.org<><>" + unutil.CreateDateString(time.Now()) + "<>スレッドを発見できませんでした。dat落ちのようです。<>dat落ち\n")
			if q.Get("charset") == "utf8" {
				// UTF-8
				cv.print_data.Write(data)
			} else {
				// Shift_JIS
				cv.print_data.Write(unutil.Utf8ToShiftJIS(data))
				encoding = "Shift_JIS"
			}
		}
	} else if name {
		var line []byte
		var title string
		index := bytes.IndexByte(data, '\n')
		if index >= 0 {
			line = data[:index:index]
		} else {
			line = []byte{}
		}
		resu := bytes.Split(unutil.StripTags(unutil.ShiftJISToUtf8(line), get2ch.HtmlTag), []byte("<>"))
		if len(resu) > 4 {
			title = string(resu[4])
		} else {
			title = "不明"
		}
		// UTF-8
		cv.print_data.WriteString(cv.path[1:] + "\n" + title + "\n" + cv.Server + "\n")
	} else {
		if q.Get("charset") == "utf8" {
			// UTF-8
			cv.print_data.ReadFrom(unutil.ShiftJISToUtf8Reader(bytes.NewReader(data)))
		} else {
			// Shift_JIS
			cv.print_data.Write(data)
			encoding = "Shift_JIS"
		}
	}
	// ヘッダー送信
	cv.Header.Set("Content-type", "text/plain; charset="+encoding)
}

// .datファイル解析
func (cv *Conv) dat_main_dummy() {
	// データの取得
	cv.mod = time.Now()
	cv.code404 = false
	cv.Server = "unkar.org"
	// キャッシュの確認
	unutil.CheckNotModified(cv.r, cv.mod)
	// Shift_JIS
	cv.print_data.Write(dummyDat)
	// ヘッダー送信
	cv.Header.Set("Content-type", "text/plain; charset=Shift_JIS")
}

func ConvMain(path string, r *http.Request) (out unutil.Output) {
	c := NewConv(path, r)
	data := c.GetData()
	mod := c.GetMod()
	code404 := c.Is404()

	if mod.IsZero() == false && code404 == false {
		req_time := time.Now()
		fs := len(data)
		// データ長の設定
		c.Header.Set("X-Unkar-Length", strconv.Itoa(fs))
		// サーバーが判明している
		if c.Server != "" {
			// 日本語名のサーバが出来たら問題かも
			c.Header.Set("X-2ch-Server", c.Server)
		}
		// ETagを設定
		c.Header.Set("ETag", unutil.CreateETag(mod))
		if req_time.Before(mod.Add(-1 * EXPIRES_HYSTERESIS_TIME)) {
			// キャッシュの有効期限を送付
			c.Header.Set("Expires", unutil.CreateModString(req_time.Add(unconf.OneYearSec)))
		} else {
			// 最終更新時刻送付
			c.Header.Set("Last-Modified", unutil.CreateModString(mod))
			ra := r.Header.Get("Range")
			if ra != "" {
				if match := RegRange.FindStringSubmatch(ra); match != nil {
					start, _ := strconv.Atoi(match[1])
					if start < fs {
						c.Code = http.StatusPartialContent
						c.Header.Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, fs-1, fs))
						data = data[start:]
					} else {
						// 範囲の指定が変
						c.Code = http.StatusRequestedRangeNotSatisfiable
						data = []byte{}
					}
				}
			}
		}
	}

	// 戻り値に設定
	if code404 {
		out.Code = http.StatusNotFound
	} else {
		out.Code = c.Code
	}
	out.Header = c.Header
	out.Reader = bytes.NewReader(data)
	out.ZFlag = len(data) > 256
	return
}
