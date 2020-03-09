package fileLog

import (
	"fmt"
	"github.com/gohouse/golib/date"
	"github.com/gohouse/golib/file"
	"os"
)

type FileLog struct {
	file string
	fp   *os.File
}

func NewFileLog(filename string) *FileLog {
	fp := file.NewFile(filename).OpenFile()
	return &FileLog{file: filename, fp: fp}
}

func (fl *FileLog) Info(arg string) (err error) {
	_, err = fmt.Fprintf(fl.fp, "\n[%s] [info] %s", date.NewDate().TodayDateTime(), arg)
	return
}

func (fl *FileLog) Infof(arg string, args ...interface{}) (err error) {
	return fl.Info(fmt.Sprintf(arg, args...))
}

func (fl *FileLog) Error(errs error) (err error) {
	_, err = fmt.Fprintf(fl.fp, "\n[%s] [error] %s", date.NewDate().TodayDateTime(), errs.Error())
	return
}

func (fl *FileLog) Read(readRows ...int64) string {
	var rows int64 = 20
	if len(readRows) > 0 {
		rows = readRows[0]
	}
	return file.Tail_f(fl.file, rows)
}
