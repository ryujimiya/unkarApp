package untidy

import (
	"./controller"
	"./util"
	"net/http"
)

func Start(path string, r *http.Request) unutil.Output {
	// 処理割り振り
	return uncontroller.Dispatch(r, path, "unkar02")
}
