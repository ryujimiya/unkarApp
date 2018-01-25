package unview

// ビューモデル

import (
	"../conf"
	"../get2ch"
	"../model"
	"../util"
	"bytes"
	"html"
	"io"
	"net/http"
	"strconv"
	"strings"
	"text/template"
)

const (
	LAST_N_DEFAULT_NUM       = 100
	ANCHOR_REFER_DEFAULT     = 1
	ANCHOR_SEE_DEFAULT       = 8
	TYPE_ASSERTION_ERROR_MSG = "型アサーションに失敗しました。どうしようも無いので処理を中断します。"
)

type ViewContainer struct {
	ViewComponent
	CanonicalStr   string
	DescriptionStr string
	HostUrl        string
	Ret            bytes.Buffer
}

type HtmlStartData struct {
	SPFlag         bool
	Title          string
	Canonical      bool
	CanonicalUrl   string
	DescriptionStr string
}

type HtmlEndData struct {
	SPFlag bool
}

type HeaderFooter struct {
	SPFlag       bool
	HostUrl      string
	CanonicalUrl string
	NichUrl      string
	Ver          string
	ByteSize     int64
	StatusCode   int
}

type ThreadHeaderData struct {
	SPFlag       bool
	HostUrl      string
	Boardname    string
	CanonicalUrl string
	Server       string
	Board        string
	Thread       string
}

type ThreadFooterData struct {
	SPFlag     bool
	ByteSize   int64
	Boardname  string
	HostUrl    string
	Board      string
	StatusCode int
	Ver        string
}

var HtmlStartTempl = template.Must(template.ParseFiles("undity/template/html_start.templ"))
var HtmlEndTempl = template.Must(template.ParseFiles("undity/template/html_end.templ"))
var ServerHeaderTempl = template.Must(template.ParseFiles("undity/template/server_header.templ"))
var ServerFooterTempl = template.Must(template.ParseFiles("undity/template/server_footer.templ"))
var BoardHeaderTempl = template.Must(template.ParseFiles("undity/template/board_header.templ"))
var BoardFooterTempl = template.Must(template.ParseFiles("undity/template/board_footer.templ"))
var ThreadHeaderTempl = template.Must(template.ParseFiles("undity/template/thread_header.templ"))
var ThreadFooterTempl = template.Must(template.ParseFiles("undity/template/thread_footer.templ"))
var ThreadDummyResBbspink = []string{"unkarからのお知らせ", "unkar", "時間とID", "dat落ちしたBBSPINKは表示できません。", ""}

func NewViewContainer(path string, r *http.Request) unutil.View {
	vc := &ViewContainer{
		ViewComponent: NewViewComponent(path, r),
	}
	vc.HostUrl = vc.GetHostUrl()
	// 128kbyteのバッファを先に確保しておく
	vc.Ret.Grow(1024 * 128)
	return vc
}

func (vc *ViewContainer) PrintData(model unutil.Model) unutil.Output {
	mod := model.GetMod()
	vc.Model = model

	// データ取得
	vc.NowData = model.GetData()

	vc.SetCommonCanonical()

	switch model.GetClassName() {
	case unmodel.ClassNameThread:
		// スレッド
		// ユーザーキャッシュを利用する
		unutil.CheckNotModified(vc.R, mod)
		vc.threadPathAnalyze()
	case unmodel.ClassNameBoard:
		// スレッド一覧
		// ユーザーキャッシュを利用する
		unutil.CheckNotModified(vc.R, mod)
		vc.boardPathAnalyze()
	case unmodel.ClassNameServer:
		// 板一覧
		vc.serverHeaderPrint()
		vc.serverPrint()
		vc.serverFooterPrint()
	case unmodel.ClassNameSpecial:
		// 特殊画面
		if match := unconf.RegInitSpecialSearch.FindStringSubmatch(vc.Path); match != nil {
			unutil.MovedPermanently("http://unkar.org/search?q=" + match[1])
		} else {
			vc.serverHeaderPrint() // serverのやつを流用
			vc.specialPrint()
			vc.serverFooterPrint() // serverのやつを流用
		}
	case unmodel.ClassNameNone:
		fallthrough
	default:
		vc.serverHeaderPrint() // serverのやつを流用
		vc.notFoundPrint()
		vc.serverFooterPrint() // serverのやつを流用
	}

	// タイトル生成
	vc.Title += " - unkar"
	// headerの出力
	vc.OutputHeader(mod)
	// HTMLの出力
	vc.printHtml()

	return vc.Output
}

