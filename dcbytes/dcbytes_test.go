package dcbytes

import (
	"log"
	"testing"
)

func TestBytesCopy(t *testing.T) {
	dst:=[]byte{0,1,2,3,4,5,6,7,8,9}
	src:=[]byte{10,11,12,13,14,15,16}
	BytesCopy(dst,0,src,5,4)
	log.Println("dst:",dst)
}
