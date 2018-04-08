package unutil

// 細かい便利機能

import (
	"../conf"
	"bytes"
	// Go ver 1.4未満
	//"code.google.com/p/go.net/html"
	//"code.google.com/p/go.net/html/atom"
	//"code.google.com/p/mahonia"
	"compress/flate"
	"compress/gzip"
	"fmt"
	"github.com/cypro666/mahonia"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"hash/fnv"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Stat_t struct {
	Size  int64
	Atime time.Time
	Mtime time.Time
}

type Output struct {
	Code   int
	Header http.Header
	Reader io.Reader
	ZFlag  bool
}

type DummyResponseWriter struct {
	Code int
	h    http.Header
	Buf  *bytes.Buffer
}

type denyMap struct {
	ipmap map[byte]*denyMap
	mask  byte
}

type Model interface {
	GetData() interface{}
	GetUrl() string
	GetMod() time.Time
	GetTitle() string
	GetClassName() string
	GetError() error
	Is404() bool
	GetServer() string
	GetByteSize() int64
	GetCode() int
}

type View interface {
	PrintData(Model) Output
	GetHostUrl() string
}

var code403Html []byte
var code503Html []byte
var code403HtmlGzip []byte
var code503HtmlGzip []byte
var denyNetwork *denyMap

var Stdlog chan<- string
var Errlog chan<- string
var RegOldUrlThread = regexp.MustCompile(`^\/(?:r(?:ead)?)\/(?:\w+\.2ch\.net|\w+\.bbspink\.com)\/(\w+)\/(\d{9,10})(\/.*)?`)
var RegOldUrlBoard = regexp.MustCompile(`^\/(?:r(?:ead)?)\/(?:\w+\.2ch\.net|\w+\.bbspink\.com)\/(\w+)(\/.*)?`)

func (out *Output) Error() (ret string) {
	if out.Reader != nil {
		buf, err := ioutil.ReadAll(out.Reader)
		if err == nil {
			ret = string(buf)
		} else {
			ret = err.Error()
		}
	} else {
		ret = "Output Error"
	}
	return
}

func init() {
	code403Html, code403HtmlGzip = MustReadFile("undity/template/code_403.templ")
	code503Html, code503HtmlGzip = MustReadFile("undity/template/code_503.templ")
	Stdlog = loggerProc(os.Stdout, 128)
	Errlog = loggerProc(os.Stderr, 1)
}

func MustReadFile(path string) (buf []byte, gzbuf []byte) {
	var err error
	buf, err = ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	// 事前圧縮なので最大圧縮率で圧縮
	data := bytes.Buffer{}
	gz, _ := gzip.NewWriterLevel(&data, gzip.BestCompression)
	io.Copy(gz, bytes.NewReader(buf))
	gz.Close()
	gzbuf = data.Bytes()
	return
}

func loggerProc(w io.Writer, bsize int) chan<- string {
	if bsize <= 0 {
		bsize = 1
	}
	c := make(chan string, bsize)
	go writeLogProc(c, w)
	return c
}

func writeLogProc(c <-chan string, w io.Writer) {
	for s := range c {
		if len(s) > 0 && s[len(s)-1] != '\n' {
			s += "\n"
		}
		io.WriteString(w, s)
	}
}

func Dispose(w http.ResponseWriter, r *http.Request) {
	if err := recover(); err != nil {
		var code, size int
		if out, ok := err.(*Output); ok {
			code = out.Code
			size = Print(w, r, *out)
		} else {
			// 500を返しておく
			code = http.StatusInternalServerError
			w.WriteHeader(code)
			s := fmt.Sprintf("%s:%v => %s", CreateDateNowLog(), err, string(Stack()))
			Errlog <- s
		}
		// ログ出力
		Putlog(r, code, size)
	}
}

func Print(resw http.ResponseWriter, r *http.Request, out Output) int {
	// ヘッダー設定
	for key, _ := range out.Header {
		resw.Header().Set(key, out.Header.Get(key))
	}
	resw.Header().Set("X-Frame-Options", "SAMEORIGIN")

	// 出力フォーマット切り替え
	var writer io.Writer
	var wc io.WriteCloser
	sw := NewSizeCountWriter(resw)
	if out.ZFlag {
		resw.Header().Set("Vary", "Accept-Encoding")
		ae := r.Header.Get("Accept-Encoding")
		if strings.Contains(ae, "gzip") {
			// gzip圧縮
			resw.Header().Set("Content-Encoding", "gzip")
			wc, _ = gzip.NewWriterLevel(sw, gzip.BestSpeed)
			writer = wc
		} else if strings.Contains(ae, "deflate") {
			// deflate圧縮
			resw.Header().Set("Content-Encoding", "deflate")
			wc, _ = flate.NewWriter(sw, flate.BestSpeed)
			writer = wc
		} else {
			// 圧縮しない
			writer = sw
		}
	} else {
		// 生データ
		writer = sw
	}
	// ステータスコード＆ヘッダー出力
	resw.WriteHeader(out.Code)
	// ボディ出力
	if r.Method != "HEAD" && out.Reader != nil {
		io.Copy(writer, out.Reader)
	}
	if wc != nil {
		// 圧縮が有効な場合
		wc.Close()
	}
	return sw.Size
}

func Putlog(r *http.Request, code, size int) {
	rh, _, _ := net.SplitHostPort(r.RemoteAddr)
	date := CreateDateNowLog()
	p := r.URL.Path
	if r.URL.RawQuery != "" {
		p += "?" + r.URL.RawQuery
	}
	s := fmt.Sprintf(`%s - - [%s] "%s %s %s" %d %d`, rh, date, r.Method, p, r.Proto, code, size)
	Stdlog <- s
}

func Stack() []byte {
	buf := make([]byte, 32*1024)
	s := runtime.Stack(buf, false)
	return buf[:s:s]
}

func Lazy(r *http.Request) {
	switch r.Method {
	case "GET", "HEAD":
		// OK
	default:
		// NG
		NotImplemented(r)
	}
	if len(r.URL.RequestURI()) >= 1024 {
		RequestURITooLong(r)
	}
}

func Move(r *http.Request) {
	p := r.URL.Path
	q := r.URL.Query()
	ef := q.Get("_escaped_fragment_")
	if ef != "" {
		if ef == "/server" {
			ef = ""
		}
		MovedPermanently("http://unkar.org/r" + ef)
	}

	if strings.Index(p, "/read.html/") == 0 {
		MovedPermanently("http://unkar.org/r" + p[10:])
	} else if p == "/read.html" {
		MovedPermanently("http://unkar.org/r")
	} else if p == "/2ch" {
		MovedPermanently("http://unkar.org/r")
	} else if strings.Index(p, "/2ch/") == 0 {
		ch := p[4:]
		switch ch {
		case "/", "/index.html", "/index2.html", "/index3.html":
			MovedPermanently("http://unkar.org/r")
		case "/search.php":
			MovedPermanently("http://unkar.org/search?" + r.URL.RawQuery)
		}
	} else if m := RegOldUrlThread.FindStringSubmatch(p); m != nil {
		MovedPermanently("http://unkar.org/r/" + m[1] + "/" + m[2])
	} else if m := RegOldUrlBoard.FindStringSubmatch(p); m != nil {
		MovedPermanently("http://unkar.org/r/" + m[1])
	}
}

func Deny(r *http.Request) {
	rh, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		// 403
		Forbidden(r)
	}
	dm := denyNetwork
	var ok bool

	for _, it := range net.ParseIP(rh).To4() {
		if dm, ok = dm.ipmap[it&dm.mask]; !ok {
			return
		}
	}
	// 制限中のネットワークからのアクセス
	// 403
	Forbidden(r)
}

