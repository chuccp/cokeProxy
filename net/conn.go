package net

import (
	"bytes"
	"github.com/chuccp/cokeProxy/core"
	"github.com/chuccp/cokeProxy/entry"
	"github.com/chuccp/utils/io"
	"github.com/chuccp/utils/log"
	"github.com/chuccp/utils/math"
	io2 "io"
	"strings"
	"sync"
)

type conn struct {
	stream    *io.NetStream
	context *core.Context
	outStream *sync.Map
	name string
}

func (c *conn) Start() {
	go c.read()
}
func (c *conn)  GetName()string{
	return c.name
}
func (c *conn) read() {
	for {
		pack, err := c.readPackage()
		log.Info("!!!!!",pack, err)
		if err != nil {
			break
		}
		if pack.GetType()==entry.LoginType{
			log.Info("有连接来了")
			var buff  bytes.Buffer
			pack.Data(&buff)
			reader:=io.NewReadStream(&buff)
			data,err:=reader.ReadLine()
			if err==nil{
				kv:=strings.Split(string(data),":")
				if kv[0]=="name"{
					name := kv[1]
					c.name = name
					c.context.GetUserManage().Add(name,c)
				}
			}
		}else{
			v, ok := c.outStream.Load(pack.GetId())

			log.Info("--------",v,ok)

			if ok {
				out := v.(*entry.Stream)
				out.Write(pack)
			}
		}
	}
	c.context.GetUserManage().Remove(c.name)
	log.Info("连接断开")
}
func (c *conn) Write(stream *entry.Stream) (*entry.Stream, error) {
	for {
		age, err := stream.Read()
		log.Info("cocococo",age, err)
		if err == nil {
			var buffer  = new(bytes.Buffer)
			age.Bytes(buffer)
			data:=buffer.Bytes()
			log.Info("~~~~",string(data))
			c.stream.Write(data)
			c.stream.Flush()
		} else {
			if err == io2.EOF {
				log.Info("$$$$$$$$",age, err,stream.Id)
				st := entry.NewStream(stream.Id)
				c.outStream.Store(stream.Id, st)
				return st, nil
			} else {
				log.Info("**********",age, err)
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
		data,err:=c.stream.ReadBytes(4)
		if err != nil {
			return nil, err
		}
		pack.Len = math.U32BE(data[0:4])
		log.Info("len",pack.Len)
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
		log.Info("pack:",pack)
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
