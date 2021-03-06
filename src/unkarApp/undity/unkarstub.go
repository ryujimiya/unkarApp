package unkarstub

import (
	"./golib/conf"
	"./golib/get2ch"
	"./golib/model"
	"./golib/util"
	"./golib/view"
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

/**
 *  Unkarのスタブ
 */

////////////////////////////////////////////////////////////
// 初期化
////////////////////////////////////////////////////////////
/**
 * Unkarの初期化処理
 */
func InitUnkar() {
	// モデルの初期化処理
	initModel()
}

/**
 * tmpディレクトリ
 */
func GetTmpHtmlDir() string {
	// 実行ファイルパスを取得
	exePath, _ := os.Executable()
	exeDir := filepath.Dir(exePath)
	// 表示用htmlファイル格納ディレクトリ
	tmpDir := exeDir + "\\tmp"
	return tmpDir
}

////////////////////////////////////////////////////////////
// index.goの抜粋
////////////////////////////////////////////////////////////
/**
  Unkarの板アイテム
*/
type BoardItem struct {
	Path string
	Name string
}

/**
  Unkarのサーバーアイテム
*/
type ServerItem struct {
	Cate  string
	Board []BoardItem
}

/**
  UnkarのIndexMainの抜粋
  @return サーバーアイテムの一覧
*/
func UnkarIndexMain() []ServerItem {
	// Get2chを生成する
	g2ch := get2ch.NewGet2ch("", "")

	// 板のマップを生成
	boardmap := make(map[string]string)
	// サーバーアイテム
	list := []ServerItem{}
	// 板一覧を取得する(BBSmenu)
	rc := g2ch.GetBBSmenu(false)

	// パース処理:板マップに格納する
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

	return list
}

////////////////////////////////////////////////////////////
// Controller
////////////////////////////////////////////////////////////
/**
 * Unkarのルーティング(モデルと正規表現の対応)
 */
type Route struct {
	Regs  *regexp.Regexp
	Model func(host string, path []string) unutil.Model
}

/**
 * モデル一覧
 */
var modelList []*Route

/**
 * モデルの初期化処理
 */
func initModel() {
	modelList = []*Route{
		&Route{
			Regs: unconf.RegInitThread,
			Model: func(host string, path []string) unutil.Model {
				return unmodel.NewThread(host, path)
			},
		},
		&Route{
			Regs: unconf.RegInitBoard,
			Model: func(host string, path []string) unutil.Model {
				return unmodel.NewBoard(host, path)
			},
		},
		&Route{
			Regs: unconf.RegInitSpecial,
			Model: func(host string, path []string) unutil.Model {
				return unmodel.NewSpecial(host, path)
			},
		},
		&Route{
			Regs: unconf.RegInitServer,
			Model: func(host string, path []string) unutil.Model {
				return unmodel.NewServer(host, path)
			},
		},
	}
}

/**
 * モデルを取得する
 * @param host ホストのurl
 *        path パス
 * @return unutilのModel
 */
func getModel(host string, path string) unutil.Model {
	var model unutil.Model
	flag := false

	for _, it := range modelList {
		if match := it.Regs.FindStringSubmatch(path); match != nil {
			model = it.Model(host, match)
			flag = true
			break
		}
	}
	if flag == false {
		model = unmodel.NewNone(host, []string{path})
	}
	return model
}

////////////////////////////////////////////////////////////
// 板モデル
////////////////////////////////////////////////////////////
/**
 * 板モデルを取得する
 * @return unutilのモデル
 */
func GetBoardModel(boardName string) unutil.Model {
	// ホスト
	host := "/r"
	// パス
	path := "/" + boardName

	// unutilのモデル
	// モデルの取得
	unutilModel := getModel(host, path)

	return unutilModel
}

////////////////////////////////////////////////////////////
// スレッドモデル
////////////////////////////////////////////////////////////
/**
 * スレッドモデルを取得する
 * @return unutilのモデル
 */
func GetThreadModel(boardName string, threadNo int64, attr string) unutil.Model {
	// ホスト
	host := "/r"
	// パス
	path := "/" + boardName + "/" + fmt.Sprintf("%d", threadNo)
	if attr != "" {
		path += "/" + attr
	}

	// unutilのモデル
	// モデルの取得
	unutilModel := getModel(host, path)

	return unutilModel
}

/**
  UnkarのViewの出力を取得する
  @param host ホストURL
  @param path パス
  @param model モデル
  @return HTML出力文字列
*/
func getViewOutput(host string, path string, model unutil.Model) string {
	outputStr := ""

	// dummy
	r, _ := http.NewRequest("GET", host+path, nil)

	viewfunc := unview.NewViewContainer
	view := viewfunc(path, r)
	//host := view.GetHostUrl()
	output := view.PrintData(model)

	bufReader := bufio.NewReader(output.Reader)
	for {
		line, _, err := bufReader.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
			break
		}
		outputStr = outputStr + string(line)
	}

	return outputStr
}

/**
  Unkarの板Viewの出力を取得する
  @param boardName 板名
  @param threadNo スレッド番号
  @param model モデル
  @return HTML出力文字列
*/
func GetBoardViewOutput(boardName string, threadNo int64, attr string, model unutil.Model) string {
	// ホスト
	host := "/r"
	// パス
	path := "/" + boardName + "/" + fmt.Sprintf("%d", threadNo)
	if attr != "" {
		path += "/" + attr
	}

	return getViewOutput(host, path, model)
}
