package unmodel

// スレッドモデル

import (
	"../conf"
	"../get2ch"
	"../util"
	"github.com/PuerkitoBio/goquery"
	//"bufio"
	"bytes"
	"fmt" // DEBUG
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	CORRUPT_DAT_STRING = "壊れています"
)

type LinkData struct {
	All     map[int]bool
	Thread  map[int]bool
	Image   map[int]bool
	Movie   map[int]bool
	Archive map[int]bool
}

type ResItem struct {
	Data [5]string
	Opt  []string
}

type ThreadData struct {
	Pink      bool
	DatFall   bool
	Server    string
	Boardname string
	Path      string
	Res       []ResItem
	Tree      map[int][]int
	Anchor    map[int][]int
	Id        map[string][]int
	Link      LinkData
}

type Thread struct {
	ModelComponent
	board  string
	option string
	data   ThreadData
}

var RegsName = regexp.MustCompile(`<b>([\s\t]*)</b>`)
var RegsHttp = regexp.MustCompile(`(s?h?ttps?)://((?:[\-\.a-zA-Z0-9]+\/?)(?:[-_.!~*"()a-zA-Z0-9;/?:\@&=+\$,%#\|]+))`)
var RegsBe = regexp.MustCompile(` BE:(\d+)\-([^\s\(]+\(\d+\))`)
var RegsIdsplit = regexp.MustCompile(`^(.*) ID:([\w!\+/]+)(.*)`)
var RegsRes = regexp.MustCompile(`(&gt;(?:&gt;)?)(\d+)([-,\d]*)`)
var RegsSssp = regexp.MustCompile(`sssp(\://img\.2ch\.net/([-_\w\./?&]+))`)

var RegsC2ch = regexp.MustCompile(`^c\.2ch\.net/test/\-/(\w+)(?:/(\d{9,10}))?(/[-,l\d]+)?`)
var RegsThread = regexp.MustCompile(`^(\w+\.2ch\.net|\w+\.bbspink\.com)/test/read\.\w+[/#](\w+)/(\d{9,10})/?([-,l\d]+)?`)
var RegsBoard = regexp.MustCompile(`^(\w+\.2ch\.net|\w+\.bbspink\.com)/(\w+)/?$`)
var RegsUnkar = regexp.MustCompile(`^www\.unkar\.org/read/(?:\w+\.2ch\.net|\w+\.bbspink\.com)/(\w+)(?:/(\d{9,10}))?`)
var RegsImage = regexp.MustCompile(`\.(?:(?:tif?|gi)f|jp(?:eg?|g)|p(?:ng|sd)|a(?:rt|i)|bmp|ico)$`)
var RegsMovie = regexp.MustCompile(`(?:nico(?:video\.jp/watch/(?:[ns]m|co|lv)|\.ms/(?:l(?:/co|v)|[ns]m))\d+|youtu(?:be\.(?:co(?:\.jp|m)|jp)/watch|\.be/[\-\w]+)|xvideos\.(?:com|jp)/video)`)
var RegsArchive = regexp.MustCompile(`\.(?:t(?:ar|gz)|[7gx]z|bz2|cab|lzh|rar|zip|Z)$`)
var RegsYoutube = regexp.MustCompile(`^(?:www\.)?youtube\.(?:com|jp|co\.jp)/watch(?:_videos)?\?.*v(?:ideo_ids)?=([\-\w]+)`)
var RegsYoutubeMin = regexp.MustCompile(`^youtu\.be/([\-\w]+)`)

// <img src="//img.5ch.net/ico/aka.gif"/> <br/> 俺「助六買ってきて。」 (再生成後のHTML)
var RegsImgTag5chIcon = regexp.MustCompile(`<img src="(//img.5ch.net/ico/[^"]+)"/>`)
var Regs5chIconUrl = regexp.MustCompile(`//img.5ch.net/ico/`)
var RegsImgurBlank = regexp.MustCompile(`(i[\s]*m[\s]*g[\s]*u[\s]*r[\s]*\.com)`)
var RegsHtmlTag = regexp.MustCompile(`<[^>]+>`)

var CorruptDatStringList = [5]string{
	CORRUPT_DAT_STRING,
	CORRUPT_DAT_STRING,
	CORRUPT_DAT_STRING,
	CORRUPT_DAT_STRING,
	CORRUPT_DAT_STRING,
}

func NewThread(host string, path []string) *Thread {
	model := &Thread{ModelComponent: CreateModelComponent(ClassNameThread, host)}
	model.board = path[1]
	if len(path) > 3 {
		model.option = path[3]
	}
	model.data.Path = model.board + "/" + path[2]
	model.url = model.data.Path
	model.g2ch = get2ch.NewGet2ch(model.board, path[2])
	model.analyzeData()
	return model
}

