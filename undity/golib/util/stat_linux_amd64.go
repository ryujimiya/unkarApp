package unutil

import (
	"syscall"
	"time"
)

func Stat(filename string) (*Stat_t, error) {
	var stat syscall.Stat_t
	err := syscall.Stat(filename, &stat)
	if err != nil {
		return nil, err
	}
	return &Stat_t{
		Size:  stat.Size,
		Atime: time.Unix(stat.Atim.Sec, stat.Atim.Nsec),
		Mtime: time.Unix(stat.Mtim.Sec, stat.Mtim.Nsec),
	}, nil
}
