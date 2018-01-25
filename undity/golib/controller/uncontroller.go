package uncontroller

// コントローラー

import (
	"../conf"
	"../model"
	"../util"
	"../view"
	"net/http"
	"regexp"
)

const (
	defaultModelKey = 3
	defaultViewKey  = "unkar02"
)

type Route struct {
	Regs  *regexp.Regexp
	Model func(host string, path []string) unutil.Model
}

var modelList []*Route
var viewList map[string]func(path string, r *http.Request) unutil.View

func init() {
	viewList = make(map[string]func(path string, r *http.Request) unutil.View)

	modelList = []*Route{
		&Route{
			Regs: unconf.RegInitThread,
			Model: func(host string, path []string) unutil.Model {
				return unmodel.NewThread(host, path)
			},
		},
		&Route{
			Regs: unconf.RegInitBoard,
			Model: func(host string, path []string) unutil.Model {
				return unmodel.NewBoard(host, path)
			},
		},
		&Route{
			Regs: unconf.RegInitSpecial,
			Model: func(host string, path []string) unutil.Model {
				return unmodel.NewSpecial(host, path)
			},
		},
		&Route{
			Regs: unconf.RegInitServer,
			Model: func(host string, path []string) unutil.Model {
				return unmodel.NewServer(host, path)
			},
		},
	}

	viewList[defaultViewKey] = unview.NewViewContainer
}

func Dispatch(r *http.Request, path string, viewname string) unutil.Output {
	var model unutil.Model
	flag := false

	viewfunc, ok := viewList[viewname]
	if !ok {
		viewfunc = viewList[defaultViewKey]
	}
	view := viewfunc(path, r)
	host := view.GetHostUrl()

	for _, it := range modelList {
		if match := it.Regs.FindStringSubmatch(path); match != nil {
			model = it.Model(host, match)
			flag = true
			break
		}
	}
	if flag == false {
		model = unmodel.NewNone(host, []string{path})
	}

	// 出力を返す
	return view.PrintData(model)
}
