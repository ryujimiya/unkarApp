package main

import (
	"./undity"
	"./undity/golib/model"
	"errors"
	"fmt"
)

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

//import (
//	"github.com/lxn/win"
//)

////////////////////////////////////////////////////////////
// BoardPage
////////////////////////////////////////////////////////////
/**
 * 板ページモデル
 */
type BoardPageModel struct {
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
 * 板ページ
 */
type BoardPage struct {
	*walk.Composite
	mainWin         *MainWin             // メインウィンドウ
	db              *walk.DataBinder     // データバインダー
	rbSort          [8]*walk.RadioButton // ソート属性ラジオボタン一覧
	listBoxThread   *walk.ListBox        // スレッド一覧リストボックス
	boardKey        string               // 板キー
	boardName       string               // 板タイトル
	boardPageModel  *BoardPageModel      // 板モデル
	threadListModel *ThreadListModel     // スレッド一覧モデル
	title           string               // タイトル
}

func newBoardPage(parent walk.Container, mainWin *MainWin) (*BoardPage, error) {
	// 板ページ生成
	boardPage := new(BoardPage)

	boardPage.mainWin = mainWin

	// ソート属性
	// "sp"
	// "勢い順↓"
	boardPage.boardPageModel = &BoardPageModel{SortAttrStr: "sp"}

	if err := (Composite{
		AssignTo: &boardPage.Composite,
		Name:     "板",
		Layout:   VBox{},
		DataBinder: DataBinder{
			AssignTo:   &boardPage.db,
			DataSource: boardPage.boardPageModel,
			AutoSubmit: true,
			OnSubmitted: func() {
				//fmt.Printf("DataBinder::OnSubmitted start\r\n")
				//fmt.Printf("%+v\r\n", boardPage.boardPageModel)

				boardPage.UpdateContents(boardPage.boardName, boardPage.boardKey, boardPage.boardPageModel.SortAttrStr)

				//fmt.Println("DataBinder::OnSubmitted end\r\n")
			},
		},
		Children: []Widget{
			GroupBox{
				Layout: HBox{},
				Children: []Widget{
					// RadioButtonGroup is needed for data binding only.
					RadioButtonGroup{
						DataMember: "SortAttrStr",
						Buttons: []RadioButton{
							RadioButton{
								AssignTo: &boardPage.rbSort[0],
								Name:     "spRB",
								Text:     "勢い順↓",
								Value:    "sp",
							},
							RadioButton{
								AssignTo: &boardPage.rbSort[1],
								Name:     "sp2RB",
								Text:     "勢い順↑",
								Value:    "sp2",
							},
							RadioButton{
								AssignTo: &boardPage.rbSort[2],
								Name:     "siRB",
								Text:     "時間順↓",
								Value:    "si",
							},
							RadioButton{
								AssignTo: &boardPage.rbSort[3],
								Name:     "si2RB",
								Text:     "時間順↑",
								Value:    "si2",
							},
							RadioButton{
								AssignTo: &boardPage.rbSort[4],
								Name:     "reRB",
								Text:     "レス数順↓",
								Value:    "re",
							},
							RadioButton{
								AssignTo: &boardPage.rbSort[5],
								Name:     "re2RB",
								Text:     "レス数順↑",
								Value:    "re2",
							},
							RadioButton{
								AssignTo: &boardPage.rbSort[6],
								Name:     "defRB",
								Text:     "番号順↑",
								Value:    "",
							},
							RadioButton{
								AssignTo: &boardPage.rbSort[7],
								Name:     "noRB",
								Text:     "番号順↓",
								Value:    "no",
							},
						},
					},
				},
			},

			ListBox{
				AssignTo:              &boardPage.listBoxThread,
				OnCurrentIndexChanged: boardPage.listBoxThreadCurrentIndexChanged,
				OnItemActivated:       boardPage.listBoxThreadItemActivated,
			},
		},
		Visible: false,
	}).Create(NewBuilder(parent)); err != nil {
		return nil, err
	}

	if err := walk.InitWrapperWindow(boardPage); err != nil {
		return nil, err
	}

	return boardPage, nil
}

func (boardPage *BoardPage) Title() string {
	return boardPage.title
}

func (boardPage *BoardPage) UpdateContents(boardName string, boardKey string, sortAttrStr string) {
	boardPage.boardName = boardName
	boardPage.boardKey = boardKey
	prevSortAttr := boardPage.boardPageModel.SortAttrStr
	boardPage.boardPageModel.SortAttrStr = sortAttrStr

	if prevSortAttr != boardPage.boardPageModel.SortAttrStr {
		// ラジオボタンのイベントではここにはこない
		// 外部から呼ばれた場合、ここにくる
		// ラジオボタンのチェック状態の更新
		for _, rb := range boardPage.rbSort {
			//rb.SetChecked(rb.Value() == boardPage.boardPageModel.SortAttrStr)
			//だとrbgSort.CheckButton()が変わらないのでCheckedValueプロパティをセットする
			prop := rb.AsWindowBase().Property("CheckedValue")
			prop.Set(boardPage.boardPageModel.SortAttrStr)
		}
		return
	}

	if len(boardName) == 0 || len(boardKey) == 0 || len(sortAttrStr) == 0 {
		//fmt.Printf("BoardPage.UpdateContents parameter is nothing\r\n")
		return
	}

	// モデルの生成
	boardPage.threadListModel = NewThreadListModel(boardPage.boardKey, boardPage.boardPageModel.SortAttrStr)
	// モデルを再設定する
	boardPage.listBoxThread.SetCurrentIndex(-1)
	boardPage.listBoxThread.SetModel(boardPage.threadListModel)

	boardPage.title = boardPage.boardName
	boardPage.mainWin.UpdateTitle(boardPage)
}

/**
 * スレッド一覧リストボックス選択インデックスが変わった
 * @param なし
 * @return なし
 */
func (boardPage *BoardPage) listBoxThreadCurrentIndexChanged() {
	//i := boardPage.listBoxThread.CurrentIndex()
	//item := &boardPage.threadListModel.items[i]

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
func (boardPage *BoardPage) listBoxThreadItemActivated() {
	i := boardPage.listBoxThread.CurrentIndex()
	item := &boardPage.threadListModel.items[i]

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

	// スレッドページの表示
	boardPage.mainWin.NavigateToThreadPage(boardPage.boardName, boardPage.boardKey, value.Sin)
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
	walk.ListModelBase                  // 派生元：リストモデルベース
	items              []ThreadListItem // アイテム一覧
	boardName          string           // 板名
}

/**
 * コンストラクタ
 * @param boardKey
 * @param sortAttrStr
 * @return スレッドリストモデル
 */
func NewThreadListModel(boardKey string, sortAttrStr string) *ThreadListModel {
	// 板のモデルを取得する
	unutilModel := unkarstub.GetBoardModel(boardKey)

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
		panic(err)
	}
	//////////////////////////////////////////

	// リストボックスのモデルを生成
	model := &ThreadListModel{
		items:     make([]ThreadListItem, len(*list)),
		boardName: unutilModel.GetTitle(),
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
