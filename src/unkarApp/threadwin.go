package main

import (
	"./undity"
	//"./undity/golib/util"
	//"./undity/golib/model"
	"fmt"
	"strings"
	//"log"
	"os"
)

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

////////////////////////////////////////////////////////////
// ThreadWin
////////////////////////////////////////////////////////////
/**
 * スレッドウィンドウ
 */
type ThreadWin struct {
	// 派生元：Walkメインウィンドウ
	*walk.MainWindow
	// 表示用WebView
	webView *walk.WebView
	// 板名
	boardName string
	// スレッド番号
	threadNo int64
	// スレッドタイトル
	threadTitle string
	// ファイルパス
	threadPageFilepath string
	// 表示用htmlのURL
	threadPageUrl string
	// 現在のURL
	url string
	// ファイルを作成する？
	fileCreate bool
}

/**
 * コンストラクタ
 * @param parentWin 親ウィンドウ
 * @param boardName URLに含まれる板名(povertyなど)
 * @param threadNo スレッド番号
 * @return (1)メインウィンドウ
 *         (2)エラー
 */
func NewThreadWin(parentWin walk.Form, boardName string, threadNo int64) (*ThreadWin, error) {

	// 板ウィンドウ生成
	threadWin := new(ThreadWin)

	// 板名の格納
	threadWin.boardName = boardName
	// スレッド番号の格納
	threadWin.threadNo = threadNo

	// アイコン
	icon := GetApplicationIcon()

	threadWin.url = ""
	threadWin.fileCreate = true

	// tmpディレクトリ
	tmpDir := unkarstub.GetTmpHtmlDir()
	threadPageFilename := fmt.Sprintf("%s_%d.html", boardName, threadNo)
	// ファイルパス
	threadWin.threadPageFilepath = tmpDir + "\\" + threadPageFilename
	fmt.Printf("threadWin.threadPageFilepath=" + threadWin.threadPageFilepath + "\r\n")
	// 表示用htmlのURL
	tmpDir = strings.Replace(tmpDir, "\\", "/", -1)
	threadWin.threadPageUrl = "file:///" + tmpDir + "/" + threadPageFilename
	fmt.Printf("threadWin.threadPageUrl=" + threadWin.threadPageUrl + "\r\n")

	// メインウィンドウのウィンドウ生成
	err := MainWindow{
		AssignTo: &threadWin.MainWindow,
		Title:    "",
		Icon:     icon,
		MinSize:  Size{850, 600},
		Layout:   VBox{},
		Children: []Widget{
			WebView{
				AssignTo:                 &threadWin.webView,
				Name:                     "スレッド",
				URL:                      threadWin.threadPageUrl,
				ShortcutsEnabled:         true,
				NativeContextMenuEnabled: true,
				OnNavigating:             threadWin.webView_OnNavigating,
				OnNavigated:              threadWin.webView_OnNavigated,
				OnDownloading:            threadWin.webView_OnDownloading,
				OnDocumentCompleted:      threadWin.webView_OnDocumentCompleted,
				OnNavigatedError:         threadWin.webView_OnNavigatedError,
				OnNewWindow:              threadWin.webView_OnNewWindow,
			},
		},
	}.Create()

	// 子ウィンドウ化
	// Note:win.SetParentだと、親WindowにクリッピングされたWindowになってしまう
	//win.SetParent(threadWin.Handle(), parentWin.Handle())
	// FormのSetOwnerを使うといいみたい
	threadWin.SetOwner(parentWin)

	return threadWin, err
}

func (threadWin *ThreadWin) webView_OnNavigating(arg *walk.WebViewNavigatingArg) {
	fmt.Printf("webView_OnNavigating\r\n")
	fmt.Printf("Url = %+v\r\n", arg.Url())
	fmt.Printf("Flags = %+v\r\n", arg.Flags())
	fmt.Printf("Headers = %+v\r\n", arg.Headers())
	fmt.Printf("TargetFrameName = %+v\r\n", arg.TargetFrameName())
	fmt.Printf("Cancel = %+v\r\n", arg.Cancel())
	// if you want to cancel
	//arg.SetCancel(true)

	// URLを格納する
	// このURLはfile:///C:/でなくC:\ で渡ってくるので注意
	threadWin.url = arg.Url()
}

func (threadWin *ThreadWin) webView_OnNavigated(arg *walk.WebViewNavigatedEventArg) {
	fmt.Printf("webView_OnNavigated\r\n")
	fmt.Printf("Url = %+v\r\n", arg.Url())
}

