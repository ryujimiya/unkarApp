package unview

// ビューモデル

import (
	"../conf"
	"../model"
	"../util"
	"bytes"
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

type ViewComponent struct {
	unutil.Output
	R           *http.Request
	Model       unutil.Model
	NowData     interface{}
	Textbrowser bool
	IsSP        bool
	Canonical   bool
	Title       string
	Path        string
}

const (
	site_url            = "http://unkar.org/" // 設置URL
	base_host           = "/"                 // 設置パス
	folder              = "/2ch/dat"          // dat格納フォルダ
	app_name            = "r"                 // 動作ファイル名
	app_js_name         = "read.html#!"
	itaname_path        = "/2ch/dat/ita.data" // 板情報格納ファイル
	affiliate_amazon_id = "unkar-22"
	gone_title          = "ページが存在しません"
)

var RegsResString = regexp.MustCompile(`(\d{1,4})(?:\-(\d{1,4}))?`)

func NewViewComponent(path string, r *http.Request) ViewComponent {
	ua := r.Header.Get("User-Agent")
	vc := ViewComponent{
		Output: unutil.Output{
			Code:   http.StatusOK,
			Header: http.Header{},
			ZFlag:  true,
		},
		R:           r,
		Path:        path,
		Textbrowser: strings.Contains(ua, "w3m"),
		IsSP:        unutil.IsMobile(r),
	}
	return vc
}

func (vc *ViewComponent) GetCode() int {
	return vc.Code
}

func (vc *ViewComponent) GetHeader() http.Header {
	return vc.Header
}

func (vc *ViewComponent) GetSiteUrl() string {
	return site_url
}

func (vc *ViewComponent) GetBaseUrl() string {
	return base_host
}

func (vc *ViewComponent) GetFolder() string {
	return folder
}

func (vc *ViewComponent) GetAppName() string {
	return app_name
}

func (vc *ViewComponent) GetAppJsName() string {
	return app_js_name
}

func (vc *ViewComponent) GetItanamePath() string {
	return itaname_path
}

func (vc *ViewComponent) GetAffiliate() string {
	var afi string
	if vc.IsSP {
		afi = vc.GetAffiliate_300x250()
	} else {
		afi = vc.GetAffiliate_728x90()
	}
	return afi
}

func (vc *ViewComponent) GetAffiliate_728x90() string {
	return unutil.AffiliateMicroad_728x90
}

func (vc *ViewComponent) GetAffiliate_300x250() string {
	return unutil.AffiliateMicroad_300x250
}

func (vc *ViewComponent) GetAmazonId() string {
	return affiliate_amazon_id
}

func (vc *ViewComponent) GetHtmlHead() string {
	return unutil.DefaultHtmlHead
}

func (vc *ViewComponent) GetSearchForm(size int) string {
	if size <= 0 {
		size = 12
	}
	text := "検索"
	return text
}

func (vc *ViewComponent) GetGoogleSearchFrom(size int) string {
	if size <= 0 {
		size = 16
	}
	text := "検索"
	return text
}

func (vc *ViewComponent) GetHostUrl() string {
	return vc.GetBaseUrl() + vc.GetAppName()
}

func (vc *ViewComponent) AnchorTree(rc, tc int) []int {
	th, ok := vc.NowData.(*unmodel.ThreadData)
	if !ok {
		return nil
	}

	size := len(th.Res)
	retdata := make([]int, 0, size)
	retdata_once := make(map[int]bool, size)
	check := func(index int) (ret bool) {
		if anc, ok := th.Anchor[index]; ok {
			// 一定数以下の安価
			ret = len(anc) <= tc
		} else {
			// 安価を送っていない
			ret = true
		}
		return
	}
	var sub func(index int, tree []int)
	sub = func(index int, tree []int) {
		if _, ok := retdata_once[index]; ok {
			// 登録済み
			return
		}
		if len(tree) < rc {
			// 一定未満の安価
			return
		}
		if check(index) == false {
			// 一定数超過の安価
			return
		}
		// 未登録
		retdata = append(retdata, index)
		retdata_once[index] = true
		for _, it := range tree {
			// tree直系
			if it > index && check(it) {
				// 一定数超過の安価と未来のレスは防ぐ
				if tr, ok := th.Tree[it]; ok {
					// 新しいtree
					sub(it, tr)
				}
				if _, ok := retdata_once[it]; !ok {
					retdata = append(retdata, it)
					retdata_once[it] = true
				}
			}
		}
		return
	}

	for i := 1; i < size; i++ {
		if _, ok := retdata_once[i]; !ok {
			sub(i, th.Tree[i])
		}
	}
	l := len(retdata)
	return retdata[:l:l]
}

func (vc *ViewComponent) AnalyzeTree(list []int) []int {
	th, ok := vc.NowData.(*unmodel.ThreadData)
	if !ok {
		return nil
	}

	size := len(th.Res)
	retdata := make([]int, 0, size)
	retdata_once := make(map[int]bool, size)
	list = vc.searchAnchor(list, size)
	retdata = vc.reslistFunc(list, retdata, retdata_once, size)
	return unutil.UniqueIntSlice(retdata)
}

func (vc *ViewComponent) reslistFunc(reslist, retdata []int, retdata_once map[int]bool, size int) []int {
	th, ok := vc.NowData.(*unmodel.ThreadData)
	if !ok {
		return nil
	}

	for _, it := range reslist {
		if _, flg := retdata_once[it]; !flg && (it < size) {
			retdata = append(retdata, it)
			retdata_once[it] = true
			if tr, f := th.Tree[it]; f {
				retdata = vc.reslistFunc(tr, retdata, retdata_once, size)
			}
			if an, f := th.Anchor[it]; f {
				retdata = vc.reslistFunc(an, retdata, retdata_once, size)
			}
		}
	}
	return retdata
}

func (vc *ViewComponent) searchAnchor(list []int, size int) []int {
	retdata := make([]int, 0, size)
	retdata_once := make(map[int]bool, size)
	for _, it := range list {
		if _, ok := retdata_once[it]; !ok && it < size {
			if r := vc.anchorFunc(retdata_once, it, size); r > 0 {
				retdata = append(retdata, r)
			}
		}
	}
	return unutil.UniqueIntSlice(retdata)
}

func (vc *ViewComponent) anchorFunc(retdata_once map[int]bool, resno, size int) int {
	th, ok := vc.NowData.(*unmodel.ThreadData)
	if !ok {
		return 0
	}
	if _, ok := retdata_once[resno]; ok {
		return resno
	}
	retdata_once[resno] = true
	an, ok := th.Anchor[resno]
	if !ok {
		return resno
	}
	for _, it := range an {
		if _, flg := retdata_once[it]; !flg && it < size {
			if r := vc.anchorFunc(retdata_once, it, size); r > 0 {
				resno = unutil.MinInt(resno, r)
			}
		}
	}
	return resno
}

func (vc *ViewComponent) SetCommonCanonical() {
	if vc.R.URL.RawQuery != "" {
		vc.Canonical = true
	}
}

func (vc *ViewComponent) GetLinkList(thread *unmodel.ThreadData, filter string) []int {
	var link map[int]bool
	if thread.Link.All == nil {
		return []int{}
	}
	if filter == "All" {
		link = thread.Link.All
	} else if filter == "Thread" {
		link = thread.Link.Thread
	} else if filter == "Image" {
		link = thread.Link.Image
	} else if filter == "Movie" {
		link = thread.Link.Movie
	} else if filter == "Archive" {
		link = thread.Link.Archive
	} else {
		link = thread.Link.All
	}

	link_list := make([]int, len(link))
	index := 0
	for key, _ := range link {
		link_list[index] = key
		index++
	}
	sort.Ints(link_list)
	return link_list
}

func createThreadLink(data *unmodel.ThreadData, url string, start, end int) []byte {
	alllen := len(data.Res) - 1
	txt := bytes.Buffer{}
	size := unutil.MinInt(alllen, 1000)
	start_start := start
	end_end := end

	// 4kbyteのバッファーを先に確保しておく
	txt.Grow(1024 * 4)
	txt.WriteString("<div class=\"nav\">\n<ul class=\"social-button\">\n")
	txt.WriteString("<li class=\"social-button-item\">　</li>" + "\n")
	txt.WriteString("<li class=\"social-button-item\">　</li>" + "\n")
	txt.WriteString("<li class=\"social-button-item\">　</li>" + "\n")
	txt.WriteString("<li class=\"social-button-item\">　</li>" + "\n")
	txt.WriteString("</ul>\n<div class=\"reslist\">\n[" + strconv.Itoa(alllen) + "res] ")
	if start != 0 && end != 0 {
		end_end += 100
		if end_end >= size {
			end_end = size - 1
		}
		start_start -= 100
		if start_start < 1 {
			start_start = 1
		}
		txt.WriteString("<a href=\"" + url + "\">全部</a> ")
		end++
		if end < end_end {
			txt.WriteString(fmt.Sprintf("<a href=\"%s/%d-%d\" rel=\"nofollow\">次100</a> ", url, end, end_end))
		}
		start--
		if start > start_start {
			txt.WriteString(fmt.Sprintf("<a href=\"%s/%d-%d\" rel=\"nofollow\">前100</a> ", url, start_start, start))
		}
		for i := 0; i < size; i += 100 {
			txt.WriteString(fmt.Sprintf("<a href=\"%s/%d-%d\" rel=\"nofollow\">%d-</a> ", url, i+1, i+100, i+1))
		}
		txt.WriteString("<a href=\"" + url + "/l50\" rel=\"nofollow\">最新50</a>")
	} else {
		txt.WriteString("<a href=\"" + url + "\">全部</a> ")
		for i := 0; i < size; i += 100 {
			txt.WriteString(fmt.Sprintf("<a href=\"%s/%d-%d\" rel=\"nofollow\">%d-</a> ", url, i+1, i+100, i+1))
		}
		txt.WriteString("<a href=\"" + url + "/l50\" rel=\"nofollow\">最新50</a>")
	}
	txt.WriteString("</div>\n")
	if (len(data.Link.All) > 0) || (len(data.Tree) > 0) {
		link := data.Link
		txt.WriteString("<div class=\"linklist\">抽出 : ")
		txt.WriteString("<a href=\"" + url + "/Anchor:@0!3\" rel=\"nofollow\">まとめ（仮）</a> ")
		if len(data.Tree) > 0 {
			txt.WriteString("<a href=\"" + url + "/Anchor:Default\" rel=\"nofollow\">アンカー</a> ")
		}
		if len(link.All) > 0 {
			txt.WriteString("<a href=\"" + url + "/Link:All\" rel=\"nofollow\">URL</a>")
			txt.WriteString("[<a href=\"" + url + "/Tree:Link:All\" rel=\"nofollow\">＋</a>] ")
		}
		if len(link.Thread) > 0 {
			txt.WriteString("<a href=\"" + url + "/Link:Thread\" rel=\"nofollow\">スレッド</a>")
			txt.WriteString("[<a href=\"" + url + "/Tree:Link:Thread\" rel=\"nofollow\">＋</a>] ")
		}
		if len(link.Image) > 0 {
			txt.WriteString("<a href=\"" + url + "/Link:Image\" rel=\"nofollow\">画像</a>")
			txt.WriteString("[<a href=\"" + url + "/Tree:Link:Image\" rel=\"nofollow\">＋</a>] ")
		}
		if len(link.Movie) > 0 {
			txt.WriteString("<a href=\"" + url + "/Link:Movie\" rel=\"nofollow\">動画</a>")
			txt.WriteString("[<a href=\"" + url + "/Tree:Link:Movie\" rel=\"nofollow\">＋</a>] ")
		}
		if len(link.Archive) > 0 {
			txt.WriteString("<a href=\"" + url + "/Link:Archive\" rel=\"nofollow\">書庫</a>")
			txt.WriteString("[<a href=\"" + url + "/Tree:Link:Archive\" rel=\"nofollow\">＋</a>] ")
		}
		txt.WriteString("</div>\n")
	}
	prefix := len(app_name) + 2
	l := strings.Index(url[prefix:], "/")
	var bname string
	if data.Boardname != "" {
		bname = " title=\"" + data.Boardname + "\""
	}
	txt.WriteString("<div class=\"optionlist\">ナビゲーション : ")
	txt.WriteString("<a href=\"" + url[:prefix+l] + "\"" + bname + ">板に戻る</a>")
	txt.WriteString("</div>\n")
	txt.WriteString("</div>\n")
	return txt.Bytes()
}

func (vc *ViewComponent) OutputHeader(mod time.Time) {
	// utf-8の文字コードである事をヘッダーで明示する
	vc.Header.Set("Content-Type", "text/html; charset=utf-8")

	if mod.IsZero() == false {
		req_time := time.Now()
		if req_time.Before(mod.Add(-1 * 172800 * time.Second)) {
			// 2日引いても現在時刻より大きい時間だった場合
			vc.Header.Set("Expires", unutil.CreateModString(req_time.Add(unconf.OneYearSec)))
			// 現在時刻を最終更新時刻とする
			vc.Header.Set("Last-Modified", unutil.CreateModString(req_time))
		} else {
			// 最終更新時刻送付
			vc.Header.Set("Last-Modified", unutil.CreateModString(mod))
		}
		// ETagを設定
		vc.Header.Set("ETag", unutil.CreateETag(mod))
	} else {
		//vc.Header.Set("Expires", "Mon, 26 Jul 1997 05:00:00 GMT")
	}
}

func AnalyzeResString(res_string string, size int) []int {
	filter := make(map[int]bool)
	sliceCallBack := func(matchstr string) string {
		match := RegsResString.FindStringSubmatch(matchstr)
		if match == nil {
			return matchstr
		}

		if len(match) > 2 && match[2] != "" {
			min, _ := strconv.Atoi(match[1])
			max, _ := strconv.Atoi(match[2])
			if min < max && min >= 1 {
				if max >= size {
					max = size - 1
				}
				for min <= max {
					if _, ok := filter[min]; ok {
						break
					}
					filter[min] = true
					min++
				}
			}
		} else {
			min, _ := strconv.Atoi(match[1])
			if min < size && min >= 1 {
				filter[min] = true
			}
		}
		return ""
	}
	// 検索処理
	RegsResString.ReplaceAllStringFunc(res_string, sliceCallBack)
	// キーを配列で取得
	ank_list := make([]int, len(filter))
	index := 0
	for key, _ := range filter {
		ank_list[index] = key
		index++
	}
	// ソート
	sort.Ints(ank_list)
	return ank_list
}

func GetCanonicalString(list []int) (canonical_str string) {
	l := len(list)
	if l == 0 {
		// 空だった場合
		return
	}
	canonical_list := make([]string, 0, l)
	value := list[0]
	list = list[1:]

	first := value      // 連番の開始値
	second := first + 1 // 期待値
	for _, value = range list {
		if value == second {
			// 期待値通りだったので期待値を更新
			second++
		} else {
			// 期待値ではない場合
			if first == (second - 1) {
				// 連番が始まる前に崩壊
				canonical_list = append(canonical_list, strconv.Itoa(first))
			} else {
				// 連番ができた
				canonical_list = append(canonical_list, fmt.Sprintf("%d-%d", first, second-1))
			}
			// 開始値と期待値を更新
			first = value
			second = value + 1
		}
	}
	if first == value {
		// 開始値を更新してすぐにループを抜けた場合
		canonical_list = append(canonical_list, strconv.Itoa(value))
	} else {
		canonical_list = append(canonical_list, fmt.Sprintf("%d-%d", first, value))
	}
	if len(canonical_list) > 0 {
		canonical_str = strings.Join(canonical_list, ",")
	}
	return
}
