package main

import (
	"./undity"
	"fmt"
	"os"
)

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

const (
	AppName   = "unkarApp"
	Version   = "1.0.0.7"
	TitleBase = AppName + " " + Version
)

/**
 *  アプリケーションのアイコン
 */
func GetApplicationIcon() *walk.Icon {
	// アイコン
	//icon, iconErr := walk.Resources.Icon("unkarApp.ico")
	icon, iconErr := walk.Resources.Icon("3")
	if iconErr != nil {
		panic(iconErr)
	}
	return icon
}

////////////////////////////////////////////////////////////
// MainWin
////////////////////////////////////////////////////////////
/**
 * ページ
 */
type Page interface {
	walk.Container
	Parent() walk.Container
	SetParent(parent walk.Container) error
	Title() string
}

/**
 * メインウィンドウ
 */
type MainWin struct {
	*walk.MainWindow
	navToolBar    *walk.ToolBar
	pageActions   []*walk.Action
	pageComposite *walk.Composite
	action2Page   map[*walk.Action]Page
	currPage      Page
	currAction    *walk.Action
	topPage       *TopPage
	boardPage     *BoardPage
	threadPage    *ThreadPage
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
	mainWin := &MainWin{
		action2Page: make(map[*walk.Action]Page),
	}

	// アイコン
	icon := GetApplicationIcon()

	// メインウィンドウのウィンドウ生成
	err := MainWindow{
		AssignTo: &mainWin.MainWindow,
		Title:    TitleBase,
		Icon:     icon,
		MinSize:  Size{950, 600},
		Layout:   HBox{MarginsZero: true, SpacingZero: true},
		Font:     Font{Family: "MS Shell Dlg 2", PointSize: 12},
		Children: []Widget{
			ScrollView{
				HorizontalFixed: true,
				Layout:          VBox{MarginsZero: true},
				Children: []Widget{
					Composite{
						Layout: VBox{MarginsZero: true},
						Children: []Widget{
							ToolBar{
								AssignTo:    &mainWin.navToolBar,
								Orientation: Vertical,
								ButtonStyle: ToolBarButtonImageAboveText,
								MaxTextRows: 2,
							},
						},
					},
				},
			},
			Composite{
				AssignTo: &mainWin.pageComposite,
				Name:     "pageComposite",
				Layout:   HBox{MarginsZero: true, SpacingZero: true},
			},
		},
	}.Create()

	mainWin.topPage, err = newTopPage(mainWin.pageComposite, mainWin)
	if err != nil {
		return nil, err
	}
	mainWin.boardPage, err = newBoardPage(mainWin.pageComposite, mainWin)
	if err != nil {
		return nil, err
	}
	mainWin.threadPage, err = newThreadPage(mainWin.pageComposite, mainWin)
	if err != nil {
		return nil, err
	}
	action, err := mainWin.newPageAction("板一覧", "./undity/public_html/img/whiteboard.png")
	if err != nil {
		return nil, err
	}
	mainWin.action2Page[action] = mainWin.topPage
	mainWin.pageActions = append(mainWin.pageActions, action)

	action, err = mainWin.newPageAction("板", "./undity/public_html/img/folder.png")
	if err != nil {
		return nil, err
	}
	mainWin.action2Page[action] = mainWin.boardPage
	mainWin.pageActions = append(mainWin.pageActions, action)

	action, err = mainWin.newPageAction("スレッド", "./undity/public_html/img/memo.png")
	if err != nil {
		return nil, err
	}
	mainWin.action2Page[action] = mainWin.threadPage
	mainWin.pageActions = append(mainWin.pageActions, action)

	mainWin.updateNavigationToolBar()

	if len(mainWin.pageActions) > 0 {
		if err := mainWin.setCurrentAction(mainWin.pageActions[0]); err != nil {
			return nil, err
		}
	}

	return mainWin, err
}

func (mainWin *MainWin) UpdateTitle(page Page) {
	s := page.Title()
	if s != "" {
		s += " - "
	}
	s += TitleBase
	mainWin.MainWindow.SetTitle(s)
}

