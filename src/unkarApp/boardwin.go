package main

import (
	"./undity"
	"./undity/golib/model"
	//"./undity/golib/util"
	"errors"
	"fmt"
	"log"
)

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

//import (
//	"github.com/lxn/win"
//)

////////////////////////////////////////////////////////////
// BoardWin
////////////////////////////////////////////////////////////
/**
 * 板ウィンドウ
 */
type BoardWin struct {
	// 派生元：Walkメインウィンドウ
	*walk.MainWindow
	// スレッド一覧リストボックス
	listBoxThread *walk.ListBox
	// 板名
	boardName string
	// 板タイトル
	boardTitle string
	// スレッド一覧モデル
	threadListModel *ThreadListModel
}

/**
 * 板ウィンドウモデル
 */
type BoardWinModel struct {
	// ソート属性  (ラジオボタングループ用データバインディングのデータ構造体)
	//////////////////////////////////////////
	// "sp"
	// "勢い順↓"
	// "sp2":
	// "勢い順↑"
	// "si":
	// "時間順↓"
	// "si2":
	// "時間順↑"
	// "re":
	// "レス数順↓"
	// "re2":
	// "レス数順↑"
	// "":
	// "番号順↑"
	// "no":
	// "番号順↓"
	SortAttrStr string
}

/**
 * コンストラクタ
 * @param parentWin 親ウィンドウ
 * @param boardName URLに含まれる板名(povertyなど)
 * @return (1)板ウィンドウ
 *         (2)エラー
 */
