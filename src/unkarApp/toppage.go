package main

import (
	"./undity"
	"fmt"
)

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

type TopPage struct {
	*walk.Composite
	mainWin        *MainWin        // メインウィンドウ
	listBoxBoard   *walk.ListBox   // 板一覧リストボックス
	boardListModel *BoardListModel // 板一覧モデル
	title          string          // タイトル
}

func newTopPage(parent walk.Container, mainWin *MainWin) (*TopPage, error) {
	// トップページの生成
	topPage := new(TopPage)

	topPage.mainWin = mainWin

	if err := (Composite{
		AssignTo: &topPage.Composite,
		Name:     "板一覧",
		Layout:   VBox{},
		Children: []Widget{
			ListBox{
				AssignTo:              &topPage.listBoxBoard,
				OnCurrentIndexChanged: topPage.listBoxBoardCurrentIndexChanged,
				OnItemActivated:       topPage.listBoxBoardItemActivated,
			},
		},
		Visible: false,
	}).Create(NewBuilder(parent)); err != nil {
		return nil, err
	}

	if err := walk.InitWrapperWindow(topPage); err != nil {
		return nil, err
	}

	topPage.UpdateContents()

	return topPage, nil
}

func (topPage *TopPage) Title() string {
	return topPage.title
}

func (topPage *TopPage) UpdateContents() {
	// モデルの生成
	topPage.boardListModel = NewBoardListModel()

	topPage.listBoxBoard.SetCurrentIndex(-1)
	// モデルを再設定する
	topPage.listBoxBoard.SetModel(topPage.boardListModel)

	topPage.title = ""
	topPage.mainWin.UpdateTitle(topPage)
}

/**
 * 板一覧リストボックス選択インデックスが変わった
 * @param なし
 * @return なし
 */
func (topPage *TopPage) listBoxBoardCurrentIndexChanged() {
	i := topPage.listBoxBoard.CurrentIndex()
	item := &topPage.boardListModel.items[i]

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
func (topPage *TopPage) listBoxBoardItemActivated() {
	i := topPage.listBoxBoard.CurrentIndex()
	item := &topPage.boardListModel.items[i]

	name := item.name
	value := item.value
	fmt.Println("listBoxBoardItemActivated")
	fmt.Println("name: ", name)
	fmt.Println("value: ", value)

	// 板ページを表示
	topPage.mainWin.NavigateToBoardPage(name, value, "sp")
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
	// Unkarの板一覧（全体）
	var boardListAll []unkarstub.BoardItem

	// Unkarのサーバー一覧
	// サーバー一覧を取得する
	serverList := unkarstub.UnkarIndexMain()

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
