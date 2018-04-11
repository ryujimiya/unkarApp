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
	"github.com/lxn/win"
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
				AssignTo:           &threadWin.webView,
				Name:               "スレッド",
				URL:                threadWin.threadPageUrl,
				ShortcutsEnabled:   true,
				ContextMenuEnabled: true,
				BeforeNavigate2:    threadWin.webView_BeforeNavigate2,
				NavigateComplete2:  threadWin.webView_NavigateComplete2,
				DownloadBegin:      threadWin.webView_DownloadBegin,
				DocumentComplete:   threadWin.webView_DocumentComplete,
				NavigateError:      threadWin.webView_NavigateError,
				NewWindow3:         threadWin.webView_NewWindow3,
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

func (threadWin *ThreadWin) webView_BeforeNavigate2(
	pDisp *win.IDispatch,
	url *win.VARIANT,
	flags *win.VARIANT,
	targetFrameName *win.VARIANT,
	postData *win.VARIANT,
	headers *win.VARIANT,
	cancel *win.VARIANT_BOOL) {

	fmt.Printf("webView_BeforeNavigate2\r\n")
	fmt.Printf("pDisp = %+v\r\n", pDisp)
	fmt.Printf("url = %+v\r\n", url)
	if url != nil && url.BstrVal() != nil {
		fmt.Printf("  url = %+v\r\n", win.BSTRToString(url.BstrVal()))
	}
	fmt.Printf("flags = %+v\r\n", flags)
	if flags != nil {
		fmt.Printf("    flags = %+v\r\n", flags.LVal())
	}
	fmt.Printf("targetFrameName = %+v\r\n", targetFrameName)
	if targetFrameName != nil && targetFrameName.BstrVal() != nil {
		fmt.Printf("  targetFrameName = %+v\r\n", win.BSTRToString(targetFrameName.BstrVal()))
	}
	fmt.Printf("postData = %+v\r\n", postData)
	if postData != nil {
		fmt.Printf("    postData = %+v\r\n", postData.PVarVal())
	}
	fmt.Printf("headers = %+v\r\n", headers)
	if headers != nil && headers.BstrVal() != nil {
		fmt.Printf("  headers = %+v\r\n", win.BSTRToString(headers.BstrVal()))
	}
	fmt.Printf("cancel = %+v\r\n", cancel)
	if cancel != nil {
		fmt.Printf("  *cancel = %+v\r\n", *cancel)
	}

	if url != nil && url.BstrVal() != nil {
		// URLを格納する
		// このURLはfile:///C:/でなくC:\ で渡ってくるので注意
		threadWin.url = win.BSTRToString(url.BstrVal())
	}
}

func (threadWin *ThreadWin) webView_NavigateComplete2(pDisp *win.IDispatch, url *win.VARIANT) {
	fmt.Printf("webView_NavigateComplete2\r\n")
	fmt.Printf("pDisp = %+v\r\n", pDisp)
	fmt.Printf("url = %+v\r\n", url)
	if url != nil && url.BstrVal() != nil {
		fmt.Printf("  url = %+v\r\n", win.BSTRToString(url.BstrVal()))
	}
}

func (threadWin *ThreadWin) webView_DownloadBegin() {
	fmt.Printf("webView_DownloadBegin\r\n")

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

func (threadWin *ThreadWin) webView_DownloadComplete() {
	fmt.Printf("webView_DownloadComplete\r\n")
}

func (threadWin *ThreadWin) webView_DocumentComplete(pDisp *win.IDispatch, url *win.VARIANT) {
	fmt.Printf("webView_DocumentComplete\r\n")
	fmt.Printf("pDisp = %+v\r\n", pDisp)
	fmt.Printf("url = %+v\r\n", url)
	if url != nil && url.BstrVal() != nil {
		fmt.Printf("  url = %+v\r\n", win.BSTRToString(url.BstrVal()))
	}
}

func (threadWin *ThreadWin) webView_NavigateError(
	pDisp *win.IDispatch,
	url *win.VARIANT,
	targetFrameName *win.VARIANT,
	statusCode *win.VARIANT,
	cancel *win.VARIANT_BOOL) {

	fmt.Printf("webView_NavigateError\r\n")
	fmt.Printf("pDisp = %+v\r\n", pDisp)
	fmt.Printf("url = %+v\r\n", url)
	if url != nil && url.BstrVal() != nil {
		fmt.Printf("  url = %+v\r\n", win.BSTRToString(url.BstrVal()))
	}
	fmt.Printf("targetFrameName = %+v\r\n", targetFrameName)
	if targetFrameName != nil && targetFrameName.BstrVal() != nil {
		fmt.Printf("  targetFrameName = %+v\r\n", win.BSTRToString(targetFrameName.BstrVal()))
	}
	fmt.Printf("statusCode = %+v\r\n", statusCode)
	if statusCode != nil {
		fmt.Printf("    statusCode = %+v\r\n", statusCode.LVal())
	}
	fmt.Printf("cancel = %+v\r\n", cancel)
	if cancel != nil {
		fmt.Printf("  *cancel = %+v\r\n", *cancel)
	}
}

func (threadWin *ThreadWin) webView_NewWindow3(
	ppDisp **win.IDispatch,
	cancel *win.VARIANT_BOOL,
	dwFlags uint32,
	bstrUrlContext *uint16,
	bstrUrl *uint16) {

	fmt.Printf("webView_NewWindow3\r\n")
	fmt.Printf("ppDisp = %+v\r\n", ppDisp)
	if ppDisp != nil {
		fmt.Printf("    *ppDisp = %+v\r\n", *ppDisp)
	}
	fmt.Printf("cancel = %+v\r\n", cancel)
	if cancel != nil {
		fmt.Printf("  *cancel = %+v\r\n", *cancel)
	}
	fmt.Printf("dwFlags = %+v\r\n", dwFlags)
	fmt.Printf("bstrUrlContext = %+v\r\n", bstrUrlContext)
	if bstrUrlContext != nil {
		fmt.Printf("  bstrUrlContext = %+v\r\n", win.BSTRToString(bstrUrlContext))
	}
	fmt.Printf("bstrUrl = %+v\r\n", bstrUrl)
	if bstrUrl != nil {
		fmt.Printf("  bstrUrl = %+v\r\n", win.BSTRToString(bstrUrl))
	}
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
