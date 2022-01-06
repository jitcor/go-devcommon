package dcsocket

import (
. "fmt"
"strconv"
"strings"
"syscall"
	"unsafe"
)


type WCSocket struct {
	ip string
	port int
	sock syscall.Handle
}


func NewWCSocket(ip string, port int) *WCSocket {
	ptr:=&WCSocket{ip:ip,port:port}
	ptr.init()
	return ptr
}

func (that *WCSocket) init() {
	var (
		addr    syscall.SockaddrInet4
		wsadata syscall.WSAData
		err     error
	)
	if err = syscall.WSAStartup(MAKEWORD(2, 2), &wsadata); err != nil {
		Println("Startup error")
		return
	}
	//defer syscall.WSACleanup()

	if that.sock, err = syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, syscall.IPPROTO_IP); err != nil {
		Println("Socket create error")
		return
	}
	//defer syscall.Closesocket(that.sock)

	addr.Addr = inet_addr(that.ip)
	addr.Port = that.port
	if err = syscall.Connect(that.sock, &addr); err != nil {
		Println("Connect error")
		return
	}

}
func (that *WCSocket)WriteBytes(bytes []byte) *WCSocket {
	var (
		data       syscall.WSABuf
		SendButes  uint32
		overlapped syscall.Overlapped
	)
	data.Len = uint32(len(bytes))
	data.Buf = (*byte)(unsafe.Pointer(&bytes[0]))
	//如果使用syscall.Sendto或syscall.Write会发送失败，原因未知
	err:= syscall.WSASend(that.sock, &data, 1, &SendButes, 0, &overlapped, nil)
	if err!=nil{
		panic(err)
	}
	//_=syscall.FlushFileBuffers(that.sock)
	return that
}
func (that *WCSocket)Close()  {
	_ = syscall.WSACleanup()
	_=syscall.Closesocket(that.sock)
}
func MAKEWORD(low, high uint8) uint32 {
	var ret uint16 = uint16(high)<<8 + uint16(low)
	return uint32(ret)
}

func inet_addr(ipaddr string) [4]byte {
	var (
		ips = strings.Split(ipaddr, ".")
		ip  [4]uint64
		ret [4]byte
	)
	for i := 0; i < 4; i++ {
		ip[i], _ = strconv.ParseUint(ips[i], 10, 8)
	}
	for i := 0; i < 4; i++ {
		ret[i] = byte(ip[i])
	}
	return ret
}

func Start() {
	var (
		sock    syscall.Handle
		addr    syscall.SockaddrInet4
		wsadata syscall.WSAData
		err     error
	)
	if err = syscall.WSAStartup(MAKEWORD(2, 2), &wsadata); err != nil {
		Println("Startup error")
		return
	}
	defer syscall.WSACleanup()

	if sock, err = syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, syscall.IPPROTO_IP); err != nil {
		Println("Socket create error")
		return
	}
	defer syscall.Closesocket(sock)

	addr.Addr = inet_addr("127.0.0.1")
	addr.Port = 8000
	if err = syscall.Connect(sock, &addr); err != nil {
		Println("Connect error")
		return
	}

	var (
		data       syscall.WSABuf
		sendstr    string = "hello"
		SendButes  uint32
		overlapped syscall.Overlapped
	)
	data.Len = uint32(len(sendstr))
	data.Buf = syscall.StringBytePtr(sendstr)
	//如果使用syscall.Sendto或syscall.Write会发送失败，原因未知
	err = syscall.WSASend(sock, &data, 1, &SendButes, 0, &overlapped, nil)
	if err != nil {
		Println("Send error")
	} else {
		Println("Send success")
	}
}
