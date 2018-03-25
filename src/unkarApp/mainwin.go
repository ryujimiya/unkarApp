package main

import (
	"./undity"
	//"./undity/golib/util"
	//"./undity/golib/model"
	"fmt"
	"log"
	"os"
)

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

const (
	Version = "1.0.0.1"
)

////////////////////////////////////////////////////////////
// MainWin
////////////////////////////////////////////////////////////
/**
 * メインウィンドウ
 */
type MainWin struct {
	// 派生元：Walkメインウィンドウ
	*walk.MainWindow
	// 板一覧リストボックス
	listBoxBoard *walk.ListBox
	// 板一覧モデル
	boardListModel *BoardListModel
}

/**
 * コンストラクタ
 * @param なし
 * @return (1)メインウィンドウ
 *         (2)エラー
 */
func NewMainWin() (*MainWin, error) {
	// tmpディレクトリ
	tmpDir := unkarstub.GetTmpHtmlDir()
	// 先ず前回起動時のディレクトリを削除
	if err := os.RemoveAll(tmpDir); err != nil {
		fmt.Println(err)
	}
	// ディレクトリ作成
	if err := os.Mkdir(tmpDir, 0777); err != nil {
		fmt.Println(err)
	}

	// Unkar初期化処理
	unkarstub.InitUnkar()

	// メインウィンドウ生成
	mainWin := new(MainWin)

	// モデルの生成
	mainWin.boardListModel = NewBoardListModel()

	// メインウィンドウのウィンドウ生成
	err := MainWindow{
		AssignTo: &mainWin.MainWindow,
		Title:    "Unkar App " + Version,
		MinSize:  Size{600, 400},
		Layout:   VBox{},
		Children: []Widget{
			ListBox{
				AssignTo: &mainWin.listBoxBoard,
				Model:    mainWin.boardListModel,
				OnCurrentIndexChanged: mainWin.listBoxBoardCurrentIndexChanged,
				OnItemActivated:       mainWin.listBoxBoardItemActivated,
			},
		},
	}.Create()

	// デフォルトのフォント(walk.Fontのinit関数参照)
	//font, err:= walk.NewFont("MS Shell Dlg 2", 8, 0x00)
	// フォントサイズを大きくする
	font, err := walk.NewFont("MS Shell Dlg 2", 12, 0x00)
	if err != nil {
		log.Fatal(err)
	}
	mainWin.listBoxBoard.SetFont(font)

	return mainWin, err
}

/**
 * 板一覧リストボックス選択インデックスが変わった
 * @param なし
 * @return なし
 */
func (mainWin *MainWin) listBoxBoardCurrentIndexChanged() {
	i := mainWin.listBoxBoard.CurrentIndex()
	item := &mainWin.boardListModel.items[i]

	name := item.name
	value := item.value
	fmt.Println("CurrentIndex: ", i)
	fmt.Println("name: ", name)
	fmt.Println("value: ", value)
}

/**
 * 板一覧リストボックスアイテムがダブルクリックされた
 * @param なし
 * @return なし
 */
func (mainWin *MainWin) listBoxBoardItemActivated() {
	i := mainWin.listBoxBoard.CurrentIndex()
	item := &mainWin.boardListModel.items[i]

	//name := item.name
	value := item.value
	//walk.MsgBox(mainWin, "Name", name, walk.MsgBoxIconInformation)
	//walk.MsgBox(mainWin, "Value", value, walk.MsgBoxIconInformation)

	// 板ウィンドウの生成
	boardWin, err := NewBoardWin(mainWin, value)
	if err != nil {
		log.Fatal(err)
	}
	// 表示
	boardWin.Run()
}

////////////////////////////////////////////////////////////
// BoardListItem
////////////////////////////////////////////////////////////
/**
 * リストボックスアイテム
 */
type BoardListItem struct {
	name  string
	value string
}

////////////////////////////////////////////////////////////
// BoardListModel
////////////////////////////////////////////////////////////
/**
 * リストボックスモデル
 */
type BoardListModel struct {
	// 派生元：リストモデルベース
	walk.ListModelBase
	// アイテム一覧
	items []BoardListItem
}

/**
 * コンストラクタ
 */
func NewBoardListModel() *BoardListModel {
	// Unkarのサーバー一覧
	var serverList []unkarstub.ServerItem
	// Unkarの板一覧（全体）
	var boardListAll []unkarstub.BoardItem

	// サーバー一覧を取得する
	serverList = unkarstub.UnkarIndexMain()
	// ボード一覧
	//boardListAll = make([]unkarstub.BoardItem, 0)

	for _, server := range serverList {
		// サーバーの板一覧の取得
		boardList := server.Board
		for _, board := range boardList {
			// 板を板一覧（全体）に追加
			boardListAll = append(boardListAll, board)
		}
	}

	// リストボックスのモデルを生成
	model := &BoardListModel{items: make([]BoardListItem, len(boardListAll))}
	for i, board := range boardListAll {
		// アイテム名(表示名)
		name := board.Name
		// アイテムの値
		value := board.Path
		// リストボックスモデルにリストボックスアイテムをセットする
		model.items[i] = BoardListItem{name, value}
	}

	return model
}

/**
 * アイテム数を取得する
 * @return アイテム数
 */
func (model *BoardListModel) ItemCount() int {
	return len(model.items)
}

/**
 * アイテムの値を取得する
 * @return アイテムの値
 */
func (model *BoardListModel) Value(index int) interface{} {
	return model.items[index].name
}