func (vc *ViewContainer) printHtml() {
	start := &bytes.Buffer{}
	HtmlStartTempl.Execute(start, &HtmlStartData{
		SPFlag:         vc.IsSP,
		Title:          vc.Title,
		Canonical:      vc.Canonical,
		CanonicalUrl:   vc.GetSiteUrl() + vc.GetAppName() + vc.CanonicalStr,
		DescriptionStr: vc.DescriptionStr,
	})
	HtmlEndTempl.Execute(&vc.Ret, &HtmlEndData{
		SPFlag: vc.IsSP,
	})
	// 出力をまとめる
	vc.Reader = io.MultiReader(start, &vc.Ret)
}

func (vc *ViewContainer) boardPathAnalyze() {
	attr := ""
	if match := unconf.RegInitBoard.FindStringSubmatch(vc.Path); match != nil {
		board := match[1]
		err := vc.Model.GetError()
		if match[3] != "" {
			attr = match[3]
		}
		// 正規化されたURL
		vc.CanonicalStr = "/" + vc.Model.GetUrl()
		cstr := "/" + board + match[2]

		// ヘッダーの生成
		vc.boardHeaderPrint(board)
		if vc.Model.Is404() || vc.Path != cstr {
			// 404専用処理
			vc.notFoundPrint()
		} else if err != nil {
			// エラーが発生していた場合
			vc.Title = err.Error()
			vc.Ret.WriteString(`<div class="topbox"><h1 class="pagetitle">` + vc.Title + `</h1></div>`)
			vc.nowThreadLinkPrint()
		} else {
			vc.boardPrint(attr)
		}
		// フッターの生成
		vc.boardFooterPrint(board)
	}
}

func (vc *ViewContainer) threadPathAnalyze() {
	attr := ""
	if match := unconf.RegInitThread.FindStringSubmatch(vc.Path); match != nil {
		board := match[1]
		thread := match[2]
		if len(match) > 3 && match[3] != "" {
			attr = match[3]
		}
		err := vc.Model.GetError()
		// 正規化されたURL
		vc.CanonicalStr = "/" + vc.Model.GetUrl()

		// ヘッダーの生成
		vc.threadHeaderPrint(board, thread)
		if vc.Model.Is404() {
			// 404とする
			vc.notFoundPrint()
		} else if err != nil {
			// エラーが発生していた場合
			vc.Title = err.Error()
			vc.Ret.WriteString(`<div class="topbox"><h1 class="pagetitle">` + vc.Title + `</h1></div>`)
			vc.nowThreadLinkPrint()
		} else if attr != "" {
			vc.threadPathAttr(attr)
		} else {
			vc.threadPrint(attr)
		}
		// フッターの生成
		vc.threadFooterPrint(board)
	}
}

func (vc *ViewContainer) threadPathAttr(attr string) {
	if m := unconf.RegThreadAttrLastn.FindStringSubmatch(attr); m != nil {
		vc.lastnThreadPrint(m[1])
	} else if m := unconf.RegThreadAttrBottom.FindStringSubmatch(attr); m != nil {
		vc.bottomThreadPrint(m[1])
	} else if m := unconf.RegThreadAttrTop.FindStringSubmatch(attr); m != nil {
		vc.topThreadPrint(m[1])
	} else if m := unconf.RegThreadAttrRes.FindStringSubmatch(attr); m != nil {
		vc.resThreadPrint(m[1])
	} else if m := unconf.RegThreadAttrId.FindStringSubmatch(attr); m != nil {
		vc.searchIdThreadPrint(m[1])
	} else if m := unconf.RegThreadAttrLink.FindStringSubmatch(attr); m != nil {
		vc.linkThreadPrint(m[1])
	} else if m := unconf.RegThreadAttrTree.FindStringSubmatch(attr); m != nil {
		vc.treeThreadPrint(m[1])
	} else if m := unconf.RegThreadAttrAnchor.FindStringSubmatch(attr); m != nil {
		vc.anchorThreadPrint(m[1])
	} else if attr == "/" {
		vc.threadPrint(attr)
	} else {
		// スペルミス等を転送するのは良いらしい
		unutil.MovedPermanently("http://unkar.org/r" + vc.CanonicalStr)
	}
}

