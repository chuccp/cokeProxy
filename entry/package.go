package entry

import (
	"bytes"
	"encoding/binary"
	"github.com/chuccp/utils/math"
)

const (
	HeaderType byte = iota
	DataType
	LoginType
	CmdType
)

type ContentType string

type Package struct {
	PType byte
	Id    uint64
	Total uint32
	Len   uint32
	DATA  []byte

	head map[string]string
}

func (p *Package) GetId() uint64 {
	return p.Id
}
func (p *Package) HeaderBytes(buffer *bytes.Buffer) {
	buffer.WriteByte(p.PType)
	var data  = []byte{0,0,0,0,0,0,0,0}
	binary.LittleEndian.PutUint64(data,p.Id)
	buffer.Write(data)
	buffer.Write(math.BEU32(p.Total))
	var buff  = new(bytes.Buffer)
	for key, value := range p.head {
		buff.WriteString(key+":"+value+"\n")
	}
	buffer.Write(math.BEU32(uint32(buff.Len())))
	buffer.Write(buff.Bytes())
}
func (p *Package) DataBytes(buffer *bytes.Buffer) {
	var len = p.Len
	buffer.WriteByte(p.PType)
	var data  = []byte{0,0,0,0,0,0,0,0}
	binary.LittleEndian.PutUint64(data,p.Id)
	buffer.Write(data)
	buffer.Write(math.BEU32(p.Len))
	buffer.Write(p.DATA[0:len])
}
func (p *Package) LoginBytes(buffer *bytes.Buffer) {
	var len = p.Len
	buffer.WriteByte(p.PType)
	buffer.Write(math.BEU32(p.Len))
	buffer.Write(p.DATA[0:len])
}
func (p *Package) Data(buffer *bytes.Buffer) {
	buffer.Write(p.DATA[0:p.Len])
}
func (p *Package) GetType() byte {
	return p.PType
}

func (p *Package) Bytes(b *bytes.Buffer) {
	if p.PType==HeaderType{
		p.HeaderBytes(b)
	}else if p.PType==LoginType{
		p.LoginBytes(b)
	}else{
		p.DataBytes(b)
	}
}
func HeaderUrlPackage(url string,id uint64) *Package {
	pack:= &Package{PType: HeaderType, Id: id, Total: 0,head: make(map[string]string)}
	pack.head["url"] = url
	return pack
}
