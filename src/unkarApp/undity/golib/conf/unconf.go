package unconf

// 設定一覧

import (
	"regexp"
	"time"
)

const (
	Ver               = "6.80"
	ResMax            = 1010
	ResMin            = 1
	AffiliateAmazonId = "unkar-22"
	OneYearSec        = 31104000 * time.Second // だいたい一年
)

var RegInitThread = regexp.MustCompile(`^\/(\w+)\/(\d{9,10})(.*)?`)
var RegInitBoard = regexp.MustCompile(`^\/(\w+)(\/((?:no|s[ip]|re)2?)?)?$`)
var RegInitSpecial = regexp.MustCompile(`^\/\*(?:\/(.*))?`)
var RegInitSpecialSearch = regexp.MustCompile(`^\/\*\/search\/([^\/]+)`)
var RegInitServer = regexp.MustCompile(`^\/?$`)
var RegThreadAttrLastn = regexp.MustCompile(`^\/l(\d{1,3})$`)
var RegThreadAttrBottom = regexp.MustCompile(`^\/(\d{1,4})-$`)
var RegThreadAttrTop = regexp.MustCompile(`^\/-(\d{1,4})$`)
var RegThreadAttrRes = regexp.MustCompile(`^\/([-,\d]+)$`)
var RegThreadAttrId = regexp.MustCompile(`^\/ID:([\w!\+\/]+)$`)
var RegThreadAttrAnchor = regexp.MustCompile(`^\/Anchor:(Default|@\d{1,4}!\d{1,4})$`)
var RegThreadAttrLink = regexp.MustCompile(`^\/Link:((?:Imag|Movi)e|A(?:rchive|ll)|Thread)$`)
var RegThreadAttrTree = regexp.MustCompile(`^\/Tree:([-,\d]+|ID:[\w!\+\/]+|Link:(?:(?:Imag|Movi)e|A(?:rchive|ll)|Thread))$`)
var RegUrl = regexp.MustCompile(`(?:s?h?ttps?|sssp):\/\/[-_.!~*'()\w;\/?:\@&=+\$,%#\|]+`)
var RegSpace = regexp.MustCompile(`\s+`)

var ServerKill = map[string]bool{
	"www.2ch.net":         true,
	"info.2ch.net":        true,
	"find.2ch.net":        true,
	"v.isp.2ch.net":       true,
	"m.2ch.net":           true,
	"test.up.bbspink.com": true,
	"stats.2ch.net":       true,
	"c-au.2ch.net":        true,
	"c-others1.2ch.net":   true,
	"movie.2ch.net":       true,
	"img.2ch.net":         true,
	"ipv6.2ch.net":        true,
	"be.2ch.net":          true,
	"p2.2ch.net":          true,
	"c.2ch.net":           true,
}

var BoardKill = map[string]bool{
	"test": true,
}