func (vc *ViewContainer) setServerCanonical() {
	vc.CanonicalStr = ""
	if vc.Path != "" {
		vc.Canonical = true
	}
}

func (vc *ViewContainer) serverPrint() {
	host_url := vc.HostUrl
	vc.Title = vc.Model.GetTitle()
	title := vc.Title
	data, ok := vc.NowData.(*[]unmodel.ServerItem)
	if !ok {
		unutil.InternalServerError(vc.R, TYPE_ASSERTION_ERROR_MSG)
	}

	// 正規化
	vc.setServerCanonical()

	vc.Ret.WriteString(`<div class="topbox"><h1 class="pagetitle">` + title + "</h1></div>\n")
	vc.Ret.WriteString("<div id=\"content\" class=\"normal\">\n")

	for _, it := range *data {
		vc.Ret.WriteString("<div class=\"cate-area\">\n<h2 class=\"subtitle\">" + it.Name + "</h2>\n<ul>\n")
		for _, list := range it.ItaList {
			vc.Ret.WriteString("<li><a href=\"" + host_url + "/" + list.Url + "\">" + list.Name + "</a></li>\n")
		}
		vc.Ret.WriteString("</ul>\n</div>\n")
	}
	vc.Ret.WriteString("</div>\n")
}

func (vc *ViewContainer) boardPrint(attr string) {
	var thread_title string

	vc.NowData = vc.Model.GetData()
	sort, dir := unmodel.BoardSort(vc.Model, attr)
	list, ok := vc.NowData.(*unmodel.ThreadItems)
	if !ok {
		unutil.InternalServerError(vc.R, TYPE_ASSERTION_ERROR_MSG)
	}
	vc.Title = vc.Model.GetTitle()
	title := vc.Title
	board := vc.Model.GetUrl()
	url_board := vc.HostUrl + "/" + board + "/"
	vc.Ret.WriteString(`<div class="topbox"><h1 class="pagetitle">` + title + "</h1></div>\n")
	vc.Ret.WriteString(`<div id="content">` + "\n")
	vc.Ret.WriteString(`<div class="board-sort">` + "\n")
	vc.Ret.WriteString(`<span class="board-sort-info">並び替え：</span>` + "\n")
	if attr == "" {
		vc.Ret.WriteString(`<span class="board-sort-no"><a href="` + url_board + sort[3] + `" rel="nofollow">No` + dir[3] + `</a></span>` + "\n")
	} else {
		vc.Ret.WriteString(`<span class="board-sort-no"><a href="` + vc.HostUrl + "/" + board + `" title="` + template.HTMLEscapeString(title) + `">No` + dir[3] + `</a></span>` + "\n")
	}
	vc.Ret.WriteString(`<span class="board-sort-res"><a href="` + url_board + sort[2] + `" rel="nofollow">レス` + dir[2] + `</a></span>` + "\n")
	vc.Ret.WriteString(`<span class="board-sort-speed"><a href="` + url_board + sort[0] + `" rel="nofollow">勢い` + dir[0] + `</a></span>` + "\n")
	vc.Ret.WriteString(`<span class="board-sort-since"><a href="` + url_board + sort[1] + `" rel="nofollow">時間` + dir[1] + `</a></span>` + "\n")
	vc.Ret.WriteString(`</div>` + "\n")
	vc.Ret.WriteString(`<ul id="board">` + "\n")
	for _, it := range *list {
		index := strings.LastIndex(it.Thread, " ")
		if index > 2 {
			thread_title = it.Thread[:index]
		} else {
			thread_title = it.Thread
		}
		vc.Ret.WriteString(`<li class="board-line"><div class="board-date">`)
		vc.Ret.WriteString(`<div class="board-since">` + it.Since + `</div>`)
		vc.Ret.WriteString(`<div class="board-res">(` + strconv.Itoa(int(it.Res)) + `)</div>`)
		vc.Ret.WriteString(`<div class="board-speed">勢い ` + strconv.Itoa(int(it.Spd)) + `</div>`)
		vc.Ret.WriteString(`</div><div class="board-data"><div class="board-thread">`)
		vc.Ret.WriteString(`<div class="board-line-no">` + strconv.Itoa(it.Num) + `</div>`)
		vc.Ret.WriteString(`<a href="` + url_board + strconv.Itoa(int(it.Sin)) + `" class="board-thread-link">` + thread_title + `</a>`)
		vc.Ret.WriteString(`</div></div></li>` + "\n")
	}
	vc.Ret.WriteString("</ul>\n</div>\n")
}

