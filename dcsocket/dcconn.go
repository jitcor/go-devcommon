package dcsocket

import (
	"encoding/hex"
	"fmt"
	"math"
	"net"
	"runtime/debug"
	"time"
)

type DCConn struct {
	base net.Conn
}

func NewDCConn(conn net.Conn)*DCConn  {
	return &DCConn{base:conn}
}
func (that *DCConn)Read(b []byte) (n int, err error)  {
	n,err= that.base.Read(b)
	fmt.Println("ReadHex:")
	debug.PrintStack()
	fmt.Println(hex.Dump(b[:int(math.Min(float64(20),float64(n)))]))
	return
}
func (that *DCConn)Write(b []byte) (n int, err error)  {
	return that.base.Write(b)
}
func (that *DCConn)Close() error  {
	return that.base.Close()
}
func (that *DCConn)LocalAddr() net.Addr {
	return that.base.LocalAddr()
}
func (that *DCConn)RemoteAddr() net.Addr{
	return that.base.RemoteAddr()
}
func (that *DCConn)SetDeadline(t time.Time) error{
	return that.base.SetDeadline(t)
}
func (that *DCConn)SetReadDeadline(t time.Time) error{
	return that.base.SetReadDeadline(t)
}
func (that *DCConn)SetWriteDeadline(t time.Time) error{
	return that.base.SetWriteDeadline(t)
}
