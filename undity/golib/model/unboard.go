package unmodel

// ボードモデル

import (
	"../get2ch"
	"../util"
	"bufio"
	"bytes"
	"regexp"
	"sort"
	"strconv"
	"time"
)

type ThreadItem struct {
	Num    int
	Res    int64
	Thread string
	Sin    int64
	Since  string
	Spd    int64
}
type ThreadItems []*ThreadItem
type ThreadItemsBySince struct {
	ThreadItems
}
type ThreadItemsBySpeed struct {
	ThreadItems
}
type ThreadItemsByRes struct {
	ThreadItems
}
type ThreadItemsByNumber struct {
	ThreadItems
}

type Board struct {
	ModelComponent
	data      ThreadItems
	sorttitle string
}

var RegsLine = regexp.MustCompile(`(\d+)\.dat<>(.*\s\((\d+)\))`)

func NewBoard(host string, path []string) *Board {
	model := &Board{
		ModelComponent: CreateModelComponent(ClassNameBoard, host),
	}
	model.url = path[1]
	model.g2ch = get2ch.NewGet2ch(model.url, "")
	model.analyzeData()
	return model
}

func (bd *Board) GetData() interface{} { return &bd.data }
func (bd *Board) GetTitle() (ret string) {
	ret = bd.title
	if bd.sorttitle != "" {
		ret += " | " + bd.sorttitle
	}
	return
}

func (bd *Board) analyzeData() {
	data := make(ThreadItems, 0, 1024)
	// データの取得
	scanner := bufio.NewScanner(unutil.ShiftJISToUtf8Reader(bytes.NewReader(bd.g2ch.GetData())))
	if bd.g2ch.GetError() != nil {
		bd.err = bd.g2ch.GetError()
		return
	}
	bd.mod = bd.g2ch.GetModified()
	bd.title = bd.g2ch.GetBoardName()
	nowtime := time.Now().Unix()

	i := 1
	for scanner.Scan() {
		if line := RegsLine.FindStringSubmatch(scanner.Text()); line != nil {
			sin, _ := strconv.ParseInt(line[1], 10, 64)
			res, _ := strconv.ParseInt(line[3], 10, 64)
			if res == 0 {
				// ゼロ除算対策
				res = 1
			}
			tmp := (nowtime - sin) / res
			if tmp == 0 {
				// ゼロ除算対策
				tmp = 1
			}
			data = append(data, &ThreadItem{
				Num:    i,
				Res:    res,
				Thread: line[2],
				Sin:    sin,
				Since:  time.Unix(sin, 0).Format("2006/01/02 15:04"),
				Spd:    86400 / tmp,
			})
		} else {
			data = append(data, &ThreadItem{
				Num:    i,
				Res:    0,
				Thread: "スレッドが壊れているみたい",
				Sin:    0,
				Since:  "故障",
				Spd:    0,
			})
		}
		i++
	}
	bd.data = data[:i-1 : i-1]
}

func BoardSort(model unutil.Model, attr string) ([]string, []string) {
	var inter sort.Interface
	var st string
	s := []string{"sp", "si", "re", "no"}
	dir := []string{"", "", "", ""}
	bd, ok := model.(*Board)
	if !ok {
		return s, dir
	}
	list := bd.data

	switch attr {
	case "sp":
		inter = &ThreadItemsBySpeed{list}
		s[0] = "sp2"
		dir[0] = "↓"
		st = "勢い順↓"
	case "sp2":
		inter = sort.Reverse(&ThreadItemsBySpeed{list})
		dir[0] = "↑"
		st = "勢い順↑"
	case "si":
		inter = &ThreadItemsBySince{list}
		s[1] = "si2"
		dir[1] = "↓"
		st = "時間順↓"
	case "si2":
		inter = sort.Reverse(&ThreadItemsBySince{list})
		dir[1] = "↑"
		st = "時間順↑"
	case "re":
		inter = &ThreadItemsByRes{list}
		s[2] = "re2"
		dir[2] = "↓"
		st = "レス数順↓"
	case "re2":
		inter = sort.Reverse(&ThreadItemsByRes{list})
		dir[2] = "↑"
		st = "レス数順↑"
	case "no":
		inter = &ThreadItemsByNumber{list}
		s[3] = "no2"
		dir[3] = "↓"
		st = "番号順↓"
	case "no2":
		fallthrough
	default:
		dir[3] = "↑"
	}

	if inter != nil {
		sort.Sort(inter)
		bd.sorttitle = st
	}
	return s, dir
}

func (t ThreadItems) Len() int      { return len(t) }
func (t ThreadItems) Swap(i, j int) { t[i], t[j] = t[j], t[i] }
func (ts ThreadItemsBySince) Less(i, j int) bool {
	// 時間降順
	return ts.ThreadItems[i].Sin > ts.ThreadItems[j].Sin
}
func (ts ThreadItemsBySpeed) Less(i, j int) bool {
	// 勢い降順
	return ts.ThreadItems[i].Spd > ts.ThreadItems[j].Spd
}
func (ts ThreadItemsByRes) Less(i, j int) bool {
	// レス数降順
	return ts.ThreadItems[i].Res > ts.ThreadItems[j].Res
}
func (ts ThreadItemsByNumber) Less(i, j int) bool {
	// 番号降順
	return ts.ThreadItems[i].Num > ts.ThreadItems[j].Num
}
