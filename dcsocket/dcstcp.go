package dcsocket

import (
	"errors"
	"io"
	"net"
	"smsonline/devcommon"
)

type DCSTcp struct {
	conn    net.Conn
	running bool
	buffer  []byte
}

func NewDCSTcp(conn net.Conn) *DCSTcp {
	ptr := &DCSTcp{conn: conn}
	ptr.running = false
	ptr.buffer = make([]byte, 32768)
	ptr.Start()
	return ptr
}
func (that *DCSTcp) Start() {
	if !that.running {
		that.running = true
	}

}

func (that *DCSTcp) ReadByte() byte {
	return that.ReadBytes(1)[0]
}
func (that *DCSTcp) ReadBytes(length uint) []byte {
	var (
		result = make([]byte, 0)
		read   int
		err    error
	)
	for that.running && length > 0 {
		if read, err = that.conn.Read(that.buffer[:devcommon.Min(len(that.buffer), int(length))]); err != nil {
			panic(err)
		} else if read >= 0 {
			length -= uint(read)
			result = append(result, that.buffer[:read]...)
		} else {
			panic(errors.New("read<0"))
		}
	}
	return result
}
func (that *DCSTcp) WriteByte(data byte) {
	that.WriteBytes([]byte{data})
}
func (that *DCSTcp) WriteBytes(data []byte) *DCSTcp {
	if data == nil || len(data) <= 0 {
		return that
	}
	if write, err := that.conn.Write(data); err != nil {
		panic(err)
	} else if write != len(data) {
		panic(errors.New("write != len(data)"))
	}
	return that
}
func (that *DCSTcp) ReadTo(writer io.Writer, length uint) *DCSTcp {
	var (
		read  int
		write int
		err   error
	)
	for that.running && length > 0 {
		if read, err = that.conn.Read(that.buffer[:devcommon.Min(len(that.buffer), int(length))]); err != nil {
			panic(err)
		} else if read >= 0 {
			length -= uint(read)
			if write, err = writer.Write(that.buffer[:read]); err != nil {
				panic(err)
			} else if write != read {
				panic(errors.New("write != read"))
			}
		} else {
			panic(errors.New("read<0"))
		}
	}
	return that
}
func (that *DCSTcp) WriteFrom(reader io.Reader, length uint) *DCSTcp {
	var (
		read    int
		write   int
		err     error
	)
	for that.running && length > 0 {
		if read, err = reader.Read(that.buffer[:devcommon.Min(len(that.buffer), int(length))]); err != nil {
			panic(err)
		} else if read >= 0 {
			length -= uint(read)
			if write, err = that.conn.Write(that.buffer[:read]); err != nil {
				panic(err)
			} else if write != read {
				panic(errors.New("write != read"))
			}
		} else {
			panic(errors.New("read<0"))
		}
	}
	return that
}

func (that *DCSTcp) Stop() {
	that.running = false
}
