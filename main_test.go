package devcommon

import (
	"encoding/hex"
	"log"
	"testing"
)

func TestGetUUID(t *testing.T) {
	if bUUID,err:=GetUUID();err!=nil{
		t.Error(err)
	}else {
		log.Println("hexUUID:",hex.EncodeToString(bUUID))
	}
}
