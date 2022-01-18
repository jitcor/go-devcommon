package dcshell

import (
	"log"
	"testing"
)

func TestNewDcShell(t *testing.T) {
	log.Println("rrr:",NewDcShell(false).LaunchApp("bin.mt.plus"))
}