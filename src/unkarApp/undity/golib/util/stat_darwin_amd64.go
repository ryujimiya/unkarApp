package unutil

import (
	"os"
)

func Stat(filename string) (*Stat_t, error) {
	s, err := os.Stat(filename)
	if err != nil {
		return nil, err
	}
	stat := &Stat_t{
		Size:  s.Size(),
		Atime: s.ModTime(),
		Mtime: s.ModTime(),
	}
	return stat, nil
}