func (vc *ViewContainer) threadPrint(attr string) {
	thread_url := vc.HostUrl + "/" + vc.Model.GetUrl()
	vc.Title = vc.Model.GetTitle()
	title := vc.Title
	data, ok := vc.NowData.(*unmodel.ThreadData)
	if !ok {
		unutil.InternalServerError(vc.R, TYPE_ASSERTION_ERROR_MSG)
	}

	vc.Canonical = attr != ""
	vc.setDescription(data, 1)
	ct := createThreadLink(data, thread_url, 0, 0)
	vc.Ret.WriteString("<div class=\"topbox\"><h1 class=\"pagetitle\">" + title + "</h1></div>\n")
	vc.Ret.Write(ct)
	vc.Ret.WriteString("\n")
	vc.Ret.WriteString("<dl id=\"content\">\n")

	// １を表示する
	afi := vc.GetAffiliate()
	if data.Pink && data.DatFall {
		// bbspink＆dat落ちの場合
		vc.threadAnchorStyle(data, 1, afi)
		vc.threadAnchorStyleDummy(data, 2, ThreadDummyResBbspink)
	} else {
		size := len(data.Res) + 1
		vc.threadAnchorStyle(data, 1, afi)
		for i := 2; i < size; i++ {
			vc.threadAnchorStyle(data, i, "")
		}
	}

	vc.Ret.WriteString("</dl>\n")
	vc.Ret.Write(ct)
	vc.Ret.WriteString("\n")
	vc.nowThreadLinkPrint()
}

func (vc *ViewContainer) lastnThreadPrint(num_last string) {
	data, ok := vc.NowData.(*unmodel.ThreadData)
	if !ok {
		unutil.InternalServerError(vc.R, TYPE_ASSERTION_ERROR_MSG)
	}
	thread_url := vc.HostUrl + "/" + vc.Model.GetUrl()
	title := vc.Model.GetTitle()
	size := len(data.Res) - 1
	var ank_list []int
	var cflag bool
	nlnum, _ := strconv.Atoi(num_last)
	if nlnum <= 0 {
		nlnum = LAST_N_DEFAULT_NUM
		num_last = strconv.Itoa(LAST_N_DEFAULT_NUM)
		cflag = true
	}
	if size > nlnum {
		ank_list = unutil.Range(size-nlnum+1, size, 1)
	} else {
		ank_list = unutil.Range(1, size, 1)
	}
	l := len(ank_list)
	if l >= 1 && ank_list[0] != 1 {
		tmp := make([]int, l+1)
		tmp[0] = 1
		copy(tmp[1:], ank_list)
		ank_list = tmp
	}
	vc.Canonical = cflag
	vc.CanonicalStr += "/l" + num_last
	vc.Title = title + "|レス抽出(後方" + strconv.Itoa(l) + "件)"
	vc.contentsPrint(data, ank_list, thread_url, 0, 0)
}

func (vc *ViewContainer) topThreadPrint(num_end string) {
	data, ok := vc.NowData.(*unmodel.ThreadData)
	if !ok {
		unutil.InternalServerError(vc.R, TYPE_ASSERTION_ERROR_MSG)
	}
	thread_url := vc.HostUrl + "/" + vc.Model.GetUrl()
	title := vc.Model.GetTitle()
	size := len(data.Res) - 1
	var ank_list []int
	var cflag bool
	nenum, _ := strconv.Atoi(num_end)
	if nenum <= 0 {
		// 404とする
		vc.Code = http.StatusNotFound
	} else if nenum < size {
		ank_list = unutil.Range(1, nenum, 1)
	} else {
		ank_list = unutil.Range(1, size, 1)
		cflag = true
	}
	if cflag {
		vc.Canonical = true
	} else {
		vc.CanonicalStr += "/-" + num_end
	}
	vc.Title = title + "|レス抽出(～" + num_end + ")"
	vc.contentsPrint(data, ank_list, thread_url, 0, 0)
}

