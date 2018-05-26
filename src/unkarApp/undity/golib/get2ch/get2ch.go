package get2ch

import (
	//"../conf" // unconf.Ver
	"../util"
	"./process"
	"bufio"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt" // DEBUG
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	CONF_FOLDER      = "./2ch/dat"    // dat保管フォルダ名
	CONF_ITAURL_HOST = "menu.2ch.net" // 板情報取得URL
	CONF_ITAURL_FILE = "bbsmenu.html"
	BOURBON_HOST     = "bg20.2ch.net" // 2chキャッシュサーバ
	BOURBON_KEY      = "bourbon"
	PARALLEL_LIMIT   = 9 // 同時アクセス制限

	DAT_CACHE_TIME_THREAD = 20 * time.Second
	DAT_CACHE_TIME_BOARD  = 7 * time.Second
	SETTING_CACHE_TIME    = 24 * 7 * time.Hour
	//DAT_MAX_SIZE              = 614400
	DAT_MAX_SIZE              = 0x7FFFFFFF
	DAT_NOT_REQUEST_ADD_MOD   = 24 * 365 * 5 * time.Hour // 5年間
	DAT_NOT_REQUEST_WAIT      = 24 * time.Hour           // 1日
	DAT_NOT_REQUEST_RES_COUNT = 1000                     // 1000レスを超えていたらリクエストしないようにする
	DAT_NOT_SIZE_LIMIT        = 524288                   // しきい値
	DAT_NOT_RENEW_SEC         = 24 * 61 * time.Hour      // 2ヶ月

	FILE_SUBJECT_TXT     = "subject.txt"
	FILE_SUBJECT_TXT_REQ = "subject.txt"
	FILE_SETTING_TXT     = "setting.txt"
	FILE_SETTING_TXT_REQ = "SETTING.TXT"

	TIMEOUT_SEC = 12 * time.Second

	//USER_AGENT = "Monazilla/1.00 (unkar/" + unconf.Ver + ")"
	// Chrome
	USER_AGENT = "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2272.89 Safari/537.36"
)

const (
	DAT_CREATE = iota
	DAT_APPEND
	DAT_BOURBON_THREAD
	DAT_BOURBON_BOARD
)

type CacheState interface {
	Size() int64
	Amod() time.Time
	Mmod() time.Time
}

type Cache interface {
	Path(s, b, t string) string
	GetData(s, b, t string) ([]byte, error)
	GetDataRC(s, b, t string) (io.ReadCloser, error)
	SetData(s, b, t string, d []byte) error
	SetDataAppend(s, b, t string, d []byte) error
	SetMod(s, b, t string, m, a time.Time) error
	Exists(s, b, t string) bool
	Stat(s, b, t string) (CacheState, error)
}

type Get2ch struct {
	size      int64     // datのデータサイズ
	mod       time.Time // datの最終更新時間
	cache_mod time.Time // datの最終更新時間
	code      int       // HTTPステータスコード
	err       error     // エラーメッセージ
	code404   bool      // code404フラグ
	server    string
	board     string
	thread    string
	req_time  time.Time
	cache     Cache
	bourbon   bool // バーボンフラグ
	numlines  int  // 行数
}

var catekill = map[string]bool{
	"特別企画":        true,
	"チャット":        true,
	"他のサイト":       true,
	"まちＢＢＳ":       true,
	"ツール類":        true,
	"チャット２ｃｈ＠ＩＲＣ": true,
	"Top10":       true,
	"2chのゴミ箱":     true,
	"BBSPINKのゴミ箱": true,
}

var sabakill = map[string]bool{
	"www.2ch.net":         true,
	"info.2ch.net":        true,
	"find.2ch.net":        true,
	"v.isp.2ch.net":       true,
	"m.2ch.net":           true,
	"test.up.bbspink.com": true,
	"stats.2ch.net":       true,
	"c-au.2ch.net":        true,
	"c-others1.2ch.net":   true,
	"movie.2ch.net":       true,
	"img.2ch.net":         true,
	"ipv6.2ch.net":        true,
	"be.2ch.net":          true,
	"p2.2ch.net":          true,
	"shop.2ch.net":        true,
	"watch.2ch.net":       true,
}

type hideData struct {
	server string
	name   string
}

var hideboard = map[string]hideData{
	"sakhalin": hideData{
		server: "toro.2ch.net",
		name:   "2ch開発室＠2ch掲示板",
	},
}

var RegServerItem = regexp.MustCompile(`<B>([^<]+)<\/B>`)
var RegServer = regexp.MustCompile(`<A HREF=http:\/\/([^\/]+)\/([^\/]+)\/>([^<]+)<\/A>`)

var HtmlTag = []string{"br", "font", "b"}

var boardServerObj *process.BoardServerBox
var boardNameObj *process.BoardNameBox
var bbnCacheObj *process.BBNCacheBox
var viewThreadListObj *process.ViewThreadBox
var dbInsertThread *process.DBInsertBox
var dbUpdateThread *process.DBUpdateBox
var boardThreadMap *process.BoardThreadMap
var parallelRequestLimitCh chan struct{}

var tanpanmanNagoyaee = unutil.Utf8ToShiftJIS([]byte("短パンマン ★名古屋はエ～エ～で"))
var tanpanman = unutil.Utf8ToShiftJIS([]byte("短パンマン ★"))
var nagoyaee = unutil.Utf8ToShiftJIS([]byte("名古屋はエ～エ～で"))

