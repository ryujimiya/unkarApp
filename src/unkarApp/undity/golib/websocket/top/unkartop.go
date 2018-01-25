package unkartop

// うんかーのTOPコマンド

import (
	"../../get2ch"
	"../../get2ch/process"
	"../../util"
	"code.google.com/p/go.net/websocket"
	"net/http"
	"sync"
	"time"
)

const (
	MaxProcs                   = 64
	MaxProcsCh                 = 16
	IgnitionTime time.Duration = 2 * time.Second
	DeadlineTime time.Duration = 8 * time.Second
)

type UnkarTop struct {
	sync.RWMutex
	path   string
	key    string
	listen map[*websocket.Conn]chan<- error
}

func New(path, key string, exitCh chan<- string) http.Handler {
	ut := UnkarTop{
		path:   path,
		key:    key,
		listen: make(map[*websocket.Conn]chan<- error),
	}
	f := func(ws *websocket.Conn) {
		ut.RLock()
		l := len(ut.listen)
		ut.RUnlock()

		if l > MaxProcs {
			// 接続数オーバー
			return
		}

		ech := make(chan error)

		ut.Lock()
		ut.listen[ws] = ech
		ut.Unlock()

		// エラー発生まで待ちうけ
		<-ech

		ut.Lock()
		delete(ut.listen, ws)
		ut.Unlock()

		ws.Close()
	}
	go ut.timeCallback(exitCh)
	return websocket.Handler(f)
}

func (ut *UnkarTop) timeCallback(exitCh chan<- string) {
	defer ut.dispose(exitCh)

	var ps uint64
	c := time.Tick(IgnitionTime)
	for _ = range c {
		// 定期的に実行
		// 接続がない時は終了する
		if ut.isConnection() == false {
			break
		}
		// データの取得
		data, size := get2ch.GetViewThreadLot()
		if data == nil {
			continue
		}
		max := unutil.MinUint64(size-ps, process.VIEW_THREAD_LIST_SIZE)
		// データの送信
		ut.Send(data[:max])
		// 今回値の保存
		ps = size
	}
}

func (ut *UnkarTop) isConnection() (ret bool) {
	ut.RLock()
	ret = len(ut.listen) != 0
	ut.RUnlock()
	return
}

func (ut *UnkarTop) dispose(exitCh chan<- string) {
	// 元のコントローラーを停止
	exitCh <- ut.key
}

func (ut *UnkarTop) Send(data []process.ViewThreadItem) {
	dl := time.Now().Add(DeadlineTime)
	sch := make(chan bool, MaxProcsCh)

	ut.RLock()
	for con, ech := range ut.listen {
		sch <- true
		con.SetDeadline(dl)
		go writeData(data, con, ech, sch)
	}
	ut.RUnlock()

	for i := MaxProcsCh; i > 0; i-- {
		sch <- true
	}
	close(sch)
}

func writeData(data []process.ViewThreadItem, con *websocket.Conn, ech chan<- error, sch <-chan bool) {
	err := websocket.JSON.Send(con, data)
	<-sch
	if err != nil {
		ech <- err
	}
}
