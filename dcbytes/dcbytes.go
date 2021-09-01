package dcbytes

import "github.com/Humenger/go-devcommon"

func BytesCopy( dst []byte, dstOffset int,src []byte, srcOffset int, length int) {
	minLength:=devcommon.Min(length,devcommon.Min(len(src)-srcOffset,len(dst)-dstOffset))
	if minLength>0&&dstOffset>=0&&srcOffset>=0{
		copy(dst[dstOffset:dstOffset+minLength],src[srcOffset:srcOffset+minLength])
	}

}