func InitDeny(denylist []string) {
	denyNetwork = newDenyMap()
	for _, it := range denylist {
		_, n, err := net.ParseCIDR(it)
		if err != nil {
			continue
		}
		var m [5]*denyMap
		var ok bool
		masklen := len(n.Mask)

		m[0] = denyNetwork
		for i := 0; i < 4; i++ {
			m[i+1], ok = m[i].ipmap[n.IP[i]]
			if !ok {
				m[i+1] = newDenyMap()
				m[i].ipmap[n.IP[i]] = m[i+1]
			}
			if masklen > i {
				m[i].mask &= n.Mask[i]
			}
		}
	}
}

func newDenyMap() *denyMap {
	return &denyMap{ipmap: make(map[byte]*denyMap), mask: 0xff}
}

func UniqueIntSlice(is []int) []int {
	l := len(is)
	ret := make([]int, l)
	i := 0
	old := 0
	sort.Ints(is)
	if l > 0 {
		old = is[0] - 1
	}
	for _, it := range is {
		if old != it {
			ret[i] = it
			old = it
			i++
		}
	}
	return ret[:i:i]
}

// Less Than for a pair of int arguments
func LTInt(v2, v1 int) bool {
	return v2 < v1
}
func LTUint64(v2, v1 uint64) bool {
	return v2 < v1
}