func (vc *ViewContainer) bottomThreadPrint(num_start string) {
	data, ok := vc.NowData.(*unmodel.ThreadData)
	if !ok {
		unutil.InternalServerError(vc.R, TYPE_ASSERTION_ERROR_MSG)
	}
	thread_url := vc.HostUrl + "/" + vc.Model.GetUrl()
	title := vc.Model.GetTitle()
	size := len(data.Res) - 1
	var ank_list []int
	var cflag bool
	nsnum, _ := strconv.Atoi(num_start)
	if nsnum <= 1 {
		ank_list = unutil.Range(1, size, 1)
		cflag = true
	} else if size > nsnum {
		ank_list = unutil.Range(nsnum, size, 1)
	} else {
		ank_list = unutil.Range(1, size, 1)
		cflag = true
	}

	l := len(ank_list)
	if l > 0 && ank_list[0] != 1 {
		tmp := make([]int, l+1)
		tmp[0] = 1
		copy(tmp[1:], ank_list)
		ank_list = tmp
	}
	if cflag {
		vc.Canonical = true
	} else {
		vc.CanonicalStr += "/" + num_start + "-"
	}
	vc.Title = title + "|レス抽出(" + num_start + "～)"
	vc.contentsPrint(data, ank_list, thread_url, 0, 0)
}

func (vc *ViewContainer) resThreadPrint(res_string string) {
	data, ok := vc.NowData.(*unmodel.ThreadData)
	if !ok {
		unutil.InternalServerError(vc.R, TYPE_ASSERTION_ERROR_MSG)
	}
	thread_url := vc.HostUrl + "/" + vc.Model.GetUrl()
	title := vc.Model.GetTitle()

	ank_list := AnalyzeResString(res_string, len(data.Res)+1)
	length := len(ank_list)
	res_canonical := GetCanonicalString(ank_list)
	subtitle := "|レス抽出(" + strconv.Itoa(length) + "件 &gt;&gt;" + res_canonical + ")"
	vc.Title = title + subtitle

	start := 0
	end := 0
	if length > 1 {
		start = ank_list[0]
		end = ank_list[length-1]
	} else if length == 0 {
		// 404とする
		vc.Code = http.StatusNotFound
	}
	vc.Canonical = res_string != res_canonical
	vc.CanonicalStr += "/" + res_canonical
	vc.contentsPrint(data, ank_list, thread_url, start, end)
}

func (vc *ViewContainer) linkThreadPrint(filter string) {
	data, ok := vc.NowData.(*unmodel.ThreadData)
	if !ok {
		unutil.InternalServerError(vc.R, TYPE_ASSERTION_ERROR_MSG)
	}
	thread_url := vc.HostUrl + "/" + vc.Model.GetUrl()
	title := vc.Model.GetTitle()

	ank_list := vc.GetLinkList(data, filter)
	length := len(ank_list)
	subtitle := "|リンク抽出(" + strconv.Itoa(length) + "件 絞り込み[" + filter + "])"
	vc.Title = title + subtitle
	if length == 0 {
		// 404とする
		vc.Code = http.StatusNotFound
	}
	vc.CanonicalStr += "/Link:" + filter
	vc.contentsPrint(data, ank_list, thread_url, 0, 0)
}

func (vc *ViewContainer) searchIdThreadPrint(id_string string) {
	data, ok := vc.NowData.(*unmodel.ThreadData)
	if !ok {
		unutil.InternalServerError(vc.R, TYPE_ASSERTION_ERROR_MSG)
	}
	thread_url := vc.HostUrl + "/" + vc.Model.GetUrl()
	title := vc.Model.GetTitle()

	var idlist []int
	if idl, ok2 := data.Id[id_string]; ok2 {
		idlist = idl
	} else {
		// 404とする
		vc.Code = http.StatusNotFound
	}
	subtitle := "|ID抽出(" + strconv.Itoa(len(idlist)) + "件 ID:" + id_string + ")"
	vc.CanonicalStr += "/ID:" + id_string
	vc.Title = title + subtitle
	vc.contentsPrint(data, idlist, thread_url, 0, 0)
}

