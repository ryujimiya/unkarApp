package unmodel

// サーバモデル

import (
	"../get2ch"
	"bufio"
	"strings"
)

type ItaItem struct {
	Server string
	Url    string
	Name   string
}

type ServerItem struct {
	Name    string
	ItaList []ItaItem
}

type Server struct {
	ModelComponent
	data []ServerItem
}

func NewServer(host string, _ []string) *Server {
	model := &Server{
		ModelComponent: CreateModelComponent(ClassNameServer, host),
	}
	model.url = ""
	model.g2ch = get2ch.NewGet2ch("", "")
	model.analyzeData()
	model.title = "板一覧"
	return model
}

func (this *Server) GetData() interface{} { return &this.data }

func (this *Server) analyzeData() {
	rc := this.g2ch.GetBBSmenu(true)
	this.mod = this.g2ch.GetModified()
	this.data = []ServerItem{}

	k := -1
	scanner := bufio.NewScanner(rc)
	for scanner.Scan() {
		it := scanner.Text()
		line := strings.Split(it, "<>")
		if len(line) == 2 {
			url := strings.Split(line[0], "/")
			if len(url) == 2 {
				this.data[k].ItaList = append(this.data[k].ItaList, ItaItem{
					Server: url[0],
					Url:    url[1],
					Name:   line[1],
				})
			}
		} else {
			if (k >= 0) && (len(this.data[k].ItaList) == 0) {
				// 何もしない
			} else {
				k++
			}
			this.data = append(this.data, ServerItem{
				Name:    it,
				ItaList: []ItaItem{},
			})
		}
	}
	rc.Close()

	if (k >= 0) && (len(this.data[k].ItaList) == 0) {
		this.data = this.data[:len(this.data)-1]
	}
	this.err = this.g2ch.GetError()
}
