package dcmemory

import "github.com/Humenger/go-devcommon"

func MemCpy(dst ,src*[]byte,length int)*[]byte {
	fLen:=devcommon.Min(len(*dst),devcommon.Min(length,len(*src)))
	for i := 0; i < fLen; i++ {
		(*dst)[i]=(*src)[i]
	}
	return dst
}
func MemCmp(p1, p2 *[]byte, length int) int {
	if *p1==nil||*p2==nil{
		return -1
	}
	if len(*p1)<length{
		return -1
	}
	if len(*p2)<length{
		return 1
	}
	for i := 0; i < length; i++ {
		if (*p1)[i]!=(*p2)[i]{
			return int((*p1)[i] - (*p2)[i])
		}
	}
	return 0
}
