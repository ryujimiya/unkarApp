package unmodel

// スペシャルモデル

import (
	"../get2ch"
	"time"
)

type Special struct {
	ModelComponent
}

func NewSpecial(host string, _ []string) *Special {
	model := &Special{ModelComponent: CreateModelComponent(ClassNameSpecial, host)}
	model.url = ""
	model.title = "特殊なページ"
	model.mod = time.Time{}
	model.err = nil
	model.g2ch = get2ch.NewGet2ch("", "")
	return model
}

func (this *Special) GetData() interface{} { return nil }
