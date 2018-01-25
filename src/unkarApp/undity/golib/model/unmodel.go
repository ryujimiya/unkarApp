package unmodel

// 基準モデル

import (
	"../get2ch"
	"time"
)

const (
	ClassNameServer  = "server"
	ClassNameBoard   = "board"
	ClassNameThread  = "thread"
	ClassNameSpecial = "special"
	ClassNameNone    = "none"
)

type ModelComponent struct {
	HostUrl   string
	url       string
	title     string
	className string
	mod       time.Time
	err       error
	g2ch      *get2ch.Get2ch
}

func CreateModelComponent(name, host string) ModelComponent {
	m := ModelComponent{
		HostUrl:   host,
		className: name,
	}
	return m
}

func (mc *ModelComponent) GetClassName() string { return mc.className }
func (mc *ModelComponent) GetTitle() string     { return mc.title }
func (mc *ModelComponent) GetMod() time.Time    { return mc.mod }
func (mc *ModelComponent) GetError() error      { return mc.err }
func (mc *ModelComponent) GetUrl() string       { return mc.url }
func (mc *ModelComponent) Is404() bool          { return mc.g2ch.Is404() }
func (mc *ModelComponent) GetServer() string    { return mc.g2ch.GetServer("") }
func (mc *ModelComponent) GetByteSize() int64   { return mc.g2ch.GetByteSize() }
func (mc *ModelComponent) GetCode() int         { return mc.g2ch.GetHttpCode() }
