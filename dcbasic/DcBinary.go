package dcbasic

import "encoding/binary"

type DcBinary struct {
	bigEndian bool
	data interface{}
}

func NewDcBinary() *DcBinary {
	return &DcBinary{
		bigEndian: false,
	}
}
func (that *DcBinary) BigEndian() *DcBinary {
	that.bigEndian=true
	return that
}
func (that *DcBinary) LittleEndian() *DcBinary {
	that.bigEndian=false
	return that
}
func (that *DcBinary) Put(data interface{}) *DcBinary {
	that.data=data
	return that
}
func (that *DcBinary) ToBytes() []byte {
	if that.data==nil{
		return nil
	}else if v,ok:=that.data.(int32);ok{
		result:=make([]byte,4)
		if that.bigEndian{
			binary.BigEndian.PutUint32(result,uint32(v))
			return result
		}
	}
	return nil
}
func (that *DcBinary) ToInt32() int32 {
	return 0
}
func (that *DcBinary) ToInt16() int16 {
	return 0
}
func (that *DcBinary) ToInt64() int64 {
	return 0
}
func (that *DcBinary) ToUint32() uint32 {
	return 0
}
func (that *DcBinary) ToUint16() uint16 {
	return 0
}
func (that *DcBinary) ToUint64() uint64 {
	return 0
}