// Minimum of a pair of int arguments
func MinInt(v1, v2 int) (m int) {
	if LTInt(v2, v1) {
		m = v2
	} else {
		m = v1
	}
	return
}
func MinUint64(v1, v2 uint64) (m uint64) {
	if LTUint64(v2, v1) {
		m = v2
	} else {
		m = v1
	}
	return
}

// Minimum of a slice of int arguments
func MinIntS(v []int) (m int) {
	l := len(v)
	if l > 0 {
		m = v[0]
	}
	for i := 1; i < l; i++ {
		m = MinInt(m, v[i])
	}
	return
}
func MinUintS64(v []uint64) (m uint64) {
	l := len(v)
	if l > 0 {
		m = v[0]
	}
	for i := 1; i < l; i++ {
		m = MinUint64(m, v[i])
	}
	return
}

// Minimum of a variable number of int arguments
func MinIntV(v1 int, vn ...int) (m int) {
	m = v1
	if len(vn) > 0 {
		m = MinInt(m, MinIntS(vn))
	}
	return
}
func MinUintV64(v1 uint64, vn ...uint64) (m uint64) {
	m = v1
	if len(vn) > 0 {
		m = MinUint64(m, MinUintS64(vn))
	}
	return
}

func Range(start, limit, step int) (ret []int) {
	if start <= limit {
		if step <= 0 {
			return
		}
		ret = make([]int, 0, (limit-start)/step)
		for i := start; i <= limit; i += step {
			ret = append(ret, i)
		}
	} else {
		if step >= 0 {
			return
		}
		ret = make([]int, 0, (start-limit)/(step*-1))
		for i := start; i >= limit; i += step {
			ret = append(ret, i)
		}
	}
	return
}

func StripTags(buf []byte, allowable_tags []string) []byte {
	if buf == nil {
		return []byte{}
	}
	data, _ := ioutil.ReadAll(StripTagReader(bytes.NewReader(buf), allowable_tags))
	return data
}

func Utf8Substr(s string, max int) string {
	r := []rune(s)
	l := len(r)
	if l > max {
		l = max
	}
	return string(r[:l:l])
}

func CreateETag(mod time.Time) string {
	h := fnv.New64a()
	// ファイル更新時間
	io.WriteString(h, strconv.FormatInt(mod.Unix(), 36))
	// うんかーのバージョン
	io.WriteString(h, unconf.Ver)
	return `W/"` + strconv.FormatUint(h.Sum64(), 16) + `"`
}

func CreateModString(mod time.Time) string {
	return mod.UTC().Format(http.TimeFormat)
}

func CreateDateString(mod time.Time) string {
	return mod.Format("2006/01/02(Mon) 15:04:05")
}

func CreateDateNowLog() string {
	return time.Now().Format("02/Jan/2006:15:04:05 -0700")
}

func GetIfNoneMatch(r *http.Request) string {
	return r.Header.Get("If-None-Match")
}

