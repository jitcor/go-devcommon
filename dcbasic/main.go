package dcbasic

import "encoding/binary"

//关于字节序，请参考https://www.ruanyifeng.com/blog/2016/11/byte-order.html
func Uint32ToBytesBigEndian(src uint32) []byte {
	ret := make([]byte, 4)
	binary.BigEndian.PutUint32(ret, src)
	return ret
}
func Int32ToBytesBigEndian(src int32) []byte {
	ret := make([]byte, 4)
	binary.BigEndian.PutUint32(ret, uint32(src))
	return ret
}
func Uint16ToBytesBigEndian(src uint16) []byte {
	ret := make([]byte, 2)
	binary.BigEndian.PutUint16(ret, src)
	return ret
}
func Int16ToBytesBigEndian(src int16) []byte {
	ret := make([]byte, 2)
	binary.BigEndian.PutUint16(ret, uint16(src))
	return ret
}
func Uint64ToBytesBigEndian(src uint64) []byte {
	ret := make([]byte, 8)
	binary.BigEndian.PutUint64(ret, src)
	return ret
}
func Int64ToBytesBigEndian(src int64) []byte {
	ret := make([]byte, 8)
	binary.BigEndian.PutUint64(ret, uint64(src))
	return ret
}

func Uint32ToBytesLittleEndian(src uint32) []byte {
	ret := make([]byte, 4)
	binary.LittleEndian.PutUint32(ret, src)
	return ret
}
func Int32ToBytesLittleEndian(src int32) []byte {
	ret := make([]byte, 4)
	binary.LittleEndian.PutUint32(ret, uint32(src))
	return ret
}
func Uint16ToBytesLittleEndian(src uint16) []byte {
	ret := make([]byte, 2)
	binary.LittleEndian.PutUint16(ret, src)
	return ret
}
func Int16ToBytesLittleEndian(src int16) []byte {
	ret := make([]byte, 2)
	binary.LittleEndian.PutUint16(ret, uint16(src))
	return ret
}
func Uint64ToBytesLittleEndian(src uint64) []byte {
	ret := make([]byte, 8)
	binary.LittleEndian.PutUint64(ret, src)
	return ret
}
func Int64ToBytesLittleEndian(src int64) []byte {
	ret := make([]byte, 8)
	binary.LittleEndian.PutUint64(ret, uint64(src))
	return ret
}