func (threadWin *ThreadWin) webView_OnDownloading() {
	fmt.Printf("webView_OnDownloading\r\n")

	if threadWin.fileCreate {
		//fmt.Printf("%+v\r\n", threadWin.url)
		//fmt.Printf("%+v\r\n", threadWin.threadPageFilepath)
		// スレッドページだったらファイルを作成する
		// NOTE: threadWin.urlが c:\と小文字で始まる現象
		//       ⇒コマンドプロンプト上で実行したときの問題
		//         cd c:\とすると発生する
		if strings.ToLower(threadWin.url) == strings.ToLower(threadWin.threadPageFilepath) {
			// HTMLファイルを作成
			threadWin.createHtmlFile()

			// このイベントが発生するときはすでにダウンロードが始まっているらしい
			// 作成したHTMLファイルが表示に反映されない
			// 仕方ないので再度更新する
			threadWin.webView.Refresh()
		}

		threadWin.fileCreate = false
	} else {
		threadWin.fileCreate = true
	}
}

func (threadWin *ThreadWin) webView_OnDownloaded() {
	fmt.Printf("webView_OnDownloaded\r\n")
}

func (threadWin *ThreadWin) webView_OnDocumentCompleted(arg *walk.WebViewDocumentCompletedEventArg) {
	fmt.Printf("webView_OnDocumentCompleted\r\n")
	fmt.Printf("Url = %+v\r\n", arg.Url())
}

func (threadWin *ThreadWin) webView_OnNavigatedError(arg *walk.WebViewNavigatedErrorEventArg) {
	fmt.Printf("webView_OnNavigatedError\r\n")
	fmt.Printf("Url = %+v\r\n", arg.Url())
	fmt.Printf("TargetFrameName = %+v\r\n", arg.TargetFrameName())
	fmt.Printf("StatusCode = %+v\r\n", arg.StatusCode())
	fmt.Printf("Cancel = %+v\r\n", arg.Cancel())
	// if you want to cancel
	//arg.SetCancel(true)
}

func (threadWin *ThreadWin) webView_OnNewWindow(arg *walk.WebViewNewWindowEventArg) {
	fmt.Printf("webView_OnNewWindow\r\n")
	fmt.Printf("Cancel = %+v\r\n", arg.Cancel())
	fmt.Printf("Flags = %+v\r\n", arg.Flags())
	fmt.Printf("UrlContext = %+v\r\n", arg.UrlContext())
	fmt.Printf("Url = %+v\r\n", arg.Url())
	// if you want to cancel
	//arg.SetCancel(true)
}

func (threadWin *ThreadWin) createHtmlFile() {
	boardName := threadWin.boardName
	threadNo := threadWin.threadNo
	// ページ属性
	attr := ""
	//attr = "Anchor:@0!3" // まとめ（仮）
	//attr = "Anchor:Default" // アンカー
	//attr = "Link:All" // URL
	//attr = "Tree:Link:All" // URL[+]
	//attr = "Link:Image" // 画像
	//attr = "Tree:Link:Image" // 画像[+]
	//attr = "Link:Movie" // 動画
	//attr = "Tree:Link:Movie" // 動画[+]

	// 本文作成
	htmlText := threadWin.createHtmlText(boardName, threadNo, attr)
	// 表示用htmlファイルに保存
	threadWin.saveToTmpHtml(threadWin.threadPageFilepath, htmlText)

	threadWin.SetTitle(threadWin.threadTitle + " - " + AppName + " " + Version)
}

func (threadWin *ThreadWin) createHtmlText(boardName string, threadNo int64, attr string) string {
	// スレッドのモデルを取得する
	unutilModel := unkarstub.GetThreadModel(boardName, threadNo, attr)

	// DEBUG
	fmt.Printf("url=%s\r\n", unutilModel.GetUrl())
	fmt.Printf("title=%s\r\n", unutilModel.GetTitle())
	fmt.Printf("className=%s\r\n", unutilModel.GetClassName())
	fmt.Printf("server=%s\r\n", unutilModel.GetServer())

	threadWin.threadTitle = unutilModel.GetTitle()

	htmlText := unkarstub.GetBoardViewOutput(boardName, threadNo, attr, unutilModel)

	return htmlText
}

func (threadWin *ThreadWin) saveToTmpHtml(filename, htmlStr string) {
	out, err := os.Create(filename)
	if err != nil {
		return
	}
	fmt.Fprint(out, htmlStr)
	out.Close()
}