func init() {
	// サーバリスト更新
	boardServerObj = process.NewBoardServerBox(setServerList)
	boardNameObj = process.NewBoardNameBox()
	bbnCacheObj = process.NewBBNCacheBox()
	viewThreadListObj = process.NewViewThreadBox()
	dbInsertThread = process.NewDBInsertBox()
	dbUpdateThread = process.NewDBUpdateBox()
	boardThreadMap = process.NewBoardThreadMap()
	parallelRequestLimitCh = make(chan struct{}, PARALLEL_LIMIT)
}

func GetViewThreadList(start, end int) []process.ViewThreadItem {
	if start < 0 || start >= end || end > process.VIEW_THREAD_LIST_SIZE {
		return []process.ViewThreadItem{}
	}
	return viewThreadListObj.GetThreadList(start, end)
}

func GetViewThreadLot() ([]process.ViewThreadItem, uint64) {
	return viewThreadListObj.GetThreadLot()
}

func NewGet2ch(board, thread string) *Get2ch {
	g2ch := &Get2ch{
		size:      0,
		mod:       time.Time{},
		cache_mod: time.Time{},
		code:      0,
		err:       nil,
		code404:   false,
		server:    "",
		board:     "",
		thread:    "",
		req_time:  time.Now(),
		cache:     NewFileCache(CONF_FOLDER),
		bourbon:   false, // バーボンフラグ
		numlines:  0,
	}
	g2ch.server = g2ch.GetServer(board)
	g2ch.board = board
	if _, err := strconv.ParseInt(thread, 10, 64); err == nil {
		g2ch.thread = thread
	}
	return g2ch
}

func (g2ch *Get2ch) GetData() (data []byte) {
	// 初期化
	g2ch.size = 0
	g2ch.mod = time.Time{}
	g2ch.cache_mod = time.Time{}
	g2ch.code = 0
	g2ch.err = nil
	g2ch.code404 = false
	// 現在のバーボン状態を取得
	g2ch.bourbon = getBourbonCache()
	g2ch.numlines = 0

	if g2ch.isThread() {
		// スレッドの場合
		if st, err := g2ch.cache.Stat(g2ch.server, g2ch.board, g2ch.thread); err == nil {
			// 現状取得できているスレッド
			mod := st.Mmod()
			if g2ch.req_time.Before(mod) || g2ch.boardThreadLookup() == false {
				// 更新時間が未来のログだった場合
				// もしくは、スレッド一覧に記載が無い場合
				data, _ = g2ch.cache.GetData(g2ch.server, g2ch.board, g2ch.thread)
				g2ch.size = st.Size()
				g2ch.mod = mod
				g2ch.cache_mod = mod
				g2ch.code = 302
				if g2ch.size < DAT_MAX_SIZE {
					g2ch.addThreadCache(data)
				} else {
					data = g2ch.dataError()
				}
			} else {
				// 通常取得
				if g2ch.bourbon {
					data = g2ch.bourbonData()
				} else {
					data = g2ch.normalData(true)
				}
			}
		} else {
			// 2015-3-17 AngeloのスレッドがNot Foundになる
			//   dat落ちはなくなったはずなので、常に通常取得するようにした
			//// 現状取得できていないスレッドで、
			//// 一定以上前に立てられたものならば取得しない
			//mod_limit := g2ch.req_time.Add(-1 * DAT_NOT_RENEW_SEC)
			//tnumber, _ := strconv.ParseInt(g2ch.thread, 10, 64)
			//if time.Unix(tnumber, 0).After(mod_limit) && g2ch.boardThreadLookup() {
			// 期間内に立てられている
			// また、スレッド一覧に記載がある場合
			// 通常取得
			if g2ch.bourbon {
				data = g2ch.bourbonData()
			} else {
				data = g2ch.normalData(true)
			}
			//} else {
			//	data = g2ch.dataErrorDat()
			//	if g2ch.bourbon == false {
			//		// バーボン状態ではない場合
			//		// 404とする
			//		g2ch.code404 = true
			//	}
			//}
		}
	} else {
		// 通常取得
		if g2ch.bourbon {
			data = g2ch.bourbonData()
		} else {
			data = g2ch.normalData(true)
		}
	}

	// SJIS-winで返す
	return
}

func (g2ch *Get2ch) GetByteSize() int64 {
	return g2ch.size
}

func (g2ch *Get2ch) GetModified() time.Time {
	return g2ch.mod
}

func (g2ch *Get2ch) GetHttpCode() int {
	return g2ch.code
}

func (g2ch *Get2ch) GetError() error {
	return g2ch.err
}

func (g2ch *Get2ch) Is404() bool {
	return g2ch.code404
}

func (g2ch *Get2ch) NumLines(data []byte) int {
	if g2ch.numlines == 0 {
		g2ch.numlines = bytes.Count(data, []byte{'\n'})
	}
	return g2ch.numlines
}

func (g2ch *Get2ch) isThread() bool {
	return g2ch.server != "" && g2ch.board != "" && g2ch.thread != ""
}

func (g2ch *Get2ch) isBoard() bool {
	return g2ch.server != "" && g2ch.board != "" && g2ch.thread == ""
}

