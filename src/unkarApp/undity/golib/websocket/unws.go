package unws

// websocketライブラリ

import (
	"./top"
	"errors"
	"net/http"
	"regexp"
)

type WsPacket struct {
	path    string
	handler http.Handler
	rch     chan<- http.Handler
}

type WsHandle struct {
	wsCh     chan<- WsPacket
	exitCh   chan<- string
	notfound http.Handler
}

type WsRoute struct {
	Path *regexp.Regexp
	New  func(string, string, chan<- string) http.Handler
}

var routeMap map[string]*WsRoute

func init() {
	routeMap = make(map[string]*WsRoute)
	routeMap["unkartop"] = &WsRoute{
		Path: regexp.MustCompile(`^\/unkartop$`),
		New:  unkartop.New,
	}
}

func searchRoute(path string) (r *WsRoute, name string, err error) {
	for key, it := range routeMap {
		if it.Path.MatchString(path) {
			name = key
			r = it
			return
		}
	}
	return nil, "", errors.New("not found")
}

func WsInit() *WsHandle {
	wch := make(chan WsPacket, 4)
	ch := make(chan string, 1)
	dh := &WsHandle{
		wsCh:     wch,
		exitCh:   ch,
		notfound: http.NotFoundHandler(),
	}
	go func(ch <-chan string, wch <-chan WsPacket) {
		webm := make(map[string]http.Handler)
		for {
			select {
			case key := <-ch:
				delete(webm, key)
			case it := <-wch:
				_, key, err := searchRoute(it.path)
				if err == nil {
					if it.rch != nil {
						wh := webm[key]
						it.rch <- wh
					} else if it.handler != nil {
						webm[key] = it.handler
					} else {
						delete(webm, key)
					}
				}
			}
		}
	}(ch, wch)
	return dh
}

func (dh *WsHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	dh.match(r.URL.Path).ServeHTTP(w, r)
}

func (dh *WsHandle) match(p string) http.Handler {
	ch := make(chan http.Handler, 1)
	dh.wsCh <- WsPacket{
		path: p,
		rch:  ch,
	}
	wh := <-ch
	close(ch)

	if wh == nil {
		h, err := dh.createWebHandler(p)
		if err == nil {
			wh = h
			// 登録
			dh.wsCh <- WsPacket{
				path:    p,
				handler: wh,
			}
		} else {
			wh = dh.notfound
		}
	}
	return wh
}

func (dh *WsHandle) createWebHandler(path string) (h http.Handler, err error) {
	if r, key, e := searchRoute(path); e == nil {
		h = r.New(path, key, dh.exitCh)
	} else {
		err = errors.New("createWebHandler error")
	}
	return
}
