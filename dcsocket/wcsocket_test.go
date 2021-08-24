package dcsocket

import "testing"

func TestStart(t *testing.T) {
	NewWCSocket("127.0.0.1",8000).WriteBytes([]byte("Hello world!")).WriteBytes([]byte("HAHAH")).Close()
	//Start()
}