func dialTimeout(network, addr string) (net.Conn, error) {
	con, err := net.DialTimeout(network, addr, TIMEOUT_SEC)
	if err == nil {
		con.SetDeadline(time.Now().Add(TIMEOUT_SEC))
	}
	return con, err
}

func newHttpClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Dial:                  dialTimeout,
			DisableKeepAlives:     true,
			DisableCompression:    true, // 圧縮解凍は全てこっちで指示する
			ResponseHeaderTimeout: TIMEOUT_SEC,
		},
		CheckRedirect: unutil.RedirectPolicy,
	}
}

func responseRead(resp *http.Response) (data []byte, err error) {
	var r io.Reader
	var gz io.ReadCloser

	ce := resp.Header.Get("Content-Encoding")
	if ce == "gzip" {
		// 解凍する
		gz, _ = gzip.NewReader(resp.Body)
		r = gz
	} else {
		// 圧縮されていない場合
		r = resp.Body
	}
	data, err = ioutil.ReadAll(io.LimitReader(r, DAT_MAX_SIZE))
	if gz != nil {
		gz.Close()
	}
	return
}

func responseWrapper(resp *http.Response) io.ReadCloser {
	var r io.Reader
	cl := make([]io.Closer, 0, 2)

	ce := resp.Header.Get("Content-Encoding")
	if ce == "gzip" {
		// 解凍する
		gz, _ := gzip.NewReader(resp.Body)
		cl = append(cl, gz)
		r = gz
	} else {
		// 圧縮されていない場合
		r = resp.Body
	}
	cl = append(cl, resp.Body)

	return struct {
		io.Reader
		io.Closer
	}{io.LimitReader(r, DAT_MAX_SIZE), unutil.MultiCloser(cl...)}
}

func getHttpBBSmenu(cache Cache) (rc io.ReadCloser, mod time.Time, err error) {
	fmt.Printf("getHttpBBSmenu\r\n") // DEBUG
	client := newHttpClient()
	// header生成
	req, nrerr := http.NewRequest("GET", "http://"+CONF_ITAURL_HOST+"/"+CONF_ITAURL_FILE, nil)
	if nrerr != nil {
		return nil, mod, nrerr
	}
	req.Header.Set("User-Agent", USER_AGENT)
	// 更新確認
	if st, merr := cache.Stat("", "", ""); merr == nil {
		req.Header.Set("If-Modified-Since", unutil.CreateModString(st.Mmod()))
	}
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("Connection", "close")
	fmt.Print("%+v\r\n", req) // DEBUG
	resp, doerr := client.Do(req)
	if doerr != nil {
		return nil, mod, doerr
	}
	code := resp.StatusCode
	if t, lerr := http.ParseTime(resp.Header.Get("Last-Modified")); lerr == nil {
		mod = t
	} else {
		mod = time.Now()
	}

	if code == 200 {
		// レスポンスボディをラップする
		rc = responseWrapper(resp)
	} else {
		err = errors.New("更新されていません")
	}
	return
}

// 板一覧取得
func saveBBSmenu(cache Cache) io.Reader {
	rc, mod, err := getHttpBBSmenu(cache)
	if err != nil {
		// errがnil以外の時、rcはnil
		return nil
	}

	// これ以降はUTF-8
	data := bytes.Buffer{}
	scanner := bufio.NewScanner(unutil.ShiftJISToUtf8Reader(rc))
	for scanner.Scan() {
		line := scanner.Text()
		if match := RegServerItem.FindStringSubmatch(line); match != nil {
			// 当てはまるものを除外
			if _, ok := catekill[match[1]]; !ok {
				data.WriteString(match[1] + "\n")
			}
		} else if strings.Contains(line, ".2ch.net/") || strings.Contains(line, ".bbspink.com/") {
			if strings.Contains(line, "TARGET") {
				continue
			}
			if match := RegServer.FindStringSubmatch(line); match != nil {
				server := match[1]
				board := match[2]
				title := match[3]
				if _, ok := sabakill[server]; ok {
					continue
				}
				data.WriteString(server + "/" + board + "<>" + title + "\n")
			}
		}
	}
	// 全部閉じる
	rc.Close()
	// ファイルにはUTF-8で保存
	cache.SetData("", "", "", data.Bytes())
	cache.SetMod("", "", "", mod, mod)
	return &data
}

func (g2ch *Get2ch) GetBBSmenu(flag bool) (rc io.ReadCloser) { // trueがデフォルト
	var r io.Reader
	if g2ch.cache.Exists("", "", "") == false {
		// 存在しない場合取得する
		r = saveBBSmenu(g2ch.cache)
	}
	if flag {
		if st, err := g2ch.cache.Stat("", "", ""); err == nil {
			g2ch.mod = st.Mmod()
		}
	}
	if r == nil {
		var err error
		rc, err = g2ch.cache.GetDataRC("", "", "")
		if err != nil {
			rc = nil
		}
	} else {
		rc = struct {
			io.Reader
			io.Closer
		}{r, unutil.MultiCloser()}
	}
	return
}

func (g2ch *Get2ch) GetServer(board_key string) string {
	retdata := ""
	if board_key == "" {
		retdata = g2ch.server
	} else {
		retdata = boardServerObj.GetServer(board_key)
	}
	// サーバー名を5chに置換
	retdata = strings.Replace(retdata, "2ch", "5ch", -1)
	return retdata
}

