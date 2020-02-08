package fileLog

import "testing"

func TestFileLog_Read(t *testing.T) {
	t.Log(NewFileLog("file.go").Read())
}
