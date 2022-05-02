package net

import (
	"bytes"
	"encoding/binary"
	"github.com/chuccp/cokeProxy/core"
	"github.com/chuccp/cokeProxy/entry"
	"github.com/chuccp/utils/io"
	"github.com/chuccp/utils/log"
	"github.com/chuccp/utils/math"
	"strconv"
	"strings"
	"sync"
)

type conn struct {
	stream    *io.NetStream
	context   *core.Context
	outStream *sync.Map
	name      string
}

func (c *conn) Start() {
	go c.read()
}
func (c *conn) GetName() string {
	return c.name
}
func (c *conn) read() {
	for {
		pack, err := c.readPackage()
		if err != nil {
			break
		}
		if pack.GetType() == entry.LoginType {
			log.Info("有连接来了")
			var buff bytes.Buffer
			pack.Data(&buff)
			reader := io.NewReadStream(&buff)
			data, err := reader.ReadLine()
			if err == nil {
				kv := strings.Split(string(data), ":")
				if kv[0] == "name" {
					name := kv[1]
					c.name = name
					c.context.GetUserManage().Add(name, c)
				}
			}
		} else {
			v, ok := c.outStream.Load(pack.GetId())
			if ok {
				out := v.(*entry.Stream)
				out.Write(pack)
			}
		}
	}
	c.context.GetUserManage().Remove(c.name)
	log.Info("连接断开")
}

func (c *conn)Write(pack *entry.Package)error  {

	return nil
}

//func (c *conn) Write(stream *entry.Stream) (*entry.Stream, error) {
//	st := entry.NewStream(stream.Id)
//	c.outStream.Store(stream.Id, st)
//	go func() {
//		for {
//			age, err := stream.Read()
//			if err == nil {
//				var buffer = new(bytes.Buffer)
//				age.Bytes(buffer)
//				data := buffer.Bytes()
//				c.stream.Write(data)
//				c.stream.Flush()
//			} else {
//				break
//			}
//		}
//	}()
//	return st, nil
//}
func (c *conn) readPackage() (*entry.Package, error) {
	b, err := c.stream.ReadByte()
	if err != nil {
		return nil, err
	}
	var pack = &entry.Package{}
	if b == entry.LoginType {
		pack.PType = entry.LoginType
		data, err := c.stream.ReadBytes(4)
		if err != nil {
			return nil, err
		}
		pack.Len = math.U32BE(data[0:4])
		log.Info("len", pack.Len)
		data, err = c.stream.ReadUintBytes(pack.Len)
		if err != nil {
			return nil, err
		}
		pack.DATA = data
		return pack, nil

	} else if b == entry.HeaderType {
		pack.PType = entry.HeaderType
		data, err := c.stream.ReadBytes(16)
		if err != nil {
			return nil, err
		}
		pack.Id =binary.LittleEndian.Uint64(data[0:8])
		pack.Total = math.U32BE(data[8:12])
		pack.Len = math.U32BE(data[12:16])
		data, err = c.stream.ReadUintBytes(pack.Len)
		if err != nil {
			return nil, err
		}
		pack.DATA = data
		log.Info("pack:", pack)
		return pack, nil
	} else {
		pack.PType = entry.DataType
		data, err := c.stream.ReadBytes(12)
		if err != nil {
			return nil, err
		}
		pack.Id =binary.LittleEndian.Uint64(data[0:8])
		pack.Len = math.U32BE(data[8:12])
		data, err = c.stream.ReadUintBytes(pack.Len)
		if err != nil {
			return nil, err
		}
		pack.DATA = data
		return pack, nil
	}

}

func newConn(stream *io.NetStream, context *core.Context) *conn {
	return &conn{stream: stream, outStream: new(sync.Map), context: context}
}
func AddrToNum(address string) (uint64, error) {
	var id uint64 = 0
	var addressPort = strings.Split(address, ":")
	var host = addressPort[0]
	var port = addressPort[1]
	var hosts = strings.Split(host, ".")
	for _, b := range hosts {
		num, err := strconv.Atoi(b)
		if err != nil {
			return 0, err
		}
		id = id<<8 | uint64(num)
	}
	p, err := strconv.Atoi(port)
	if err != nil {
		return 0, err
	}
	id = (uint64)(p) | (id << 16)
	return id, nil
}