func setServerList() map[string]string {
	var r io.Reader
	m := make(map[string]string, 1024)
	cache := NewFileCache(CONF_FOLDER)
	if cache.Exists("", "", "") == false {
		// 存在しない場合取得する
		r = saveBBSmenu(cache)
	}
	if r == nil {
		rc, err := cache.GetDataRC("", "", "")
		if err != nil {
			return m
		}
		defer rc.Close()
		r = rc
	}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		sp := strings.Split(scanner.Text()+"<>", "<>")
		dat, name := sp[0], sp[1]
		if name != "" {
			u := strings.Split(dat, "/")
			server, board := u[0], u[1]
			if _, ok := m[board]; !ok {
				// 存在しなかったらセットする
				m[board] = server
			}
		}
	}
	// 隠し板をロードする
	for board, it := range hideboard {
		if _, ok := m[board]; !ok {
			m[board] = it.server
		}
	}
	return m
}

func getBoardNameSub(bd string) string {
	rc, err := NewFileCache(CONF_FOLDER).GetDataRC("", "", "")
	if err != nil {
		return ""
	}
	defer rc.Close()
	scanner := bufio.NewScanner(rc)
	for scanner.Scan() {
		sp := strings.Split(scanner.Text()+"<>", "<>")
		dat, name := sp[0], sp[1]
		if name != "" {
			u := strings.Split(dat, "/")
			board := u[1]
			if bd == board {
				return name
			}
		}
	}
	// 隠し板をロードする
	for board, it := range hideboard {
		if bd == board {
			return it.name
		}
	}
	return ""
}

// 板名取得
func (g2ch *Get2ch) GetBoardName() (boardname string) {
	boardname = boardNameObj.GetName(g2ch.board)

	if boardname == "" {
		boardname = g2ch.sliceBoardName()
		if boardname == "" {
			boardname = getBoardNameSub(g2ch.board)
		}

		// 空白でも登録
		boardNameObj.SetName(g2ch.board, boardname)
	}
	return
}

func (g2ch *Get2ch) getSettingFile() ([]byte, error) {
	fmt.Printf("Get2ch::getSettingFile\r\n") // DEBUG
	server := g2ch.server
	board := g2ch.board
	req_time := g2ch.req_time

	if server == "" {
		// サーバが見つからない場合
		return nil, errors.New("サーバが無いよ")
	}

	var cf bool
	if g2ch.bourbon {
		// バーボン中
		cf = true
	} else {
		// 未来の時間
		if st, err := g2ch.cache.Stat(server, board, BOARD_SETTING); err == nil {
			cf = st.Mmod().After(req_time)
		}
	}
	if cf {
		// Cacheを返す
		// UTF-8に変換
		cdata, err := g2ch.cache.GetData(server, board, BOARD_SETTING)
		if err != nil {
			cdata = []byte{}
		}
		return unutil.ShiftJISToUtf8(cdata), nil
	}

	client := newHttpClient()
	// header生成
	req, nrerr := http.NewRequest("GET", "http://"+server+"/"+board+"/"+FILE_SETTING_TXT_REQ, nil)
	if nrerr != nil {
		return nil, nrerr
	}
	req.Header.Set("User-Agent", USER_AGENT)
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("Connection", "close")
	fmt.Print("%+v\r\n", req) // DEBUG

	resp, doerr := client.Do(req)
	if doerr != nil {
		return nil, doerr
	}

	var data []byte
	var err error
	code := resp.StatusCode
	if code == 200 {
		// 読み込む
		if data, err = responseRead(resp); err == nil {
			g2ch.cache.SetData(server, board, BOARD_SETTING, data)
			mod := req_time.Add(SETTING_CACHE_TIME)
			g2ch.cache.SetMod(server, board, BOARD_SETTING, mod, mod)
		}
	} else {
		// 板名取得失敗
		// 特にエラーとしない
		if data, err = g2ch.cache.GetData(server, board, BOARD_SETTING); err != nil {
			// ファイルが存在しない場合
			data = []byte{}
		}
	}
	// 閉じる
	resp.Body.Close()
	// 返す際にUTF-8に変換
	return unutil.ShiftJISToUtf8(data), nil
}

func (g2ch *Get2ch) sliceBoardName() (bname string) {
	stf, err := g2ch.getSettingFile()
	if err != nil {
		return
	}
	start_text := []byte("BBS_TITLE=")
	if start := bytes.Index(stf, start_text); start >= 0 {
		start += len(start_text)
		if end := bytes.IndexByte(stf[start:], '\n'); end >= 0 {
			stf = stf[start : start+end]
			var name string
			if i := bytes.Index(stf, []byte("＠")); i >= 0 {
				name = string(stf[:i])
			} else {
				name = string(stf)
			}
			bname = strings.Trim(name, " \t")
		}
	}
	return
}