func GetIfModifiedSince(r *http.Request) time.Time {
	if m := r.Header.Get("If-Modified-Since"); m != "" {
		since, err := http.ParseTime(m)
		if err == nil {
			return since
		}
	}
	return time.Time{}
}

func CheckNotModified(r *http.Request, mod time.Time) {
	if mod.IsZero() == false {
		data := GetIfModifiedSince(r)
		inm := GetIfNoneMatch(r)
		etag := CreateETag(mod)
		if mod.Before(data.Add(1*time.Second)) || etag == inm || inm == "*" {
			// 更新なし or ETagが同じ
			panic(&Output{
				Code: http.StatusNotModified,
			})
		}
	}
}

func MovedPermanently(u string) {
	// 301
	h := http.Header{}
	h.Set("Location", u)
	msg := `<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<title>301 Moved Permanently</title>
</head>
<body>
<a href="` + u + `">` + u + `</a>
</body>
</html>`
	panic(&Output{
		Code:   http.StatusMovedPermanently,
		Header: h,
		Reader: bytes.NewReader([]byte(msg)),
	})
}

func Forbidden(r *http.Request) {
	// 403
	var msg []byte
	h := http.Header{}
	if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		// gzip圧縮を使う
		h.Set("Content-Encoding", "gzip")
		msg = code403HtmlGzip
	} else {
		msg = code403Html
	}
	h.Set("Vary", "Accept-Encoding")
	h.Set("Content-Type", "text/html; charset=utf-8")
	panic(&Output{
		Code:   http.StatusOK,
		Header: h,
		Reader: bytes.NewReader(msg),
	})
}

func ServiceUnavailable(r *http.Request) {
	// 503
	var msg []byte
	h := http.Header{}
	if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		// gzip圧縮を使う
		h.Set("Content-Encoding", "gzip")
		msg = code503HtmlGzip
	} else {
		msg = code503Html
	}
	h.Set("Vary", "Accept-Encoding")
	h.Set("Content-Type", "text/html; charset=utf-8")
	panic(&Output{
		Code:   http.StatusServiceUnavailable,
		Header: h,
		Reader: bytes.NewReader(msg),
	})
}

func InternalServerError(r *http.Request, msg string) {
	// 500
	h := http.Header{}
	h.Set("Content-Type", "text/html; charset=utf-8")
	panic(&Output{
		Code:   http.StatusInternalServerError,
		Header: h,
		Reader: bytes.NewReader([]byte(msg)),
	})
}

func NotImplemented(r *http.Request) {
	// 501
	h := http.Header{}
	h.Set("Content-Type", "text/html; charset=utf-8")
	h.Set("Public", "GET, HEAD")
	panic(&Output{
		Code:   http.StatusNotImplemented,
		Header: h,
		Reader: bytes.NewReader([]byte(NotImplementedMessage)),
	})
}

func RequestURITooLong(r *http.Request) {
	// 414
	h := http.Header{}
	h.Set("Content-Type", "text/html; charset=utf-8")
	panic(&Output{
		Code:   http.StatusRequestURITooLong,
		Header: h,
		Reader: bytes.NewReader([]byte(RequestURITooLongMessage)),
	})
}

func ShiftJISToUtf8(data []byte) []byte {
	return []byte(ShiftJISToUtf8String(string(data)))
}
func ShiftJISToUtf8String(data string) string {
	return mahonia.NewDecoder("cp932").ConvertString(data)
}
func ShiftJISToUtf8Reader(r io.Reader) io.Reader {
	return mahonia.NewDecoder("cp932").NewReader(r)
}

func Utf8ToShiftJIS(data []byte) []byte {
	buf := bytes.NewBuffer(make([]byte, 0, len(data)))
	enc := mahonia.NewEncoder("cp932")
	enc.NewWriter(buf).Write(data)
	return buf.Bytes()
}
func Utf8ToShiftJISWriter(w io.Writer) io.Writer {
	return mahonia.NewEncoder("cp932").NewWriter(w)
}

