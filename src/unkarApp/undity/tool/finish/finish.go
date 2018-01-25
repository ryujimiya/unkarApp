package main

import (
	"log"
	"net/http"
	"runtime/debug"
	"time"
)

type untidyHandle struct {
	fs http.Handler
}

const (
	TIMEOUT_HANDLER_SEC = 25 * time.Second
	TIMEOUT_READ_SEC    = 15 * time.Second
	TIMEOUT_WRITE_SEC   = 30 * time.Second
	MAX_HEADER_SIZE     = 1024 * 100
)

const message = `<!DOCTYPE html>
<html lang="ja">
<head>
<meta charset="utf-8">
<title>unkarだったらしい</title>
</head>
<body>
<p>unkarは閉鎖しました。<br>
今までありがとうございました。</p>
<p><a href="http://d.hatena.ne.jp/heiwaboke/20140723/1406121144">unkarを閉鎖します</a></p>
</body>
</html>`

func main() {
	debug.SetMaxThreads(512) // 512Thread

	server := &http.Server{
		Addr: ":80",
		Handler: http.TimeoutHandler(&untidyHandle{
			fs: http.FileServer(http.Dir("public_html")),
		}, TIMEOUT_HANDLER_SEC, "タイムアウトしました。"),
		ReadTimeout:    TIMEOUT_READ_SEC,
		WriteTimeout:   TIMEOUT_WRITE_SEC,
		MaxHeaderBytes: MAX_HEADER_SIZE,
	}
	log.Printf("listen start %s\n", server.Addr)
	// サーバ起動
	log.Fatal(server.ListenAndServe())
}

func (uh *untidyHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if r.Method != "HEAD" {
		w.Write([]byte(message))
	}
}