// header送信
func (g2ch *Get2ch) request(flag bool) (data []byte) {
	fmt.Printf("Get2ch::request\r\n") // DEBUG
	select {
	case parallelRequestLimitCh <- struct{}{}:
		// 同時実行数制限
	default:
		// 詰まってる
		g2ch.code = 0
		return
	}
	defer func() {
		<-parallelRequestLimitCh
	}()

	var req *http.Request
	var err error
	server := g2ch.server
	board := g2ch.board
	thread := g2ch.thread
	req_time := g2ch.req_time

	if server == "" {
		// サーバが分からない
		g2ch.code = 0
		return
	} else if g2ch.isThread() {
		// 2015-3-14 2ch新仕様対応で削除
		//// dat取得用header生成
		//req, err = http.NewRequest("GET", "http://"+server+"/"+board+"/dat/"+thread+".dat", nil)
		// 2015-3-14 2ch新仕様対応
		fmt.Printf("get thread html of 2ch directly:\r\n")                                   // DEBUG
		fmt.Printf("https://" + server + "/test/read.cgi/" + board + "/" + thread + "/\r\n") // DEBUG
		req, err = http.NewRequest("GET", "https://"+server+"/test/read.cgi/"+board+"/"+thread+"/", nil)
		if err != nil {
			return
		}
		req.Header.Set("User-Agent", USER_AGENT)

		/*2015-3-14 2ch新仕様対応で削除
		st, err := g2ch.cache.Stat(server, board, thread)
		if flag && err == nil {
			timem := st.Mmod()
			timea := st.Amod()
			if req_time.Before(timem) {
				// dat落ちしている場合
				g2ch.code = 0
				g2ch.cache_mod = timem
				return
			} else if req_time.Before(timem.Add(DAT_CACHE_TIME_THREAD)) || req_time.Before(timea.Add(DAT_CACHE_TIME_THREAD)) {
				// 前回の取得から数秒しか経過していない場合
				// 変更なしとする
				g2ch.code = 429
				return
			} else {
				// 多重書き込みを防ぐため早めに取得待機時間を延長しておく
				g2ch.cache.SetMod(server, board, thread, timem, req_time)
				size := st.Size()
				if size > 1 {
					// 1バイト引いても差分取得ができる場合
					// 1バイト引いて取得する
					req.Header.Set("Range", "bytes="+strconv.Itoa(int(size-1))+"-")
				}
				req.Header.Set("If-Modified-Since", unutil.CreateModString(timem))
			}
		} else {
			// 差分取得は使えないためここで設定
			req.Header.Set("Accept-Encoding", "gzip")
		}
		*/
		// 2015-3-14 2ch新仕様対応 常に取得する
		req.Header.Set("Accept-Encoding", "gzip")
	} else if g2ch.isBoard() {
		// スレッド一覧取得用header生成
		req, err = http.NewRequest("GET", "http://"+server+"/"+board+"/"+FILE_SUBJECT_TXT_REQ, nil)
		if err != nil {
			return
		}
		req.Header.Set("User-Agent", USER_AGENT)

		if st, err := g2ch.cache.Stat(server, board, ""); err == nil {
			timem := st.Mmod()
			timea := st.Amod()
			if req_time.Before(timem.Add(DAT_CACHE_TIME_BOARD)) || req_time.Before(timea.Add(DAT_CACHE_TIME_BOARD)) {
				// 前回の取得から数秒しか経過していない場合
				// 変更なしとする
				g2ch.code = 429
				return
			}
			// 早めに取得待機時間を延長しておく
			g2ch.cache.SetMod(server, board, "", timem, req_time)
			req.Header.Set("If-Modified-Since", unutil.CreateModString(timem))
		}
		req.Header.Set("Accept-Encoding", "gzip")
	} else {
		g2ch.code = 0
		return
	}
	req.Header.Set("Connection", "close")
	fmt.Print("%+v\r\n", req) // DEBUG

	// リクエスト送信
	var resp *http.Response
	client := newHttpClient()
	resp, err = client.Do(req)
	if err != nil {
		// errがnil以外の場合、resp.Bodyは閉じられている
		if resp == nil {
			g2ch.code = 0
		} else {
			if rerr := unutil.GetRedirectError(err); rerr != nil {
				// RedirectErrorだった場合は処理続行
				// バーボン判定
				if strings.Contains(rerr.Path, "403") {
					// バーボン状態
					g2ch.bourbon = true
				}
				g2ch.code = resp.StatusCode
			} else {
				g2ch.code = 0
			}
		}
		// 終了
		return
	}

	// 読み込み
	data, err = responseRead(resp)
	resp.Body.Close()
	if err != nil {
		g2ch.code = 0
		return nil
	}

	g2ch.code = resp.StatusCode
	g2ch.size = int64(len(data))
	var mod time.Time
	if t, perr := http.ParseTime(resp.Header.Get("Last-Modified")); perr == nil {
		mod = t
	}

	if g2ch.code == 304 {
		// データは空
		data = []byte{}
	} else if flag && (g2ch.code == 206) && (g2ch.size > 1) {
		// あぼーん検知
		data = lfCheck(data)
		if data == nil {
			g2ch.code = 416
		}
	}
	if mod.IsZero() == false {
		g2ch.cache_mod = mod
	} else {
		g2ch.cache_mod = req_time
	}
	return
}

