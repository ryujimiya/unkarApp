package main

import (
	"log"
)

////////////////////////////////////////////////////////////
// unkarApp
////////////////////////////////////////////////////////////
/**
 * アプリケーションエントリーポイント
 */
func main() {
	var mw *MainWin
	var err error

	// メインウィンドウの生成
	mw, err = NewMainWin()
	if err != nil {
		log.Fatal(err)
		return
	}
	// 表示
	mw.Run()
}

