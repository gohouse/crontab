package gobPersist

import (
	"encoding/gob"
	"errors"
	"github.com/gohouse/golib/file"
	"os"
)
var FileNotExists = errors.New("file not exists.")

type GobPersist struct{
	file string
}

func NewGobPersist(file string) *GobPersist {
	return &GobPersist{file: file}
}

func (gp *GobPersist) Store(key, arg interface{}) error {
	//f, err := os.Create(filename)
	f := file.NewFile(gp.file).OpenFile(os.O_CREATE | os.O_WRONLY | os.O_TRUNC)
	//if err != nil {
	//	panic(err.Error())
	//}

	enc := gob.NewEncoder(f)
	return enc.Encode(arg)
}

func (gp *GobPersist) Load(key, arg interface{}) error {
	if !file.FileExists(gp.file) {
		return FileNotExists
	}
	f, err := os.Open(gp.file)
	if err != nil {
		panic(err.Error())
	}
	dec := gob.NewDecoder(f)
	return dec.Decode(arg)
}