func (g2ch *Get2ch) bourbonRequest() (data []byte) {
	fmt.Printf("Get2ch::bourbonRequest\r\n")
	select {
	case parallelRequestLimitCh <- struct{}{}:
		// 同時実行数制限
	default:
		// 詰まってる
		g2ch.code = 0
		return
	}
	defer func() {
		<-parallelRequestLimitCh
	}()

	var req *http.Request
	var err error
	server := g2ch.server
	board := g2ch.board
	thread := g2ch.thread
	strerr := append(make([]byte, 0, len(tanpanman)), tanpanman...)

	if server == "" {
		// サーバが分からない
		return
	} else if g2ch.isThread() {
		// dat取得用header生成
		req, err = http.NewRequest("GET", "http://"+BOURBON_HOST+"/test/r.so/"+server+"/"+board+"/"+thread+"/", nil)
		if err != nil {
			return
		}
	} else if g2ch.isBoard() {
		// スレッド一覧取得用header生成
		req, err = http.NewRequest("GET", "http://"+BOURBON_HOST+"/test/p.so/"+server+"/"+board+"/", nil)
		if err != nil {
			return
		}
		// 更新確認
		if st, err := g2ch.cache.Stat(server, board, ""); err == nil {
			req.Header.Set("If-Modified-Since", unutil.CreateModString(st.Mmod()))
		}
	} else {
		return strerr
	}
	req.Header.Set("User-Agent", USER_AGENT)
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("Connection", "close")
	fmt.Print("%+v\r\n", req) // DEBUG

	// リクエスト送信
	var resp *http.Response
	client := newHttpClient()
	resp, err = client.Do(req)
	if err != nil {
		// errがnil以外の場合、resp.Bodyは閉じられている
		g2ch.code = 0
		return strerr
	}

	// 読み込み
	data, err = responseRead(resp)
	resp.Body.Close()
	if err != nil {
		g2ch.code = 0
		return strerr
	}

	g2ch.code = resp.StatusCode
	g2ch.size = int64(len(data))
	g2ch.mod, g2ch.cache_mod = g2ch.req_time, g2ch.req_time
	return data
}

func (g2ch *Get2ch) normalData(reget bool) []byte {
	var err error
	// データ取得
	data := g2ch.request(reget)
	if g2ch.isThread() {
		switch g2ch.code {
		case 200:
			g2ch.createCache(data, DAT_CREATE)
			g2ch.addThreadCache(data)
		case 206:
			g2ch.createCache(data, DAT_APPEND)
			data, err = g2ch.readThread()
			if err != nil {
				data = g2ch.dataErrorDat()
			}
		case 416:
			if reget {
				// もう一回取得
				data = g2ch.normalData(false)
			} else {
				data, err = g2ch.readThread()
				if err != nil {
					data = g2ch.dataErrorDat()
				}
			}
		case 301, 302, 404:
			if st, staterr := g2ch.cache.Stat(g2ch.server, g2ch.board, g2ch.thread); staterr == nil {
				g2ch.size = st.Size()
				g2ch.mod = st.Mmod()
				if g2ch.size < DAT_MAX_SIZE {
					data, err = g2ch.cache.GetData(g2ch.server, g2ch.board, g2ch.thread)
					if err != nil {
						data = g2ch.dataError()
						break
					}
					g2ch.addThreadCache(data)
					// バーボンキャッシュ更新
					updateBourbonCache(g2ch.bourbon)
					// dat落ち判定
					if g2ch.bourbon == false && g2ch.mod.Before(g2ch.req_time.Add(-1*DAT_NOT_REQUEST_WAIT)) {
						// バーボン状態でなければ5年後まで読みに行かないようにする
						mod := g2ch.req_time.Add(DAT_NOT_REQUEST_ADD_MOD)
						g2ch.cache.SetMod(g2ch.server, g2ch.board, g2ch.thread, mod, mod)
						g2ch.mod = mod
					}
				} else {
					data = g2ch.dataError()
				}
			} else {
				data = g2ch.dataErrorDat()
				if g2ch.bourbon == false {
					// バーボン状態でなければ404
					g2ch.code404 = true
				}
			}
		default:
			// case 0, 304, 400, 409, 429, 503, 504:
			// 更新無し or salamiのエラー or キャッシュ利用
			data, err = g2ch.readThread()
			if err != nil {
				data = g2ch.dataErrorDat()
			}
		}
	} else {
		switch g2ch.code {
		case 200:
			g2ch.createCache(data, DAT_CREATE)
			g2ch.setBoardThread(data)
		case 301, 302, 404:
			// 鯖情報取得
			saveBBSmenu(g2ch.cache)
			data = []byte{}
			g2ch.err = errors.New("２ちゃんねるにアクセスできなかったので、サーバー移転チェックを行いました。")
		default:
			// case 206, 304, 400, 409, 429, 503, 504:
			// 更新無し or salamiのエラー or キャッシュ利用
			data, err = g2ch.readBoard()
			if err != nil {
				data = g2ch.dataErrorDat()
			}
		}
	}
	return data
}

