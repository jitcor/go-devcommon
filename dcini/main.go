package dcini

import (
	"github.com/go-ini/ini"
	"log"
	"strconv"
)

type DcIni struct {
	name string
	section *ini.Section
}
var preference=getIni("preference")


func NewDcIni(name string,err *error) *DcIni {

	return &DcIni{name: name,section: preference.Section(name)}
}
func (that *DcIni) GetString(key,def string)string  {
	return that.section.Key(key).String()
}
func (that *DcIni) GetBool(key string,def bool)bool  {
	if val,e:= that.section.Key(key).Bool();e!=nil{
		return def
	}else {
		return val
	}
}
func (that *DcIni) GetInt(key string,def int32)int32  {
	if val,e:= that.section.Key(key).Int();e!=nil{
		return def
	}else {
		return int32(val)
	}
}
func (that *DcIni) GetLong(key string,def int64)int64  {
	if val,e:= that.section.Key(key).Int64();e!=nil{
		return def
	}else {
		return val
	}
}
func (that *DcIni) GetShort(key string,def int16)int16  {
	if val,e:= that.section.Key(key).Int();e!=nil{
		return def
	}else {
		return int16(val)
	}
}
func (that *DcIni) PutString(key,value string)*DcIni  {
	that.section.Key(key).SetValue(value)
	return that
}
func (that *DcIni) PutBool(key string,value bool)*DcIni  {
	that.section.Key(key).SetValue(strconv.FormatBool(value))
	return that
}
func (that *DcIni) PutInt(key string,value int32)*DcIni  {
	that.section.Key(key).SetValue(strconv.FormatInt(int64(value),10))
	return that
}


func getIni(name string) *ini.File {
	if preference!=nil{
		return preference
	}
	cfg, e := ini.Load(name+".ini")
	if e != nil {
		log.Println("Fail to read file: %w", e)
		return nil
	}
	return cfg
}