func (th *Thread) GetData() interface{} { return &th.data }

func CreateThreadDataDummy(data *ThreadData, size int) ThreadData {
	return ThreadData{
		Pink:      data.Pink,
		DatFall:   data.DatFall,
		Server:    data.Server,
		Boardname: data.Boardname,
		Path:      data.Path,
		Res:       make([]ResItem, size+1),
		Tree:      map[int][]int{},
		Anchor:    map[int][]int{},
		Id:        map[string][]int{},
		Link: LinkData{
			All:     map[int]bool{},
			Thread:  map[int]bool{},
			Image:   map[int]bool{},
			Movie:   map[int]bool{},
			Archive: map[int]bool{},
		},
	}
}

func linkChecker(match, m []string) string {
	if _, ok := unconf.ServerKill[m[1]]; ok {
		return hCheck(match[1], match[0])
	}
	if _, ok := unconf.BoardKill[m[2]]; ok {
		return hCheck(match[1], match[0])
	}
	return ""
}

func hAdd(ttp, url string) string {
	if !strings.Contains(ttp, "h") {
		return "h" + url
	}
	return url
}

func hCheck(ttp, matchurl string) string {
	return `<a href="` + hAdd(ttp, matchurl) + `" class="thread-data-link" target="_blank" rel="nofollow">` + matchurl + `</a>`
}

func createImage(ttp, matchurl string) string {
	return `<img src="` + hAdd(ttp, matchurl) + "\" style=\"max-width: 100%;\"><br/>"
}

func createYoutubeThumb(id string) string {
	link := `<div class="youtube-thumb">`
	for j := 1; j <= 3; j++ {
		jstr := strconv.Itoa(j)
		link += `<img src="https://img.youtube.com/vi/` + id + `/` + jstr + `.jpg" alt="YouTubeサムネイル[` + id + `] - ` + jstr + `枚目" width="120" height="90">`
	}
	link += `</div>`
	return link
}

