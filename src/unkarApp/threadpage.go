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
type ThreadPage struct {
	*walk.Composite
	mainWin            *MainWin
	webView            *walk.WebView // 表示用WebView
	boardName          string        // 板名
	boardKey           string        // 板キー
	threadNo           int64         // スレッド番号
	threadTitle        string        // スレッドタイトル
	threadPageFilepath string        // ファイルパス
	threadPageUrl      string        // 表示用htmlのURL
	url                string        // 現在のURL
	fileCreate         bool          // ファイルを作成する？
}

func newThreadPage(parent walk.Container, mainWin *MainWin) (*ThreadPage, error) {
	// 板ページ生成
	threadPage := new(ThreadPage)

	threadPage.mainWin = mainWin

	threadPage.url = ""
	threadPage.fileCreate = true

	if err := (Composite{
		AssignTo: &threadPage.Composite,
		Name:     "スレッド",
		Layout:   VBox{},
		Children: []Widget{
			WebView{
				AssignTo:                 &threadPage.webView,
				Name:                     "スレッド",
				URL:                      "",
				ShortcutsEnabled:         true,
				NativeContextMenuEnabled: true,
				OnNavigating:             threadPage.webView_OnNavigating,
				OnNavigated:              threadPage.webView_OnNavigated,
				OnDownloading:            threadPage.webView_OnDownloading,
				OnDownloaded:             threadPage.webView_OnDownloaded,
				OnDocumentCompleted:      threadPage.webView_OnDocumentCompleted,
				OnNavigatedError:         threadPage.webView_OnNavigatedError,
				OnNewWindow:              threadPage.webView_OnNewWindow,
			},
		},
		Visible: false,
	}).Create(NewBuilder(parent)); err != nil {
		return nil, err
	}

	if err := walk.InitWrapperWindow(threadPage); err != nil {
		return nil, err
	}

	return threadPage, nil
}

func (threadPage *ThreadPage) UpdateContents(boardName string, boardKey string, threadNo int64) {
	// 板名の格納
	threadPage.boardName = boardName
	// 板キーの格納
	threadPage.boardKey = boardKey
	// スレッド番号の格納
	threadPage.threadNo = threadNo

	// tmpディレクトリ
	tmpDir := unkarstub.GetTmpHtmlDir()
	threadPageFilename := fmt.Sprintf("%s_%d.html", boardKey, threadNo)
	// ファイルパス
	threadPage.threadPageFilepath = tmpDir + "\\" + threadPageFilename
	fmt.Printf("threadPage.threadPageFilepath=" + threadPage.threadPageFilepath + "\r\n")
	// 表示用htmlのURL
	tmpDir = strings.Replace(tmpDir, "\\", "/", -1)
	threadPage.threadPageUrl = "file:///" + tmpDir + "/" + threadPageFilename
	fmt.Printf("threadPage.threadPageUrl=" + threadPage.threadPageUrl + "\r\n")

	// HTMLファイルを作成
	threadPage.createHtmlFile()
	threadPage.fileCreate = false

	threadPage.webView.SetURL(threadPage.threadPageUrl)
}

func (threadPage *ThreadPage) webView_OnNavigating(eventData *walk.WebViewNavigatingEventData) {
	fmt.Printf("webView_OnNavigating\r\n")
	fmt.Printf("Url = %+v\r\n", eventData.Url())
	fmt.Printf("Flags = %+v\r\n", eventData.Flags())
	fmt.Printf("Headers = %+v\r\n", eventData.Headers())
	fmt.Printf("TargetFrameName = %+v\r\n", eventData.TargetFrameName())
	fmt.Printf("Canceled = %+v\r\n", eventData.Canceled())
	// if you want to cancel
	//eventData.SetCanceled(true)

	// URLを格納する
	// このURLはfile:///C:/でなくC:\ で渡ってくるので注意
	threadPage.url = eventData.Url()
}

func (threadPage *ThreadPage) webView_OnNavigated(url string) {
	fmt.Printf("webView_OnNavigated\r\n")
	fmt.Printf("url = %+v\r\n", url)
}

