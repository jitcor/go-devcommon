package dcsocket

import (
	"github.com/Humenger/go-devcommon/dcdebug"
	"testing"
)

func TestGetUDPConn(t *testing.T) {
	if conn,err:= GetUDPConn(48081);err!=nil{
		t.Error(err)
	}else {
		dcdebug.LogI("conn",conn)
	}
}