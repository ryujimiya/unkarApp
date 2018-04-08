package process

import (
	"../../search"
	"../../util"
	"../../util/kill"
	"errors"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

const (
	BOURBON_TIME           = 1 * time.Minute
	SERVER_UPDATE_CYCLE    = 1 * time.Hour
	BOARDNAME_UPDATE_CYCLE = 24 * time.Hour
	DB_INSERT_CYCLE        = 1 * time.Minute
	DB_INSERT_LIST_BUFSIZE = 32
	DB_INSERT_BUFSIZE      = 128
	DB_UPDATE_BUFSIZE      = 256
	VIEW_THREAD_LIST_SIZE  = 100
)

type ViewThreadItem struct {
	Board     string    `json:"board"`
	Thread    string    `json:"thread"`
	Boardname string    `json:"boardname"`
	Title     string    `json:"title"`
	Res       int       `json:"res"`
	Lastdate  time.Time `json:"lastdate"`
	Addtime   time.Time `json:"addtime"`
}

type BoardServerBox struct {
	sync.RWMutex
	m map[string]string
	f func() map[string]string
	t time.Time
}

func NewBoardServerBox(f func() map[string]string) *BoardServerBox {
	return &BoardServerBox{
		m: f(),
		f: f,
		t: time.Now().Add(SERVER_UPDATE_CYCLE),
	}
}
func (bs *BoardServerBox) GetServer(board string) (server string) {
	bs.RLock()
	server = bs.m[board]
	old := bs.t
	bs.RUnlock()

	now := time.Now()
	if old.Before(now) {
		// 今よりも古い時間になった場合更新
		m := bs.f()
		t := now.Add(SERVER_UPDATE_CYCLE)

		bs.Lock()
		bs.m = m
		bs.t = t
		bs.Unlock()
	}
	return
}

type boardNameItem struct {
	name string
	t    time.Time
}
type BoardNameBox struct {
	sync.RWMutex
	m map[string]boardNameItem
}

func NewBoardNameBox() *BoardNameBox {
	bn := &BoardNameBox{
		m: make(map[string]boardNameItem, 1024),
	}
	return bn
}
func (bn *BoardNameBox) SetName(board, bname string) {
	item := boardNameItem{
		name: bname,
		t:    time.Now().Add(BOARDNAME_UPDATE_CYCLE),
	}

	bn.Lock()
	bn.m[board] = item
	bn.Unlock()
}
func (bn *BoardNameBox) GetName(board string) string {
	bn.RLock()
	item := bn.m[board]
	bn.RUnlock()

	if item.name != "" {
		now := time.Now()
		if item.t.Before(now) {
			bn.Lock()
			delete(bn.m, board)
			bn.Unlock()
		}
	}
	return item.name
}

type BBNCacheBox struct {
	sync.RWMutex
	cm map[string]time.Time
}

func NewBBNCacheBox() *BBNCacheBox {
	bbn := &BBNCacheBox{
		cm: make(map[string]time.Time),
	}
	return bbn
}
func (bbn *BBNCacheBox) SetBourbon(key string) {
	bbn.Lock()
	bbn.cm[key] = time.Now().Add(BOURBON_TIME)
	bbn.Unlock()
}
func (bbn *BBNCacheBox) GetBourbon(key string) bool {
	bbn.RLock()
	t, ok := bbn.cm[key]
	bbn.RUnlock()
	if ok {
		if time.Now().After(t) {
			// 期間経過
			bbn.Lock()
			delete(bbn.cm, key)
			bbn.Unlock()
			ok = false
		}
	}
	return ok
}

type ViewThreadBox struct {
	sync.RWMutex
	l [VIEW_THREAD_LIST_SIZE]ViewThreadItem
	c uint64
}

func NewViewThreadBox() *ViewThreadBox {
	return &ViewThreadBox{}
}
func (vt *ViewThreadBox) SetThread(item ViewThreadItem) {
	vt.Lock()
	// データを設定
	copy(vt.l[1:], vt.l[0:])
	vt.l[0] = item
	vt.Unlock()

	atomic.AddUint64(&vt.c, 1)
}
func (vt *ViewThreadBox) GetThreadList(start, end int) (ret []ViewThreadItem) {
	max := unutil.MinUintV64(uint64(end), VIEW_THREAD_LIST_SIZE, atomic.LoadUint64(&vt.c))
	if uint64(start) < max {
		ret = make([]ViewThreadItem, max-uint64(start))

		vt.RLock()
		copy(ret, vt.l[start:])
		vt.RUnlock()
	}
	return
}
func (vt *ViewThreadBox) GetThreadLot() (ret []ViewThreadItem, c uint64) {
	c = atomic.LoadUint64(&vt.c)
	end := unutil.MinUint64(VIEW_THREAD_LIST_SIZE, c)
	ret = make([]ViewThreadItem, end)

	vt.RLock()
	copy(ret, vt.l[:])
	vt.RUnlock()
	return
}

type DBInsertBox struct {
	wch chan<- *unsearch.DBItem
}

func NewDBInsertBox() *DBInsertBox {
	ch := make(chan *unsearch.DBItem, DB_INSERT_BUFSIZE)
	dbin := &DBInsertBox{
		wch: ch,
	}
	go func(rch <-chan *unsearch.DBItem) {
		in := unsearch.NewInsert(DB_INSERT_LIST_BUFSIZE)
		c := time.Tick(DB_INSERT_CYCLE)
		kc := kill.CreateKillChan()
		for {
			select {
			case it := <-rch:
				// バッファにためる
				in.Push(it)
			case <-kc:
				// 外部から殺された
				in.Exec()
				os.Exit(1)
			case <-c:
				// 時間により実行
				in.Exec()
			}
		}
	}(ch)
	return dbin
}
func (in *DBInsertBox) SetItem(data []byte, board, thread string, resnum int) {
	item := unsearch.CreateDBItem(data, board, thread, resnum)
	if item != nil {
		in.wch <- item
	}
}

type dbUpdateItem struct {
	Board  string
	Number string
	Resnum int
}
type DBUpdateBox struct {
	wch chan<- *dbUpdateItem
}

func NewDBUpdateBox() *DBUpdateBox {
	ch := make(chan *dbUpdateItem, DB_UPDATE_BUFSIZE)
	dbup := &DBUpdateBox{
		wch: ch,
	}
	go func(rch <-chan *dbUpdateItem) {
		up := unsearch.NewUpdate()
		for it := range rch {
			up.Update(it.Board, it.Number, it.Resnum)
		}
	}(ch)
	return dbup
}
func (up *DBUpdateBox) SetItem(board, thread string, resnum int) {
	item := &dbUpdateItem{
		Board:  board,
		Number: thread,
		Resnum: resnum,
	}
	up.wch <- item
}

type BoardThreadMap struct {
	sync.RWMutex
	m map[string]map[int64]struct{}
}

func NewBoardThreadMap() *BoardThreadMap {
	return &BoardThreadMap{
		m: make(map[string]map[int64]struct{}, 1024),
	}
}
func (btm *BoardThreadMap) Lookup(board, thread string) (bool, error) {
	num, nerr := strconv.ParseInt(thread, 10, 64)
	if nerr != nil {
		return false, nerr
	}

	btm.RLock()
	bm, ok := btm.m[board]
	btm.RUnlock()

	var ret bool
	var err error
	if ok {
		if _, ok := bm[num]; ok {
			ret = true
		} else {
			ret = false
		}
	} else {
		ret = true
		err = errors.New("not found")
	}
	return ret, err
}
func (btm *BoardThreadMap) SetMap(board string, m map[int64]struct{}) {
	btm.Lock()
	btm.m[board] = m
	btm.Unlock()
}