func (threadPage *ThreadPage) webView_OnDownloading() {
	fmt.Printf("webView_OnDownloading\r\n")

	if threadPage.fileCreate {
		//fmt.Printf("%+v\r\n", threadPage.url)
		//fmt.Printf("%+v\r\n", threadPage.threadPageFilepath)
		// スレッドページだったらファイルを作成する
		// NOTE: threadPage.urlが c:\と小文字で始まる現象
		//       ⇒コマンドプロンプト上で実行したときの問題
		//         cd c:\とすると発生する
		if strings.ToLower(threadPage.url) == strings.ToLower(threadPage.threadPageFilepath) {
			// HTMLファイルを作成
			threadPage.createHtmlFile()

			// このイベントが発生するときはすでにダウンロードが始まっているらしい
			// 作成したHTMLファイルが表示に反映されない
			// 仕方ないので再度更新する
			threadPage.webView.Refresh()
		}

		threadPage.fileCreate = false
	} else {
		threadPage.fileCreate = true
	}
}

func (threadPage *ThreadPage) webView_OnDownloaded() {
	fmt.Printf("webView_OnDownloaded\r\n")
}

func (threadPage *ThreadPage) webView_OnDocumentCompleted(url string) {
	fmt.Printf("webView_OnDocumentCompleted\r\n")
	fmt.Printf("url = %+v\r\n", url)
}

func (threadPage *ThreadPage) webView_OnNavigatedError(eventData *walk.WebViewNavigatedErrorEventData) {
	fmt.Printf("webView_OnNavigatedError\r\n")
	fmt.Printf("Url = %+v\r\n", eventData.Url())
	fmt.Printf("TargetFrameName = %+v\r\n", eventData.TargetFrameName())
	fmt.Printf("StatusCode = %+v\r\n", eventData.StatusCode())
	fmt.Printf("Canceled = %+v\r\n", eventData.Canceled())
	// if you want to cancel
	//eventData.SetCanceled(true)
}

func (threadPage *ThreadPage) webView_OnNewWindow(eventData *walk.WebViewNewWindowEventData) {
	fmt.Printf("webView_OnNewWindow\r\n")
	fmt.Printf("Canceled = %+v\r\n", eventData.Canceled())
	fmt.Printf("Flags = %+v\r\n", eventData.Flags())
	fmt.Printf("UrlContext = %+v\r\n", eventData.UrlContext())
	fmt.Printf("Url = %+v\r\n", eventData.Url())
	// if you want to cancel
	//eventData.SetCancel(true)
}

func (threadPage *ThreadPage) createHtmlFile() {
	boardKey := threadPage.boardKey
	threadNo := threadPage.threadNo
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
	htmlText := threadPage.createHtmlText(boardKey, threadNo, attr)
	// 表示用htmlファイルに保存
	threadPage.saveToTmpHtml(threadPage.threadPageFilepath, htmlText)

	//threadPage.mainWin.SetTitle(threadPage.threadTitle + " - " + AppName + " " + Version)
}

func (threadPage *ThreadPage) createHtmlText(boardKey string, threadNo int64, attr string) string {
	// スレッドのモデルを取得する
	unutilModel := unkarstub.GetThreadModel(boardKey, threadNo, attr)

	// DEBUG
	fmt.Printf("url=%s\r\n", unutilModel.GetUrl())
	fmt.Printf("title=%s\r\n", unutilModel.GetTitle())
	fmt.Printf("className=%s\r\n", unutilModel.GetClassName())
	fmt.Printf("server=%s\r\n", unutilModel.GetServer())

	threadPage.threadTitle = unutilModel.GetTitle()

	htmlText := unkarstub.GetBoardViewOutput(boardKey, threadNo, attr, unutilModel)

	return htmlText
}

func (threadPage *ThreadPage) saveToTmpHtml(filename, htmlStr string) {
	out, err := os.Create(filename)
	if err != nil {
		return
	}
	fmt.Fprint(out, htmlStr)
	out.Close()
}
