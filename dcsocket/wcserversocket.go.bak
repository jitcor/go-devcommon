package dcsocket

import "syscall"

type WCServerSocket struct {
	ip string
	port int
	fd syscall.Handle
}

func NewWCServerSocket(ip string, port int) *WCServerSocket {
	ptr:=&WCServerSocket{ip:ip,port:port}
	return ptr
}
func (that *WCServerSocket)Accept()*WCSocket  {
	return nil
}