func NewBoardWin(parentWin walk.Form, boardName string) (*BoardWin, error) {
	// ソート属性
	// "sp"
	// "勢い順↓"
	boardWinModel := &BoardWinModel{SortAttrStr: "sp"}

	// 板ウィンドウ生成
	boardWin := new(BoardWin)

	// 板名の格納
	boardWin.boardName = boardName
	// モデルの生成
	boardWin.threadListModel = NewThreadListModel(boardName, boardWinModel.SortAttrStr)

	boardWin.boardTitle = boardWin.threadListModel.boardTitle

	// アイコン
	icon := GetApplicationIcon()

	// メインウィンドウのウィンドウ生成
	err := MainWindow{
		AssignTo: &boardWin.MainWindow,
		Title:    boardWin.boardTitle + " - " + AppName + " " + Version,
		Icon:     icon,
		//MinSize:	Size{600, 400},
		MinSize: Size{600, 800},
		Layout:  VBox{},
		DataBinder: DataBinder{
			DataSource: boardWinModel,
			AutoSubmit: true,
			OnSubmitted: func() {
				//fmt.Println("DataBinder::OnSubmitted start\r\n")
				//fmt.Println(boardWinModel)

				// モデルの生成
				boardWin.threadListModel = NewThreadListModel(boardWin.boardName, boardWinModel.SortAttrStr)
				// リストボックスのカレントインデックスを非選択に設定する
				// これをやらないとモデル再設定中にインデックスの変更イベントが発生して
				//   panic: runtime error: index out of range
				// で落ちるみたい
				boardWin.listBoxThread.SetCurrentIndex(-1)
				// モデルを再設定する
				boardWin.listBoxThread.SetModel(boardWin.threadListModel)
				//fmt.Println("DataBinder::OnSubmitted end\r\n")
			},
		},
		Children: []Widget{
			///////////////////
			GroupBox{
				Layout: HBox{},
				Children: []Widget{
					///////////////////
					// RadioButtonGroup is needed for data binding only.
					RadioButtonGroup{
						DataMember: "SortAttrStr",
						Buttons: []RadioButton{
							RadioButton{
								Name:  "spRB",
								Text:  "勢い順↓",
								Value: "sp",
							},
							RadioButton{
								Name:  "sp2RB",
								Text:  "勢い順↑",
								Value: "sp2",
							},
							RadioButton{
								Name:  "siRB",
								Text:  "時間順↓",
								Value: "si",
							},
							RadioButton{
								Name:  "si2RB",
								Text:  "時間順↑",
								Value: "si2",
							},
							RadioButton{
								Name:  "reRB",
								Text:  "レス数順↓",
								Value: "re",
							},
							RadioButton{
								Name:  "re2RB",
								Text:  "レス数順↑",
								Value: "re2",
							},
							RadioButton{
								Name:  "defRB",
								Text:  "番号順↑",
								Value: "",
							},
							RadioButton{
								Name:  "noRB",
								Text:  "番号順↓",
								Value: "no",
							},
						},
					},
					///////////////////
				},
			},
			///////////////////

			ListBox{
				AssignTo: &boardWin.listBoxThread,
				Model:    boardWin.threadListModel,
				OnCurrentIndexChanged: boardWin.listBoxThreadCurrentIndexChanged,
				OnItemActivated:       boardWin.listBoxThreadItemActivated,
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
	boardWin.listBoxThread.SetFont(font)

	// 子ウィンドウ化
	// Note:win.SetParentだと、親WindowにクリッピングされたWindowになってしまう
	//win.SetParent(boardWin.Handle(), parentWin.Handle())
	// FormのSetOwnerを使うといいみたい
	boardWin.SetOwner(parentWin)

	return boardWin, err
}

/**
 * スレッド一覧リストボックス選択インデックスが変わった
 * @param なし
 * @return なし
 */
func (boardWin *BoardWin) listBoxThreadCurrentIndexChanged() {
	//i := boardWin.listBoxThread.CurrentIndex()
	//item := &boardWin.threadListModel.items[i]

	//name := item.name
	//value := item.value
	//fmt.Println("CurrentIndex: ", i)
	//fmt.Println("name: ", name)
	//fmt.Println("value: ", value)
}

/**
 * スレッド一覧リストボックスアイテムがダブルクリックされた
 * @param なし
 * @return なし
 */
func (boardWin *BoardWin) listBoxThreadItemActivated() {
	i := boardWin.listBoxThread.CurrentIndex()
	item := &boardWin.threadListModel.items[i]

	//name := item.name
	value := item.value
	//fmt.Println(name)
	//fmt.Println(value)
	//fmt.Println("Num=", value.Num)
	//fmt.Println("Res=", value.Res)
	//fmt.Println("Thread=", value.Thread)
	//fmt.Println("Sin=", value.Sin)
	//fmt.Println("Since=", value.Since)
	//fmt.Println("Spd=", value.Spd)

	//model := unkarstub.GetThreadModel(boardWin.boardName, value.Sin)
	//fmt.Println(model)

	// スレッドウィンドウの生成
	threadWin, err := NewThreadWin(boardWin, boardWin.boardName, value.Sin)
	if err != nil {
		log.Fatal(err)
	}
	// 表示
	threadWin.Run()
}

////////////////////////////////////////////////////////////
// ThreadListItem
////////////////////////////////////////////////////////////
/**
 * リストボックスアイテム
 */
type ThreadListItem struct {
	name  string
	value *unmodel.ThreadItem
}

////////////////////////////////////////////////////////////
// ThreadListModel
////////////////////////////////////////////////////////////
/**
 * リストボックスモデル
 */
type ThreadListModel struct {
	// 派生元：リストモデルベース
	walk.ListModelBase
	// アイテム一覧
	items []ThreadListItem
	// 板タイトル
	boardTitle string
}

/**
 * コンストラクタ
 * @param boardName
 * @param sortAttrStr
 * @return スレッドリストモデル
 */
func NewThreadListModel(boardName string, sortAttrStr string) *ThreadListModel {
	// 板のモデルを取得する
	unutilModel := unkarstub.GetBoardModel(boardName)

	// DEBUG
	fmt.Printf("url=%s\r\n", unutilModel.GetUrl())
	fmt.Printf("title=%s\r\n", unutilModel.GetTitle())
	fmt.Printf("className=%s\r\n", unutilModel.GetClassName())
	fmt.Printf("server=%s\r\n", unutilModel.GetServer())

	attr := sortAttrStr
	nowData := unutilModel.GetData()
	sort, dir := unmodel.BoardSort(unutilModel, attr)
	fmt.Printf("sort=\r\n")
	fmt.Print(sort)
	fmt.Print("\r\n")
	fmt.Printf("dir=\r\n")
	fmt.Print(dir)
	fmt.Printf("\r\n")

	// スレッド一覧を取得する
	list, ok := nowData.(*unmodel.ThreadItems)
	if !ok {
		err := errors.New("failed to get nowData")
		log.Fatal(err)
	}
	//////////////////////////////////////////

	// リストボックスのモデルを生成
	model := &ThreadListModel{
		items:      make([]ThreadListItem, len(*list)),
		boardTitle: unutilModel.GetTitle(),
	}
	for i, thread := range *list {
		// DEBUG
		//fmt.Println("thread=", thread)

		name := thread.Thread
		value := thread
		model.items[i] = ThreadListItem{name, value}
	}

	return model
}

/**
 * アイテム数を取得する
 * @return アイテム数
 */
func (model *ThreadListModel) ItemCount() int {
	return len(model.items)
}

/**
 * アイテムの値を取得する
 * @return アイテムの値
 */
func (model *ThreadListModel) Value(index int) interface{} {
	return model.items[index].name
}
