package dcadb

import (
	"log"
	"testing"
)

func TestDcAdb_LaunchApp(t *testing.T) {
	log.Println("",NewDcAdb("adb").CheckPortIsListen("55272"))
}