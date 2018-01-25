package get2ch

import (
	"../util"
	"io"
	"io/ioutil"
	"os"
	"path"
	"time"
	"fmt" // DEBUG
)

const (
	BOARD_SETTING = "setting"
	// 板情報格納ファイル
	tBOARD_LIST_NAME = "ita.data"
	// スレッド一覧格納ファイル名
	tBOARD_SUBJECT_NAME = "subject.txt"
	// 板情報格納ファイル名
	tBOARD_SETTING_NAME = "setting.txt"
)

type State struct {
	fsize int64
	atime time.Time
	mtime time.Time
}

func (s *State) Size() int64     { return s.fsize }
func (s *State) Amod() time.Time { return s.atime }
func (s *State) Mmod() time.Time { return s.mtime }

type FileCache struct {
	Folder string // dat保管フォルダ名
}

func NewFileCache(root string) *FileCache {
	return &FileCache{
		Folder: root,
	}
}

func (fc FileCache) Path(s, b, t string) string {
	if s == "" && b == "" && t == "" {
		return fc.Folder + "/" + tBOARD_LIST_NAME
	} else if t == BOARD_SETTING {
		return fc.Folder + "/" + b + "/" + tBOARD_SETTING_NAME
	} else if t == "" {
		return fc.Folder + "/" + b + "/" + tBOARD_SUBJECT_NAME
	}
	return fc.Folder + "/" + b + "/" + t[0:4] + "/" + t + ".dat"
}

func (fc FileCache) GetData(s, b, t string) ([]byte, error) {
	fmt.Printf("FileCache::GetData: s=%s b=%s t=%s\r\n", s, b, t) // DEBUG
	logfile := fc.Path(s, b, t)
	fmt.Print(logfile) // DEBUG
	fmt.Printf("\r\n") // DEBUG
	return ioutil.ReadFile(logfile)
}
func (fc FileCache) GetDataRC(s, b, t string) (io.ReadCloser, error) {
	fmt.Printf("FileCache::GetDataRC: s=%s b=%s t=%s\r\n", s, b, t) // DEBUG
	logfile := fc.Path(s, b, t)
	fmt.Print(logfile) // DEBUG
	fmt.Printf("\r\n") // DEBUG
	rc, err := os.Open(logfile)
	if err != nil {
		fmt.Printf("failed\r\n") // DEBUG
		return nil, err
	}
	fmt.Printf("success\r\n") // DEBUG
	return rc, nil
}

func (fc FileCache) SetData(s, b, t string, d []byte) error {
	fmt.Printf("FileCache::SetData: s=%s b=%s t=%s\r\n", s, b, t) // DEBUG
	logfile := fc.Path(s, b, t)
	fmt.Print(logfile) // DEBUG
	fmt.Printf("\r\n") // DEBUG
	os.MkdirAll(path.Dir(logfile), 0777)
	return ioutil.WriteFile(logfile, d, 0777)
}

func (fc FileCache) SetDataAppend(s, b, t string, d []byte) error {
	fmt.Printf("FileCache::SetDataAppend: s=%s b=%s t=%s\r\n", s, b, t) // DEBUG
	logfile := fc.Path(s, b, t)
	fmt.Print(logfile) // DEBUG
	fmt.Printf("\r\n") // DEBUG
	fp, err := os.OpenFile(logfile, os.O_WRONLY|os.O_APPEND, 0777)
	if err != nil {
		fmt.Printf("failed\r\n") // DEBUG
		return err
	}
	_, err = fp.Write(d)
	fp.Close()
	fmt.Printf("success\r\n") // DEBUG
	return err
}

func (fc FileCache) SetMod(s, b, t string, m, a time.Time) error {
	// atimeとmtimeの順番に注意
	return os.Chtimes(fc.Path(s, b, t), a, m)
}

func (fc FileCache) Exists(s, b, t string) bool {
	_, err := os.Stat(fc.Path(s, b, t))
	return err == nil
}

func (fc FileCache) Stat(s, b, t string) (CacheState, error) {
	st, err := unutil.Stat(fc.Path(s, b, t))
	if err != nil {
		return nil, err
	}
	return &State{
		fsize: st.Size,
		atime: st.Atime,
		mtime: st.Mtime,
	}, nil
}
