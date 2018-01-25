package unmodel

// 存在しないページ

import (
	"../get2ch"
	"time"
)

type None struct {
	ModelComponent
}

func NewNone(host string, _ []string) *None {
	model := &None{ModelComponent: CreateModelComponent(ClassNameNone, host)}
	model.url = ""
	model.title = "そんなページないよ"
	model.mod = time.Time{}
	model.err = nil
	model.g2ch = get2ch.NewGet2ch("", "")
	return model
}

func (this *None) GetData() interface{} { return nil }
