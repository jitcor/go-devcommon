package dcsocket

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

type SocketInfo struct {
	Conn      net.Conn
	IsStopped bool
}

func NewSocketInfo(conn net.Conn) *SocketInfo {
	return &SocketInfo{Conn: conn}
}
func GetServerListener(port int) (listener net.Listener, err error) {
	service := ":" + strconv.Itoa(port)
	if listener, err = net.Listen("tcp", service); err != nil {
		return nil, err
	} else {
		return listener, nil
	}
}
func GetTCPServerListener(port int) (listener *net.TCPListener, err error) {
	service := ":" + strconv.Itoa(port)
	if tcpAddr, err := net.ResolveTCPAddr("tcp4", service); err != nil {
		return nil, err
	} else if listener, err = net.ListenTCP("tcp", tcpAddr); err != nil {
		return nil, err
	} else {
		return listener, nil
	}
	//if tmp,err:=net.Listen("tcp",service);err!=nil{
	//	return
	//}else {
	//	listener=tmp
	//}
}

func GetUDPConn(port int) (udpConn *net.UDPConn, err error) {
	service := ":" + strconv.Itoa(port)
	if tcpAddr, err := net.ResolveUDPAddr("udp", service); err != nil {
		return nil, err
	} else if udpConn, err = net.ListenUDP("udp", tcpAddr); err != nil {
		return nil, err
	} else {
		return udpConn, nil
	}
}
func GetTCPConn(host string, port int) (conn net.Conn, err error) {
	return net.Dial("tcp", host+":"+strconv.Itoa(port))
}

func CheckPort(host string, port int) bool {
	timeout := time.Second
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, strconv.Itoa(port)), timeout)
	if err != nil {
		fmt.Println("Connecting error:", err)
		return false
	}
	if conn != nil {
		defer conn.Close()
		//fmt.Println("Opened", net.JoinHostPort(host, strconv.Itoa(port)))
		return true
	}
	return false
}