func (vc *ViewContainer) treeThreadPrint(res_string string) {
	data, ok := vc.NowData.(*unmodel.ThreadData)
	if !ok {
		unutil.InternalServerError(vc.R, TYPE_ASSERTION_ERROR_MSG)
	}
	thread_url := vc.HostUrl + "/" + vc.Model.GetUrl()
	title := vc.Model.GetTitle()
	canonical_tree := ""

	var ank_list []int
	if strings.Index(res_string, "ID:") >= 0 {
		if idl, ok2 := data.Id[res_string[3:]]; ok2 {
			ank_list = idl
			canonical_tree = res_string
		} else {
			ank_list = []int{}
		}
	} else if strings.Index(res_string, "Link:") >= 0 {
		ank_list = vc.GetLinkList(data, res_string[5:])
		canonical_tree = res_string
	} else {
		ank_list = AnalyzeResString(res_string, len(data.Res)+1)
	}
	// ツリー構造の解析
	tree := vc.AnalyzeTree(ank_list)

	res_count := len(tree)
	subtitle := ""
	// 正規化する
	if res_count > 0 {
		if canonical_tree == "" {
			canonical_tree = strconv.Itoa(tree[0])
		}
		vc.Canonical = res_string != canonical_tree
		subtitle = "|ツリー抽出(" + strconv.Itoa(res_count) + "件 &gt;&gt;" + canonical_tree + ")"
	} else {
		canonical_tree = res_string
		subtitle = "|ツリー抽出(0件)"
		// 404とする
		vc.Code = http.StatusNotFound
	}
	vc.CanonicalStr += "/Tree:" + canonical_tree
	vc.Title = title + subtitle
	vc.contentsPrint(data, tree, thread_url, 0, 0)
}

func (vc *ViewContainer) anchorThreadPrint(res_string string) {
	data, ok := vc.NowData.(*unmodel.ThreadData)
	if !ok {
		unutil.InternalServerError(vc.R, TYPE_ASSERTION_ERROR_MSG)
	}
	thread_url := vc.HostUrl + "/" + vc.Model.GetUrl()
	title := vc.Model.GetTitle()
	canonical_tree := res_string

	var ref int
	var see int
	if strings.Index(res_string, "@") == 0 {
		list := strings.Split(res_string[1:], "!")
		ref, _ = strconv.Atoi(list[0])
		see, _ = strconv.Atoi(list[1])
	} else {
		if res_string != "Default" {
			canonical_tree = "Default"
		}
		ref = ANCHOR_REFER_DEFAULT
		see = ANCHOR_SEE_DEFAULT
	}
	// ツリー構造の解析
	tree := vc.AnchorTree(ref, see)

	res_count := len(tree)
	subtitle := ""
	// 正規化する
	if res_count > 0 {
		vc.Canonical = res_string != canonical_tree
		subtitle = "|アンカー抽出(" + strconv.Itoa(res_count) + "件 絞り込み[被参照" + strconv.Itoa(ref) + "以上、参照" + strconv.Itoa(see) + "以下])"
	} else {
		subtitle = "|アンカー抽出(0件)"
		if see <= 0 {
			vc.Code = http.StatusNotFound
		}
	}
	vc.CanonicalStr += "/Anchor:" + canonical_tree
	vc.Title = title + subtitle
	vc.contentsPrint(data, tree, thread_url, 0, 0)
}

// $listはarrayかSplFixedArray
func (vc *ViewContainer) contentsPrint(data *unmodel.ThreadData, list []int, thread_url string, start, end int) {
	ct := createThreadLink(data, thread_url, start, end)
	vc.Ret.WriteString("<div class=\"topbox\">\n")
	vc.Ret.WriteString(vc.GetAffiliate())
	vc.Ret.WriteString("</div>\n")
	vc.Ret.WriteString(`<h1 class="pagetitle">` + vc.Title + "</h1>\n")
	vc.Ret.Write(ct)
	vc.Ret.WriteString("\n<dl id=\"content\">\n")

	if data.Pink && data.DatFall {
		// bbspink＆dat落ちの場合
		vc.setDescription(data, 1)
		vc.threadAnchorStyle(data, 1, "")
		vc.threadAnchorStyleDummy(data, 2, ThreadDummyResBbspink)
	} else if len(list) > 0 {
		vc.setDescription(data, list[0])
		for _, it := range list {
			vc.threadAnchorStyle(data, it, "")
		}
	}

	vc.Ret.WriteString("</dl>\n")
	vc.Ret.Write(ct)
	vc.Ret.WriteString("\n")
	vc.nowThreadLinkPrint()
}

func (vc *ViewContainer) setDescription(data *unmodel.ThreadData, i int) {
	if len(data.Res) > i {
		db := unutil.StripTags([]byte(data.Res[i].Data[3]), nil) // タグ消去
		db = unconf.RegUrl.ReplaceAllLiteral(db, []byte{})       // url消去
		db = unconf.RegSpace.ReplaceAllLiteral(db, []byte{' '})  // スペースをまとめる
		ds := string(db)
		desc := unutil.Utf8Substr(ds, 150) // 切り落とし

		if len(ds) != len(desc) {
			desc += "..."
		}
		vc.DescriptionStr = html.UnescapeString(desc)
	}
}

