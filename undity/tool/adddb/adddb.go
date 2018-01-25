package main

/*

スレッド情報の挿入

CREATE DATABASE unkar DEFAULT CHARACTER SET utf8;

CREATE TABLE thread_title(
id INT PRIMARY KEY AUTO_INCREMENT,
board VARCHAR(32),
number INT UNSIGNED,
title VARCHAR(100),
master VARCHAR(100),
resnum SMALLINT UNSIGNED,
INDEX board_index (board),
INDEX number_index (number),
FULLTEXT INDEX title_index (title) COMMENT 'parser "TokenMecab"'
) ENGINE = Mroonga DEFAULT CHARSET utf8;

*/

import (
	"../../golib/search"
	"io/ioutil"
	"log"
	"os"
)

const (
	BASE_PATH      = "/2ch/dat"
	CH_BUFF_SIZE   = 256
	UNIT_SIZE      = 64
	BOARD_PARALLEL = 8
)

var stdlog *log.Logger = log.New(os.Stdout, "", log.LstdFlags)

func main() {
	dat, err := ioutil.ReadDir(BASE_PATH)
	if err != nil {
		return
	}
	ch, exitCh := registerProc()
	for _, it := range dat {
		if it.IsDir() {
			name := it.Name()
			board(name, ch)
		}
	}
	ch <- &unsearch.DBItem{}
	<-exitCh
}

func registerProc() (chan<- *unsearch.DBItem, <-chan bool) {
	exitCh := make(chan bool, 1)
	ch := make(chan *unsearch.DBItem, CH_BUFF_SIZE)
	go func(ch <-chan *unsearch.DBItem, exitCh chan<- bool) {
		in := unsearch.NewInsert(UNIT_SIZE)
		for it := range ch {
			if it.Resnum == 0 {
				in.Exec()
				exitCh <- true
				break
			}
			in.Push(it)
		}
	}(ch, exitCh)
	return ch, exitCh
}

func board(boardname string, ch chan<- *unsearch.DBItem) {
	board_path := BASE_PATH + "/" + boardname
	board, boarderr := ioutil.ReadDir(board_path)
	if boarderr != nil {
		return
	}

	parallel := BOARD_PARALLEL
	sync := make(chan bool, parallel)
	for _, it := range board {
		if it.IsDir() == false {
			continue
		}
		indexname := it.Name()
		index_path := board_path + "/" + indexname
		index, indexerr := ioutil.ReadDir(index_path)
		if indexerr != nil {
			break
		}

		sync <- true
		go func(index []os.FileInfo, index_path, boardname string, ch chan<- *unsearch.DBItem) {
			indexLoop(index, index_path, boardname, ch)
			stdlog.Println(index_path)
			<-sync
		}(index, index_path, boardname, ch)
	}
	for ; parallel > 0; parallel-- {
		sync <- true
	}
	close(sync)
}

func indexLoop(index []os.FileInfo, index_path, boardname string, ch chan<- *unsearch.DBItem) {
	for _, line := range index {
		if line.IsDir() {
			continue
		}
		// ディレクトリではない
		thread := line.Name()
		tlen := len(thread) - 4
		if thread[tlen:] == ".dat" {
			p := index_path + "/" + thread
			data, err := ioutil.ReadFile(p)
			if err == nil {
				it := unsearch.CreateDBItem(data, boardname, thread[:tlen], 0)
				if it != nil {
					ch <- it
				}
			}
		}
	}
}