func IsMobile(r *http.Request) bool {
	useragents := []string{
		"iPhone",  // Apple iPhone
		"iPad",    // Apple iPad
		"iPod",    // Apple iPod touch
		"Android", // 1.5+ Android
		"Windows Phone",
		"dream",   // Pre 1.5 Android
		"CUPCAKE", // 1.5+ Android
		"PlayStation Vita",
	}
	ua := r.Header.Get("User-Agent")
	if ua != "" {
		for _, it := range useragents {
			if strings.Contains(ua, it) {
				return true
			}
		}
	}
	return false
}

func IsBot(r *http.Request) bool {
	bots := []string{
		"googlebot",
		"slurp",
		"y!j",
		"msnbot",
		"spider",
		"robot",
		"crawl",
	}
	ua := r.Header.Get("User-Agent")
	if ua != "" {
		for _, it := range bots {
			if strings.Contains(ua, it) {
				return true
			}
		}
	}
	return false
}

func NewResponseWriter() *DummyResponseWriter {
	return &DummyResponseWriter{
		Code: http.StatusOK,
		h:    http.Header{},
		Buf:  &bytes.Buffer{},
	}
}

func (drw *DummyResponseWriter) Header() http.Header {
	return drw.h
}

func (drw *DummyResponseWriter) Write(p []byte) (int, error) {
	return drw.Buf.Write(p)
}

func (drw *DummyResponseWriter) WriteHeader(code int) {
	drw.Code = code
}

type RedirectError struct {
	Host string
	Path string
	Msg  string
}

func (e *RedirectError) Error() string {
	return e.Msg
}

func RedirectPolicy(r *http.Request, _ []*http.Request) error {
	return &RedirectError{r.URL.Host, r.URL.Path, "redirect error"}
}

func GetRedirectError(err error) *RedirectError {
	uerr, uok := err.(*url.Error)
	if !uok {
		return nil
	}
	rerr, rok := uerr.Err.(*RedirectError)
	if !rok {
		return nil
	}
	return rerr
}

type sizeCountWriter struct {
	w    io.Writer
	Size int
}

func NewSizeCountWriter(w io.Writer) *sizeCountWriter {
	return &sizeCountWriter{w: w}
}

func (scw *sizeCountWriter) Write(p []byte) (n int, err error) {
	n, err = scw.w.Write(p)
	scw.Size += n
	return
}

type multiCloser struct {
	closers []io.Closer
}

func MultiCloser(closers ...io.Closer) io.Closer {
	return &multiCloser{closers}
}

func (mc *multiCloser) Close() (err error) {
	for _, c := range mc.closers {
		err = c.Close()
	}
	return
}

type stripTag struct {
	d   *html.Tokenizer
	a   map[atom.Atom]struct{}
	buf []byte
}

func StripTagReader(r io.Reader, allowable_tags []string) io.Reader {
	allowed := make(map[atom.Atom]struct{})
	for _, it := range allowable_tags {
		if a := atom.Lookup([]byte(strings.ToLower(it))); a != 0 {
			allowed[a] = struct{}{}
		}
	}
	return &stripTag{
		d:   html.NewTokenizer(r),
		a:   allowed,
		buf: []byte{},
	}
}

func (st *stripTag) Read(buf []byte) (i int, err error) {
	l := len(buf)
	data := bytes.NewBuffer(st.buf)
	data.Grow(l)
	for tt := st.d.Next(); tt != html.ErrorToken; tt = st.d.Next() {
		switch tt {
		case html.StartTagToken, html.EndTagToken, html.SelfClosingTagToken:
			name, _ := st.d.TagName()
			if _, ok := st.a[atom.Lookup(name)]; ok {
				data.Write(st.d.Raw())
			}
		//case html.TextToken, html.CommentToken, html.DoctypeToken:
		case html.TextToken:
			data.Write(st.d.Raw())
		}
		if data.Len() >= l {
			break
		}
	}
	// 読み込む
	i, err = data.Read(buf)
	st.buf = data.Bytes()
	return
}