func (vc *ViewContainer) specialPrint() {
	vc.Code = http.StatusNotFound

	vc.Title = vc.Model.GetTitle()
	vc.Ret.WriteString(`<div class="topbox"><h1 class="pagetitle">` + vc.Title + `</h1></div>`)
	vc.Ret.WriteString("<div id=\"content\" class=\"normal\">\n")
	vc.Ret.WriteString(`うんかーJSモードで使用している特殊なURLです。`)
	vc.Ret.WriteString("</div>\n")
	vc.nowThreadLinkPrint()
}

func (vc *ViewContainer) notFoundPrint() {
	vc.Code = http.StatusNotFound

	vc.Title = gone_title
	vc.Ret.WriteString(`<div class="topbox"><h1 class="pagetitle">` + http.StatusText(vc.Code) + `</h1></div>`)
	vc.Ret.WriteString("<div id=\"content\" class=\"normal\">\n")
	vc.Ret.WriteString(gone_title)
	vc.Ret.WriteString("</div>\n")
	vc.nowThreadLinkPrint()
}

func (vc *ViewContainer) threadAnchorStyle(data *unmodel.ThreadData, i int, afi string) {
	if i >= len(data.Res) {
		return
	}
	aname := vc.GetAppName()
	istr := strconv.Itoa(i)
	opt := data.Res[i].Opt
	res := data.Res[i].Data
	path := data.Path
	var link string
	var name string
	var idlink string
	var cls string
	var mini string

	if opt != nil {
		id := opt[1]
		idmap, ok := data.Id[id]

		idlink = opt[0]
		if ok {
			idlen := len(idmap)
			if idlen <= 1 {
				idlink += ` ID:` + id
			} else {
				if idlen <= 1 {
					cls = `thread-tan`
				} else if idlen < 5 {
					cls = `thread-ninki`
				} else {
					cls = `thread-makka`
				}
				idlink += ` <a href="/` + aname + `/` + path + `/ID:` + id + `" data-id="` + id + `" class="` + cls + `" rel="nofollow">ID:</a>` + id
			}
			idlink += opt[2]
		} else {
			idlink = res[2]
		}
	} else {
		idlink = res[2]
	}
	if tr, ok := data.Tree[i]; ok {
		if len(tr) <= 3 {
			cls = "thread-ninki"
		} else {
			cls = "thread-makka"
		}
		link = `<a href="/` + aname + `/` + path + `/Tree:` + istr + `" class="` + cls + `" rel="nofollow">` + istr + `</a>`
	} else {
		link = istr
	}
	if res[1] != "" {
		name = `<span class="thread-resdate-mail">` + res[0] + `</span> <span class="thread-resdate-thin">[` + res[1] + `]</span> `
	} else {
		name = `<span class="thread-resdate-name">` + res[0] + `</span> : `
	}
	vc.Ret.WriteString(`<dt id="l` + istr + `" class="thread-resdate">`)
	vc.Ret.WriteString(`<span class="thread-resdate-no">` + link + `</span>`)
	vc.Ret.WriteString(` : ` + name + `<span class="thread-resdate-info">` + idlink + `</span>`)
	vc.Ret.WriteString(`</dt>`)
	if len(res[3]) > 255 {
		mini = ` thread-mini`
	}
	vc.Ret.WriteString(`<dd id="b` + istr + `" class="thread-data` + mini + `">`)
	vc.Ret.WriteString(res[3])
	if vc.Textbrowser {
		vc.Ret.WriteString(`<br><br>`)
	}
	if afi != "" {
		vc.Ret.WriteString(`<div>` + afi + `</div>`)
	}
	vc.Ret.WriteString("</dd>\n")
}

func (vc *ViewContainer) serverHeaderPrint() {
	ServerHeaderTempl.Execute(&vc.Ret, &HeaderFooter{
		SPFlag:       vc.IsSP,
		CanonicalUrl: vc.GetSiteUrl() + vc.GetAppName(),
	})
}

func (vc *ViewContainer) serverFooterPrint() {
	ServerFooterTempl.Execute(&vc.Ret, &HeaderFooter{
		SPFlag:  vc.IsSP,
		HostUrl: vc.HostUrl,
		Ver:     unconf.Ver,
	})
}