func (mainWin *MainWin) newPageAction(title, image string) (*walk.Action, error) {
	img, err := walk.Resources.Bitmap(image)
	if err != nil {
		return nil, err
	}

	action := walk.NewAction()
	action.SetCheckable(true)
	action.SetExclusive(true)
	action.SetImage(img)
	action.SetText(title)

	action.Triggered().Attach(func() {
		mainWin.setCurrentAction(action)
	})

	return action, nil
}

func (mainWin *MainWin) setCurrentAction(action *walk.Action) error {
	defer func() {
		if !mainWin.pageComposite.IsDisposed() {
			mainWin.pageComposite.RestoreState()
			mainWin.pageComposite.Layout().Update(false)
		}
	}()

	mainWin.SetFocus()

	if prevPage := mainWin.currPage; prevPage != nil {
		mainWin.pageComposite.SaveState()
		prevPage.SetVisible(false)
		prevPage.SetParent(nil)
	}

	page := mainWin.action2Page[action]
	page.SetParent(mainWin.pageComposite)
	page.SetVisible(true)

	action.SetChecked(true)

	mainWin.currPage = page
	mainWin.currAction = action
	mainWin.UpdateTitle(page)

	return nil
}

func (mainWin *MainWin) updateNavigationToolBar() error {
	mainWin.navToolBar.SetSuspended(true)
	defer mainWin.navToolBar.SetSuspended(false)

	actions := mainWin.navToolBar.Actions()

	if err := actions.Clear(); err != nil {
		return err
	}

	for _, action := range mainWin.pageActions {
		if err := actions.Add(action); err != nil {
			return err
		}
	}

	if mainWin.currAction != nil {
		if !actions.Contains(mainWin.currAction) {
			for _, action := range mainWin.pageActions {
				if action != mainWin.currAction {
					if err := mainWin.setCurrentAction(action); err != nil {
						return err
					}

					break
				}
			}
		}
	}

	return nil
}

func (mainWin *MainWin) topPageAction() *walk.Action {
	var tgtAction *walk.Action = nil
	tgtPage := mainWin.topPage
	for workAction, workPage := range mainWin.action2Page {
		if workPage == tgtPage {
			tgtAction = workAction
			break
		}
	}
	return tgtAction
}

func (mainWin *MainWin) boardPageAction() *walk.Action {
	var tgtAction *walk.Action = nil
	tgtPage := mainWin.boardPage
	for workAction, workPage := range mainWin.action2Page {
		if workPage == tgtPage {
			tgtAction = workAction
			break
		}
	}
	return tgtAction
}

func (mainWin *MainWin) threadPageAction() *walk.Action {
	var tgtAction *walk.Action = nil
	tgtPage := mainWin.threadPage
	for workAction, workPage := range mainWin.action2Page {
		if workPage == tgtPage {
			tgtAction = workAction
			break
		}
	}
	return tgtAction
}

func (mainWin *MainWin) NavigateToTopPage() {
	action := mainWin.topPageAction()
	mainWin.topPage.UpdateContents()
	mainWin.changePage(action)
}

func (mainWin *MainWin) NavigateToBoardPage(boardName string, boardKey string, sortAttrStr string) {
	action := mainWin.boardPageAction()
	mainWin.boardPage.UpdateContents(boardName, boardKey, sortAttrStr)
	mainWin.changePage(action)
}

func (mainWin *MainWin) NavigateToThreadPage(boardName string, boardKey string, threadNo int64) {
	action := mainWin.threadPageAction()
	mainWin.threadPage.UpdateContents(boardName, boardKey, threadNo)
	mainWin.changePage(action)
}

func (mainWin *MainWin) changePage(action *walk.Action) {
	mainWin.clearToolBarChecked()
	mainWin.setCurrentAction(action)
}

func (mainWin *MainWin) clearToolBarChecked() {
	for _, workAction := range mainWin.pageActions {
		workAction.SetChecked(false)
	}
	mainWin.currAction = nil
}