func (th *Thread) analyzeData() {
	fmt.Printf("Thread,analyzeData\n") // DEBUG
	check := map[int]struct{}{}
	resNo := 1
	host_url := th.HostUrl
	//fmt.Printf("host_url" + host_url + "\r\n")
	thread_base_url := host_url + "/" + th.url
	treelist := map[int][]int{}
	anchorlist := map[int][]int{}
	idlist := map[string][]int{}
	linklist := map[int]bool{}
	threadlist := map[int]bool{}
	imagelist := map[int]bool{}
	movielist := map[int]bool{}
	archivelist := map[int]bool{}

	th.data.Server = th.g2ch.GetServer("")
	// 板名の設定
	th.data.Boardname = th.g2ch.GetBoardName()
	// bbspinkのフラグ設定
	th.data.Pink = strings.Contains(th.data.Server, "bbspink.com")
	// データの取得
	data := th.g2ch.GetData()
	if th.g2ch.GetError() != nil {
		th.err = th.g2ch.GetError()
		fmt.Printf("th.err=%+v\r\n", th.err)
		return
	}
	th.mod = th.g2ch.GetModified()
	// dat落ち確認
	th.data.DatFall = th.mod.After(time.Now())
	fmt.Println("th.data.DatFall=%+v\r\n", th.data.DatFall)

	ankerNumberColor := func(matchstr string) string {
		match := RegsRes.FindStringSubmatch(matchstr)
		if match == nil {
			return matchstr
		}

		anc, _ := strconv.Atoi(match[2])
		_, ok2 := check[anc]
		if tl, ok := treelist[anc]; ok {
			if !ok2 {
				// >>1>>1>>1等をカウントしてしまうのを防ぐ
				treelist[anc] = append(tl, resNo)
			}
		} else {
			treelist[anc] = append(make([]int, 0, 3), resNo)
		}
		if al, ok := anchorlist[resNo]; ok {
			if !ok2 {
				// >>1>>1>>1等をカウントしてしまうのを防ぐ
				anchorlist[resNo] = append(al, anc)
			}
		} else {
			anchorlist[resNo] = append(make([]int, 0, 3), anc)
		}

		check[anc] = struct{}{}
		url := match[2]
		if len(match) > 3 {
			url += match[3]
		}
		return `<a href="` + thread_base_url + `/` + url + `" rel="nofollow">` + matchstr + `</a>`
	}

	url_callback := func(matchstr string) string {
		match := RegsHttp.FindStringSubmatch(matchstr)
		if match == nil {
			return matchstr
		}
		linklist[resNo] = true

		if RegsImage.MatchString(match[2]) {
			// 画像ファイルっぽい
			link := ""
			if !Regs5chIconUrl.MatchString(match[0]) { // アイコンのときはリンクを表示しない
				link = hCheck(match[1], match[0])
			}
			img := createImage(match[1], match[0])
			//fmt.Printf("img=" + img + "\r\n")
			link += "<br>" + img
			imagelist[resNo] = true
			return link
		} else if m := RegsThread.FindStringSubmatch(match[2]); m != nil {
			// スレッドだった場合
			text := linkChecker(match, m)
			if text != "" {
				return text
			}
			threadlist[resNo] = true
			if len(m) > 4 {
				return `<a href="` + host_url + "/" + m[2] + "/" + m[3] + "/" + m[4] + `" class="thread-data-link">` + match[0] + "</a>"
			}
			return `<a href="` + host_url + "/" + m[2] + "/" + m[3] + `" class="thread-data-link">` + match[0] + "</a>"
		} else if m := RegsBoard.FindStringSubmatch(match[2]); m != nil {
			// 板だった場合
			text := linkChecker(match, m)
			if text != "" {
				return text
			}
			return `<a href="` + host_url + `/` + m[2] + `" class="thread-data-link">` + match[0] + `</a>`
		} else if m := RegsUnkar.FindStringSubmatch(match[2]); m != nil {
			// unkarReplaceAllLiteralString
			if len(m) > 2 {
				// スレッド
				return `<a href="` + host_url + `/` + m[1] + `/` + m[2] + `" class="thread-data-link">` + match[0] + `</a>`
			}
			// 板
			return `<a href="` + host_url + `/` + m[1] + `" class="thread-data-link">` + match[0] + `</a>`
		} else if m := RegsC2ch.FindStringSubmatch(match[2]); m != nil {
			if _, ok := unconf.BoardKill[m[1]]; ok {
				return hCheck(match[1], match[0])
			}
			ml := len(m)
			if ml > 2 {
				// スレッド
				threadlist[resNo] = true
				if ml > 3 {
					return `<a href="` + host_url + "/" + m[1] + "/" + m[2] + m[3] + `" class="thread-data-link">` + match[0] + "</a>"
				} else {
					return `<a href="` + host_url + "/" + m[1] + "/" + m[2] + `" class="thread-data-link">` + match[0] + "</a>"
				}
			}
			// 板
			return `<a href="` + host_url + "/" + m[1] + `" class="thread-data-link">` + match[0] + "</a>"
		} else if RegsArchive.MatchString(match[2]) {
			// アーカイブファイルっぽい
			link := hCheck(match[1], match[0])
			archivelist[resNo] = true
			return link
		} else if RegsMovie.MatchString(match[2]) {
			// 動画サイトのURLを含んでいるURLも引っかかる可能性があるので条件の優先度は最低にしておく
			movielist[resNo] = true
			link := hCheck(match[1], match[0])
			// youtube
			// サムネイルを付ける
			if m := RegsYoutube.FindStringSubmatch(match[2]); m != nil {
				link += createYoutubeThumb(m[1])
			} else if m := RegsYoutubeMin.FindStringSubmatch(match[2]); m != nil {
				link += createYoutubeThumb(m[1])
			}
			return link
		}
		return hCheck(match[1], match[0])
	}

	r := unutil.ShiftJISToUtf8Reader(bytes.NewReader(data))
	/*
		scanner := bufio.NewScanner(unutil.StripTagReader(r, get2ch.HtmlTag))
	*/

	/*2015-3-14 2ch新仕様対応で削除
	// datの行数取得
	l := unutil.MinInt(th.g2ch.NumLines(data)+1, 1100)
	reslist := make([]ResItem, l)
	for resNo < l && scanner.Scan() {
		check = map[int]struct{}{}
		// datの行の文字列取得
		line := scanner.Text()

		///////////////////////
		// 整形処理
		// URL処理
		line = RegsHttp.ReplaceAllStringFunc(line, url_callback)
		// 画像リンク処理
		line = RegsSssp.ReplaceAllString(line, `<img src="http${1}" alt="${2}">`)
		//
		line = RegsBe.ReplaceAllString(line, ` <a href="http://be.2ch.net/test/p.php?i=${1}" target="_blank" rel="nofollow">?${2}</a>`)
		// アンカー番号処理
		line = RegsRes.ReplaceAllStringFunc(line, ankerNumberColor)

		// レスアイテム生成
		tmpdata := ResItem{}
		// datの行を分解
		// ハンドルネーム<>メール<>日付時刻、ID<>本文<>タイトル
		tmpdataarray := strings.Split(line, "<>")
		if len(tmpdataarray) == 5 {
			copy(tmpdata.Data[:], tmpdataarray)
		} else {
			tmpdata.Data = CorruptDatStringList
		}
		if m := RegsIdsplit.FindStringSubmatch(tmpdata.Data[2]); m != nil {
			rdata := m[2]
			idl, ok := idlist[rdata]
			if !ok {
				idl = make([]int, 0, 3)
			}
			idlist[rdata] = append(idl, resNo)
			tmpdata.Opt = []string{m[1], rdata, m[3]}
		}
		tmpdata.Data[0] = RegsName.ReplaceAllString("<b>"+tmpdata.Data[0]+"</b>", "${1}")
		reslist[resNo] = tmpdata
		resNo++
	}
	if l <= 1 {
		// 文字化けなどで正常なデータがない場合は上書き
		reslist = make([]ResItem, 2)
		reslist[1].Data = CorruptDatStringList
	}
	*/

	reslist := make([]ResItem, 1) // 先頭要素は空にする
	threadTitleStr := ""
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		fmt.Print("url scraping failed\r\n")
	}
	doc.Find("title").Each(func(_ int, s *goquery.Selection) {
		threadTitleStr = s.Text()
	})
	fmt.Printf("threadTitleStr=" + threadTitleStr)
	doc.Find(".post").Each(func(_ int, s *goquery.Selection) {
		//threadNum := s.Find(".number").Text()
		handleName := s.Find(".name").Find("a").Text()
		if handleName == "" {
			handleName = s.Find(".name").Text()
		}
		emailStr, _ := s.Find(".name").Find("a").Attr("href")
		data1 := s.Find(".date").Text()
		data2 := s.Find(".uid").Text()
		dateStr := data1 + " " + data2
		//messageStr := s.Find(".message").Text()
		messageStr, _ := s.Find(".message").Html()
		// NOTE: このHTMLは再生成されており原文とはことなる 特に <br>が<br/> <img>が<img/>になるので注意
		messageStr = strings.Replace(messageStr, "\r", "", -1)
		messageStr = strings.Replace(messageStr, "\n", "", -1)
		messageStr = strings.Replace(messageStr, "<br/>", "\r\n", -1)
		messageStr = RegsImgTag5chIcon.ReplaceAllString(messageStr, `https:${1} `) // アイコン対応
		messageStr = RegsHtmlTag.ReplaceAllString(messageStr, "")                  // HTMLタグ削除
		messageStr = RegsImgurBlank.ReplaceAllString(messageStr, `imgur.com`)      // im gur.com対応

		//fmt.Printf("threadNum=" + threadNum + "\r\n")
		//fmt.Printf("handleName=" + handleName + "\r\n")
		//fmt.Printf("emailStr=" + emailStr + "\r\n")
		//fmt.Printf("dateStr=" + dateStr + "\r\n")
		//fmt.Printf("messageStr=" + messageStr + "\r\n")

		///////////////////////
		// 整形処理
		// 改行
		messageStr = strings.Replace(messageStr, "\r\n", "<br/>", -1)
		// URL処理
		messageStr = RegsHttp.ReplaceAllStringFunc(messageStr, url_callback)
		// 画像リンク処理
		messageStr = RegsSssp.ReplaceAllString(messageStr, `<img src="http${1}" alt="${2}">`)
		//
		messageStr = RegsBe.ReplaceAllString(messageStr, ` <a href="http://be.2ch.net/test/p.php?i=${1}" target="_blank" rel="nofollow">?${2}</a>`)
		// アンカー番号処理
		messageStr = RegsRes.ReplaceAllStringFunc(messageStr, ankerNumberColor)

		// ハンドルネーム<>メール<>日付時刻、ID<>本文<>タイトル (1の場合)
		// ハンドルネーム<>メール<>日付時刻、ID<>レス<> (2以降)
		tmpdataarray := make([]string, 5)
		tmpdataarray[0] = handleName
		tmpdataarray[1] = emailStr
		tmpdataarray[2] = dateStr
		tmpdataarray[3] = messageStr
		tmpdataarray[4] = ""
		if resNo == 1 {
			tmpdataarray[4] = threadTitleStr
		}
		tmpdata := ResItem{}
		copy(tmpdata.Data[:], tmpdataarray)
		if m := RegsIdsplit.FindStringSubmatch(tmpdata.Data[2]); m != nil {
			rdata := m[2]
			idl, ok := idlist[rdata]
			if !ok {
				idl = make([]int, 0, 3)
			}
			idlist[rdata] = append(idl, resNo)
			tmpdata.Opt = []string{m[1], rdata, m[3]}
		}

		reslist = append(reslist, tmpdata)

		resNo++
	})

	// タイトル設定
	th.title = string(unutil.StripTags([]byte(reslist[1].Data[4]), get2ch.HtmlTag))
	th.data.Res = reslist
	th.data.Tree = treelist
	th.data.Anchor = anchorlist
	th.data.Id = idlist
	th.data.Link.All = linklist
	th.data.Link.Thread = threadlist
	th.data.Link.Image = imagelist
	th.data.Link.Movie = movielist
	th.data.Link.Archive = archivelist
}