func (g2ch *Get2ch) bourbonData() (data []byte) {
	time := g2ch.req_time
	g2ch.bourbon = true

	// タイムスタンプがいじられている場合、見に行かない
	if st, staterr := g2ch.cache.Stat(g2ch.server, g2ch.board, g2ch.thread); staterr == nil {
		timem := st.Mmod()
		if time.Before(timem) {
			// 現在時刻よりもファイルの更新時刻のほうが大きい場合
			// dat落ちしていることとする
			g2ch.code = 0
			g2ch.mod = timem
			g2ch.size = st.Size()
			if g2ch.size < DAT_MAX_SIZE {
				data, _ = g2ch.cache.GetData(g2ch.server, g2ch.board, g2ch.thread)
				if g2ch.isThread() {
					// スレッドの場合はキャッシュにセット
					g2ch.addThreadCache(data)
				}
			} else {
				data = g2ch.dataError()
			}
			return
		} else if time.Before(timem.Add(DAT_CACHE_TIME_THREAD)) {
			// 前回の取得から数秒しか経過していない場合
			g2ch.code = 429
			g2ch.mod = timem
			g2ch.size = st.Size()
			if g2ch.size < DAT_MAX_SIZE {
				data, _ = g2ch.cache.GetData(g2ch.server, g2ch.board, g2ch.thread)
				if g2ch.isThread() {
					// スレッドの場合はキャッシュにセット
					g2ch.addThreadCache(data)
				}
			} else {
				data = g2ch.dataError()
			}
			return
		}
	}

	if strings.Contains(g2ch.server, ".bbspink.com") {
		// BBSPINKだった場合
		data = append(make([]byte, 0, len(tanpanmanNagoyaee)), tanpanmanNagoyaee...)
	} else {
		data = g2ch.bourbonRequest()
	}
	tp := append(make([]byte, 0, len(tanpanman)), tanpanman...)
	ne := append(make([]byte, 0, len(nagoyaee)), nagoyaee...)
	checklen := len(data)
	if checklen > 1024 {
		checklen = 1024
	}
	if bytes.Contains(data[:checklen], tp) || bytes.Contains(data[:checklen], ne) {
		// 取得に失敗した場合
		g2ch.code = 302
		if st, staterr := g2ch.cache.Stat(g2ch.server, g2ch.board, g2ch.thread); staterr == nil {
			g2ch.mod = st.Mmod()
			g2ch.size = st.Size()
			if g2ch.size < DAT_MAX_SIZE {
				data, _ = g2ch.cache.GetData(g2ch.server, g2ch.board, g2ch.thread)
				if g2ch.isThread() {
					g2ch.addThreadCache(data)
					// これから先にリクエストを送る必要がないか判断する
					g2ch.checkNoRequest(data, true)
				}
			} else {
				data = g2ch.dataError()
			}
		} else {
			data = g2ch.dataErrorDat()
		}
	} else {
		// 取得に成功した場合
		if g2ch.isThread() {
			g2ch.createCache(data, DAT_BOURBON_THREAD)
			g2ch.addThreadCache(data)
		} else {
			g2ch.createCache(data, DAT_BOURBON_BOARD)
			g2ch.setBoardThread(data)
		}
	}
	return
}

func (g2ch *Get2ch) dataError() []byte {
	data := bytes.Buffer{}
	g2ch.err = errors.New("壊れているため表示できません。")
	if g2ch.isThread() {
		data.WriteString("unkar.org<><>")
		data.WriteString(unutil.CreateDateString(g2ch.req_time))
		data.WriteString("<>DATが壊れているため表示できません。<>なんかえらーだって\n")
	} else {
		data.WriteString(strconv.Itoa(int(g2ch.req_time.Unix())))
		data.WriteString(".dat<>板が壊れているため表示できません (1)\n")
	}
	return unutil.Utf8ToShiftJIS(data.Bytes())
}

func (g2ch *Get2ch) dataErrorDat() []byte {
	data := bytes.Buffer{}
	g2ch.err = errors.New("アクセス不可(dat落ち)")
	if g2ch.isThread() {
		data.WriteString("unkar.org<><>")
		data.WriteString(unutil.CreateDateString(g2ch.req_time))
		data.WriteString("<>スレッドを発見できませんでした。dat落ちのようです。<>アクセス不可(dat落ち)\n")
	} else {
		data.WriteString(strconv.Itoa(int(g2ch.req_time.Unix())))
		data.WriteString(".dat<>２ちゃんねるにアクセスできませんでした。 (1)\n")
	}
	return unutil.Utf8ToShiftJIS(data.Bytes())
}

// この関数の引数はSJIS-winであること
func (g2ch *Get2ch) addThreadCache(data []byte) error {
	if !g2ch.isThread() {
		return errors.New("thread")
	}

	i := bytes.IndexByte(data, '\n')
	if i < 0 {
		return errors.New("LF error")
	}
	resu := bytes.Split(unutil.StripTags(data[:i:i], HtmlTag), []byte("<>"))
	if len(resu) > 4 {
		var lastdate time.Time
		st, err := g2ch.cache.Stat(g2ch.server, g2ch.board, g2ch.thread)
		if err == nil {
			lastdate = st.Mmod()
		}
		titem := process.ViewThreadItem{
			Board:     g2ch.board,
			Boardname: g2ch.GetBoardName(),
			Thread:    g2ch.thread,
			Title:     string(unutil.ShiftJISToUtf8(resu[4])),
			Res:       g2ch.NumLines(data),
			Lastdate:  lastdate,
			Addtime:   g2ch.req_time,
		}
		// スレッド情報を詰める
		viewThreadListObj.SetThread(titem)
	} else {
		// データ破損
		return errors.New("data error")
	}
	return nil
}

// 必ずSJIS-winの状態で渡す
func lfCheck(data []byte) []byte {
	// ソースはUTF-8で文字列はSJIS-win
	// 改行コードはASCIIの範囲なので問題なし
	if data[0] == '\n' {
		return data[1:]
	}
	return nil
}

