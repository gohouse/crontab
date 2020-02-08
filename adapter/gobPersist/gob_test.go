package gobPersist

import "testing"

type User struct {
	Name string
	Age int
}
var ps Persist
func TestStore(t *testing.T) {
	var u = User{
		Name: "kevin",
		Age:  18,
	}
	err := ps.Store("xxx.gob",&u)
	if err!=nil {
		t.Error(err.Error())
		return
	}
	t.Log("store success: xxx.gob")
}

func TestLoad(t *testing.T) {
	var u User
	err := ps.Load("xx2x.gob",&u)
	if err!=nil {
		t.Error(err.Error())
		return
	}
	t.Logf("load success: %+v", u)
}
