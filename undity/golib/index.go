package untidy

import (
	"./get2ch"
	"./util"
	"bufio"
	"bytes"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"
)

type BoardItem struct {
	Path string
	Name string
}

type ServerItem struct {
	Cate  string
	Board []BoardItem
}

type IndexData struct {
	Title       string
	Desc        string
	Afi         string
	AfiFlag     bool
	ItaList     *bytes.Buffer
	NowList     *bytes.Buffer
	SelectBoard *bytes.Buffer
	Year        int
}

const (
	AppName                = "r"
	DAT_FAILURE_HYSTERESIS = 30 * time.Second
)

var IndexPCTempl *template.Template = template.Must(template.ParseFiles("template/index_pc.templ"))

func IndexMain(r *http.Request, notfound bool) unutil.Output {
	out := unutil.Output{
		Header: http.Header{},
		ZFlag:  true,
	}
	buf := &bytes.Buffer{}
	now := time.Now()
	g2ch := get2ch.NewGet2ch("", "")

	boardmap := make(map[string]string)
	list := []ServerItem{}
	rc := g2ch.GetBBSmenu(false)
	scanner := bufio.NewScanner(rc)
	for scanner.Scan() {
		var dat, name string
		if tmp := strings.Split(scanner.Text()+"<>", "<>"); tmp != nil {
			dat = tmp[0]
			name = tmp[1]
		}
		l := len(list)
		if name == "" {
			if l > 0 && len(list[l-1].Board) == 0 {
				list[l-1].Cate = dat
			} else {
				list = append(list, ServerItem{
					Cate:  dat,
					Board: []BoardItem{},
				})
			}
		} else {
			var board string
			if tmp := strings.Split(dat, "/"); tmp != nil {
				board = tmp[1]
			}
			list[l-1].Board = append(list[l-1].Board, BoardItem{
				Path: board,
				Name: name,
			})
			boardmap[board] = name
		}
	}
	rc.Close()

	l := len(list)
	if l > 0 && len(list[l-1].Board) == 0 {
		list = list[:l-1]
	}
	list = append([]ServerItem{
		ServerItem{
			Cate: "人気",
			Board: []BoardItem{
				BoardItem{Path: "news4vip", Name: boardmap["news4vip"]},
				BoardItem{Path: "livejupiter", Name: boardmap["livejupiter"]},
				BoardItem{Path: "poverty", Name: boardmap["poverty"]},
				BoardItem{Path: "news", Name: boardmap["news"]},
				BoardItem{Path: "morningcoffee", Name: boardmap["morningcoffee"]},
				BoardItem{Path: "newsplus", Name: boardmap["newsplus"]},
				BoardItem{Path: "mnewsplus", Name: boardmap["mnewsplus"]},
				BoardItem{Path: "akb", Name: boardmap["akb"]},
			},
		},
	}, list...)

	out.Header.Set("Content-Type", "text/html; charset=utf-8")
	data := IndexData{}
	data.Title = "unkar - 2ちゃんねるの閲覧と検索"
	data.Year = now.Year()

	if notfound {
		out.Code = http.StatusNotFound
		data.Title = "ページが見つかりません - 404 Not Found ::unkar"
	} else {
		out.Code = http.StatusOK
	}

	if notfound {
		data.Desc = "<h2>ページが見つかりません。</h2><p>お目当てのページは以下から見つかるかもしれません。</p>"
		data.AfiFlag = true
		data.Afi = unutil.AffiliateMicroad_300x250
	} else {
		data.Desc = "<h2>2ちゃんねるを検索</h2>"
		data.AfiFlag = false
	}
	bl := bytes.Buffer{}
	for _, cate := range list {
		bl.WriteString("<div class=\"cate-area\">\n<h3>" + cate.Cate + "</h3>\n<ul>\n")
		c := 0
		for _, it := range cate.Board {
			if c == 8 {
				bl.WriteString("</ul>\n<p class=\"showbutton\">[全部表示]</p>\n<ul class=\"italist-hide\">\n")
			}
			bl.WriteString(`<li><a href="/` + AppName + `/` + it.Path + `">` + it.Name + "</a></li>\n")
			c++
		}
		bl.WriteString("</ul>\n</div>\n")
	}
	data.ItaList = &bl
	data.NowList = nowList(now, 70)
	data.SelectBoard = printSelectBoard(list[1:])
	IndexPCTempl.Execute(buf, &data)

	out.Reader = buf
	return out
}

func printSelectBoard(si []ServerItem) *bytes.Buffer {
	buf := &bytes.Buffer{}
	buf.WriteString(`<select name="board" id="search-select-board">` + "\n")
	buf.WriteString(`<optgroup label="説明">` + "\n")
	buf.WriteString(`<option value="">板の絞り込み</option>` + "\n")
	for _, s := range si {
		buf.WriteString(`</optgroup>` + "\n" + `<optgroup label="` + s.Cate + `">` + "\n")
		for _, it := range s.Board {
			buf.WriteString(`<option value="` + it.Path + `">` + it.Name + `</option>` + "\n")
		}
	}
	buf.WriteString("</optgroup>\n</select>\n")
	return buf
}

func nowList(mod time.Time, l int) *bytes.Buffer {
	nw := &bytes.Buffer{}
	compmod := mod.Add(DAT_FAILURE_HYSTERESIS)
	for _, it := range get2ch.GetViewThreadList(0, l) {
		sin, _ := strconv.ParseInt(it.Thread, 10, 64)
		nw.WriteString(`<dt><a href="/` + AppName + `/` + it.Board + `/` + it.Thread + `">` + it.Title + ` (` + strconv.Itoa(int(it.Res)) + ")</a></dt>\n")
		nw.WriteString(`<dd>` + unutil.CreateDateString(time.Unix(sin, 0)))
		if it.Boardname != "" {
			nw.WriteString(`　<a href="/` + AppName + `/` + it.Board + `">` + it.Boardname + `</a>`)
		}
		if it.Lastdate.After(compmod) {
			nw.WriteString(`　(dat落ち)`)
		}
		nw.WriteString("</dd>\n")
	}
	return nw
}
