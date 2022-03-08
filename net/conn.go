package net

import (
	"bytes"
	"github.com/chuccp/cokeProxy/core"
	"github.com/chuccp/cokeProxy/entry"
	"github.com/chuccp/utils/io"
	"github.com/chuccp/utils/math"
	io2 "io"
	"strings"
	"sync"
)

type conn struct {
	stream    *io.NetStream
	context *core.Context
	outStream *sync.Map
}

func (c *conn) Start() {
	go c.read()
}
func (c *conn) read() {
	for {
		pack, err := c.readPackage()
		if err != nil {
			break
		}
		if pack.GetType()==entry.LoginType{
			var buff  bytes.Buffer
			pack.Data(&buff)
			reader:=io.NewReadStream(&buff)
			data,err:=reader.ReadLine()
			if err==nil{
				kv:=strings.Split(string(data),":")
				if kv[0]=="name"{
					name := kv[1]
					c.context.AddUser(name,c)
				}
			}
			continue
		}else{
			v, ok := c.outStream.Load(pack.GetId())
			if ok {
				out := v.(*entry.Stream)
				out.Write(pack)
			}
		}
	}
}
func (c *conn) Write(stream *entry.Stream) (*entry.Stream, error) {
	for {
		age, err := stream.Read()
		if err == nil {
			var buffer bytes.Buffer
			age.Bytes(&buffer)
			c.stream.Write(buffer.Bytes())
		} else {
			if err == io2.EOF {
				st := entry.NewStream()
				c.outStream.Store(st.Id, st)
				return st, nil
			} else {
				return nil, err
			}
		}
	}
	return nil, nil
}
func (c *conn) readPackage() (*entry.Package, error) {
	b, err := c.stream.ReadByte()
	if err != nil {
		return nil, err
	}
	var pack = &entry.Package{}
	if b == entry.LoginType{
		pack.PType = entry.LoginType
		data,err:=c.stream.ReadBytes(8)
		if err != nil {
			return nil, err
		}
		pack.Len = math.U32BE(data[4:8])
		data, err = c.stream.ReadUintBytes(pack.Len)
		if err != nil {
			return nil, err
		}
		pack.DATA = data
		return pack, nil

	}else if b == entry.HeaderType {
		pack.PType = entry.HeaderType
		data,err:=c.stream.ReadBytes(12)
		if err != nil {
			return nil, err
		}
		pack.Id = math.U32BE(data[0:4])
		pack.Total = math.U32BE(data[4:8])
		pack.Len = math.U32BE(data[8:12])
		data, err = c.stream.ReadUintBytes(pack.Len)
		if err != nil {
			return nil, err
		}
		pack.DATA = data
		return pack, nil
	} else {
		pack.PType = entry.DataType
		data,err:=c.stream.ReadBytes(8)
		if err != nil {
			return nil, err
		}
		pack.Id = math.U32BE(data[0:4])
		pack.Len = math.U32BE(data[4:8])
		data, err = c.stream.ReadUintBytes(pack.Len)
		if err != nil {
			return nil, err
		}
		pack.DATA = data
		return pack, nil
	}

}

func newConn(stream *io.NetStream,context *core.Context) *conn {
	return &conn{stream: stream, outStream: new(sync.Map),context:context}
}
