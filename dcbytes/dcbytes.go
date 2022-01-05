package dcbytes

import (
	"errors"
	"github.com/Humenger/go-devcommon"
	"io"
)

func BytesCopy( dst []byte, dstOffset int,src []byte, srcOffset int, length int) {
	minLength:=devcommon.Min(length,devcommon.Min(len(src)-srcOffset,len(dst)-dstOffset))
	if minLength>0&&dstOffset>=0&&srcOffset>=0{
		copy(dst[dstOffset:dstOffset+minLength],src[srcOffset:srcOffset+minLength])
	}

}
func ReaderCopy(reader io.Reader, out []byte) error {
	outBuf:=make([]byte,0)
	l:= len(out)
	buf:=make([]byte,devcommon.Min(l,1024))
	if n,err:=reader.Read(buf);err!=nil&&err!=io.EOF{
		return err
	}else {
		for n > 0 {
			outBuf=append(outBuf,buf[0:n]...)
			l-=n
			if l==0{
				break
			}
			if l<1024{
				buf=make([]byte,l)
			}
			if n,err=reader.Read(buf);err!=nil&&err!=io.EOF{
				return err
			}
		}
		if len(outBuf)> len(out){
			return errors.New("ReaderCopy is error")
		}
		BytesCopy(out,0,outBuf,0, len(outBuf))
		return nil
	}
}