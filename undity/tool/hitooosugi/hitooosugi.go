package main

// 人大杉datの削除

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
)

const RootDir = "/2ch/dat"

func main() {
	parallel := 4
	moto, err := ioutil.ReadFile("hitooosugi.txt")
	if err != nil {
		return
	}
	motosize := int64(len(moto))
	fmt.Printf("hitooosugi size: %dbyte\n", motosize)

	boardlist, err := ioutil.ReadDir(RootDir)
	if err != nil {
		return
	}
	sync := make(chan bool, parallel)
	for _, it := range boardlist {
		if it.IsDir() {
			fmt.Println(it.Name())
			sync <- true
			go func(it os.FileInfo) {
				board(RootDir+"/"+it.Name(), moto, motosize)
				<-sync
			}(it)
		}
	}
	for ; parallel > 0; parallel-- {
		sync <- true
	}
	close(sync)
}

func board(indexpath string, moto []byte, motosize int64) {
	indexlist, err := ioutil.ReadDir(indexpath)
	if err != nil {
		return
	}
	for _, index := range indexlist {
		if !index.IsDir() {
			continue
		}
		datpath := indexpath + "/" + index.Name()
		datlist, err := ioutil.ReadDir(datpath)
		if err != nil {
			continue
		}
		for _, dat := range datlist {
			if dat.IsDir() {
				continue
			}
			name := dat.Name()
			l := len(name)
			if l >= 13 && name[l-4:] == ".dat" {
				if dat.Size() == motosize {
					p := datpath + "/" + name
					f, err := ioutil.ReadFile(p)
					if err != nil {
						continue
					}
					if bytes.Equal(moto, f) {
						os.Remove(p)
						fmt.Printf("delete:%s\n", p)
					}
				}
			}
		}
	}
}