// 必ずSJIS-winの状態で渡す
func (g2ch *Get2ch) createCache(data []byte, switch_data int) error {
	mod := g2ch.cache_mod
	append_data := false
	renew := true

	if data == nil {
		return errors.New("data nil")
	}

	switch switch_data {
	case DAT_CREATE:
		append_data = false
	case DAT_APPEND:
		if len(data) > 0 {
			// データが存在するので追記
			append_data = true
		} else {
			// データが更新されていない
			renew = false
		}
	case DAT_BOURBON_THREAD:
		append_data = false
		st, err := g2ch.cache.Stat(g2ch.server, g2ch.board, g2ch.thread)
		if err == nil {
			if int64(len(data)) <= st.Size() {
				// データが更新されていない
				renew = false
			}
		}
	case DAT_BOURBON_BOARD:
		// 何もしない
	default:
		return errors.New("invalid arguments")
	}

	if renew {
		// バーボン中ではない、またはデータが更新されている場合
		// ファイルに書き込む
		if append_data {
			// 追記する
			g2ch.cache.SetDataAppend(g2ch.server, g2ch.board, g2ch.thread, data) // 追記
			if g2ch.isThread() {
				// データの追記がある
				data, _ = g2ch.cache.GetData(g2ch.server, g2ch.board, g2ch.thread)
				// DBを更新
				dbUpdateThread.SetItem(g2ch.board, g2ch.thread, g2ch.NumLines(data))
			}
		} else {
			flag := g2ch.cache.Exists(g2ch.server, g2ch.board, g2ch.thread)
			g2ch.cache.SetData(g2ch.server, g2ch.board, g2ch.thread, data) // 上書き
			if g2ch.isThread() && flag == false {
				// ファイルが存在しない
				// DBに挿入
				dbInsertThread.SetItem(data, g2ch.board, g2ch.thread, g2ch.NumLines(data))
			}
		}
	}
	if g2ch.isThread() {
		// スレッドの場合
		// これから先にリクエストを送る必要がないか判断する
		mod = g2ch.checkNoRequest(data, false)
	}
	// If-Modified-Sinceをセット
	if mod.IsZero() == false {
		g2ch.mod = mod
		g2ch.cache.SetMod(g2ch.server, g2ch.board, g2ch.thread, mod, mod)
	}
	return nil
}

func (g2ch *Get2ch) checkNoRequest(data []byte, flag bool) time.Time {
	mod := g2ch.cache_mod

	if g2ch.isThread() {
		// スレッドの場合
		cacheflag := false

		if g2ch.NumLines(data) >= DAT_NOT_REQUEST_RES_COUNT {
			// 1000res以上
			cacheflag = true
		} else if len(data) > DAT_NOT_SIZE_LIMIT {
			// 512kbyteよりも大きい
			cacheflag = true
		}

		if cacheflag {
			// 5年後まで読みに行かないようにする
			mod = g2ch.req_time.Add(DAT_NOT_REQUEST_ADD_MOD)
			if flag {
				// 未来の時間を設定する
				g2ch.cache.SetMod(g2ch.server, g2ch.board, g2ch.thread, mod, mod)
			}
		}
	}
	return mod
}

func getBourbonCache() (bourbon bool) {
	return bbnCacheObj.GetBourbon(BOURBON_KEY)
}

func updateBourbonCache(bin bool) {
	if bin {
		// バーボン状態
		bbnCacheObj.SetBourbon(BOURBON_KEY)
	}
}

func (g2ch *Get2ch) readThread() (data []byte, err error) {
	var st CacheState
	if st, err = g2ch.cache.Stat(g2ch.server, g2ch.board, g2ch.thread); err == nil {
		g2ch.size = st.Size()
		g2ch.mod = st.Mmod()
		data, err = g2ch.cache.GetData(g2ch.server, g2ch.board, g2ch.thread)
		if err == nil {
			g2ch.addThreadCache(data)
		}
	}
	return
}

func (g2ch *Get2ch) readBoard() (data []byte, err error) {
	var st CacheState
	if st, err = g2ch.cache.Stat(g2ch.server, g2ch.board, ""); err == nil {
		g2ch.size = st.Size()
		g2ch.mod = st.Mmod()
		data, err = g2ch.cache.GetData(g2ch.server, g2ch.board, "")
	}
	return
}

func (g2ch *Get2ch) boardThreadLookup() bool {
	f, err := boardThreadMap.Lookup(g2ch.board, g2ch.thread)
	if f && err != nil {
		// まだmapが登録されていない
		if data, err := g2ch.cache.GetData(g2ch.server, g2ch.board, ""); err == nil {
			// TODO:atimeが更新される？影響の調査が必要
			// 登録する
			g2ch.setBoardThread(data)
			// もう一度探索
			f, _ = boardThreadMap.Lookup(g2ch.board, g2ch.thread)
		} else {
			// ファイルが無い
			f = false
		}
	}
	return f
}
func (g2ch *Get2ch) setBoardThread(data []byte) {
	lines := 0
	m := make(map[int64]struct{}, 64)
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		line := scanner.Text()
		lines++
		i := strings.Index(line, ".dat")
		if i <= 8 || i >= 11 {
			continue
		}
		if num, err := strconv.ParseInt(line[:i], 10, 64); err == nil {
			m[num] = struct{}{}
		}
	}
	if len(m) > (lines / 2) {
		boardThreadMap.SetMap(g2ch.board, m)
	}
}