func (vc *ViewContainer) boardHeaderPrint(board string) {
	server := vc.Model.GetServer()

	BoardHeaderTempl.Execute(&vc.Ret, &HeaderFooter{
		SPFlag:       vc.IsSP,
		HostUrl:      vc.HostUrl,
		CanonicalUrl: vc.GetSiteUrl() + vc.GetAppName() + vc.CanonicalStr,
		NichUrl:      "http://" + server + "/" + board + "/",
	})
}

func (vc *ViewContainer) boardFooterPrint(board string) {
	data := HeaderFooter{
		SPFlag:  vc.IsSP,
		HostUrl: vc.HostUrl,
		Ver:     unconf.Ver,
	}
	if vc.Model != nil {
		data.ByteSize = vc.Model.GetByteSize() / 1000
		data.StatusCode = vc.Model.GetCode()
	}
	BoardFooterTempl.Execute(&vc.Ret, &data)
}

func (vc *ViewContainer) threadHeaderPrint(board, thread string) {
	data, tdok := vc.NowData.(*unmodel.ThreadData)
	if !tdok {
		unutil.InternalServerError(vc.R, TYPE_ASSERTION_ERROR_MSG)
	}
	thd := ThreadHeaderData{
		SPFlag:       vc.IsSP,
		HostUrl:      vc.HostUrl,
		Boardname:    data.Boardname,
		CanonicalUrl: vc.GetSiteUrl() + vc.GetAppName() + vc.CanonicalStr,
		Server:       data.Server,
		Board:        board,
		Thread:       thread,
	}
	ThreadHeaderTempl.Execute(&vc.Ret, &thd)
}

func (vc *ViewContainer) threadFooterPrint(board string) {
	data, tdok := vc.NowData.(*unmodel.ThreadData)
	if !tdok {
		unutil.InternalServerError(vc.R, TYPE_ASSERTION_ERROR_MSG)
	}
	tfd := ThreadFooterData{
		SPFlag:    vc.IsSP,
		Boardname: data.Boardname,
		HostUrl:   vc.HostUrl,
		Board:     board,
		Ver:       unconf.Ver,
	}
	if vc.Model != nil {
		tfd.ByteSize = vc.Model.GetByteSize() / 1000
		tfd.StatusCode = vc.Model.GetCode()
	}
	ThreadFooterTempl.Execute(&vc.Ret, &tfd)
}

func (vc *ViewContainer) nowThreadLinkPrint() {
	host_url := vc.HostUrl
	vc.Ret.WriteString("<aside>\n")
	vc.Ret.WriteString(vc.GetAffiliate() + "\n")
	vc.Ret.WriteString("<div id=\"search\">\n")
	vc.Ret.WriteString("<span class=\"subtitle\">【スレッド検索】</span>\n")
	vc.Ret.WriteString(vc.GetGoogleSearchFrom(60) + "\n")
	vc.Ret.WriteString("</div>\n")
	vc.Ret.WriteString(vc.GetAffiliate_300x250() + "\n")
	list := get2ch.GetViewThreadList(1, 10)
	if len(list) > 0 {
		vc.Ret.WriteString("<div id=\"hip\">\n")
		vc.Ret.WriteString("<span class=\"subtitle\">【最近見られたスレッド】</span>\n<table>\n")
		for _, it := range list {
			bnshort := it.Boardname
			base := host_url + "/" + it.Board
			vc.Ret.WriteString("<tr>")
			if bnshort != "" {
				vc.Ret.WriteString(`<td><a href="` + base + `/` + it.Thread + `">` + it.Title + `</a></td><td><a href="` + base + `">` + bnshort + `</a></td>`)
			} else {
				vc.Ret.WriteString(`<td><a href="` + base + `/` + it.Thread + `">` + it.Title + `</a></td><td></td>`)
			}
			vc.Ret.WriteString("</tr>\n")
		}
		vc.Ret.WriteString("</table>\n</div>\n")
	}
	vc.Ret.WriteString("</aside>\n<br>\n")
}

func (vc *ViewContainer) threadAnchorStyleDummy(data *unmodel.ThreadData, no int, res []string) {
	dummy := unmodel.CreateThreadDataDummy(data, no)
	copy(dummy.Res[no].Data[:], res)
	vc.threadAnchorStyle(&dummy, no, "")
}
