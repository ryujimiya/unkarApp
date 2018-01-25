package main

import (
	"./undity"
	"./undity/golib/util"
	//"./undity/golib/model"
	"fmt"
	"strings"
	//"log"
	"os"
	"path/filepath"
)

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

//import (
//	"github.com/lxn/win"
//)


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
	// 本文表示用エディットボックス
	//textEditBody *walk.TextEdit
	// 板名
	boardName string
	// スレッド番号
	threadNo int64
}

/**
 * コンストラクタ
 * @param なし
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
	// 本文作成
	htmlText := CreateHtmlText(boardName, threadNo)

	// 実行ファイルパスを取得
	exePath, _ := os.Executable()
	exeDir := filepath.Dir(exePath)

	// 表示用htmlファイルに保存
	tmpHtmlFilename := fmt.Sprintf("%s_%d.html", boardName, threadNo)
	tmpHtmlFilepath := exeDir + "\\tmp\\" + tmpHtmlFilename
	fmt.Printf("tmpHtmlFilepath=" + tmpHtmlFilepath + "\r\n")
	SaveToTmpHtml(tmpHtmlFilepath, htmlText)

	// 表示用htmlのURL
	exeDir = strings.Replace(exeDir, "\\", "/", -1)
	tmpHtmlUrl := "file:///" + exeDir + "/tmp/" + tmpHtmlFilename
	fmt.Printf("tmpHtmlUrl=" + tmpHtmlUrl + "\r\n")

	// メインウィンドウのウィンドウ生成
	err := MainWindow {
		AssignTo:	&threadWin.MainWindow,
		Title:	"Unkar App",
		MinSize:	Size{850, 600},
		Layout:	VBox{},
		Children: []Widget {
			WebView{
				AssignTo: &threadWin.webView,
				Name:     "スレッド",
				URL: tmpHtmlUrl,
			},
			/*
			TextEdit {
				AssignTo: &threadWin.textEditBody,
				ReadOnly: true,
				Text: htmlText,
			},
			*/
		},
	}.Create()
	
	// 子ウィンドウ化
	// Note:win.SetParentだと、親WindowにクリッピングされたWindowになってしまう
	//win.SetParent(threadWin.Handle(), parentWin.Handle())
	// FormのSetOwnerを使うといいみたい
	threadWin.SetOwner(parentWin)

	return threadWin, err
}

func CreateHtmlText(boardName string, threadNo int64) string {
	// 本文
	var htmlText string
	// unutilのモデル
	var unutilModel unutil.Model

	// 初期化
	htmlText = ""

	// スレッドのモデルを取得する
	unutilModel = unkarstub.GetThreadModel(boardName, threadNo)

	// DEBUG
	fmt.Printf("url=%s\r\n", unutilModel.GetUrl())
	fmt.Printf("title=%s\r\n", unutilModel.GetTitle())
	fmt.Printf("className=%s\r\n", unutilModel.GetClassName())
	fmt.Printf("server=%s\r\n", unutilModel.GetServer())

    htmlText = unkarstub.GetBoardViewOutput(boardName, threadNo, unutilModel)

	return htmlText
}

func SaveToTmpHtml(filename, htmlStr string) {
	out, err := os.Create(filename)
	if err != nil {
		return
	}
	fmt.Fprint(out, htmlStr)
	out.Close()
}

