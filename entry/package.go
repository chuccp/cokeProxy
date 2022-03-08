package entry

import (
	"bytes"
	"github.com/chuccp/utils/math"
)

const (
	HeaderType byte = iota
	DataType
	LoginType
)

type ContentType string

const (
	FILE = ContentType("file")
	JSON = ContentType("json")
	HTML = ContentType("html")
)




type Package struct {
	PType byte
	Id    uint32
	Total uint32
	Len   uint32
	DATA  []byte
}

func (p *Package) GetId() uint32 {
	return p.Id
}
func (p *Package) HeaderBytes(buffer *bytes.Buffer) {
	var len = p.Len
	buffer.WriteByte(p.PType)
	buffer.Write(math.BEU32(p.Id))
	buffer.Write(math.BEU32(p.Total))
	buffer.Write(math.BEU32(p.Len))
	buffer.Write(p.DATA[0:len])
}
func (p *Package) DataBytes(buffer *bytes.Buffer) {
	var len = p.Len
	buffer.WriteByte(p.PType)
	buffer.Write(math.BEU32(p.Id))
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
func HeaderUrlPackage(url string) *Package {
	value := "url:" + url + "\n"
	data := []byte(value)
	id := math.RandInt()
	return &Package{PType: HeaderType, Id: id, Total: uint32(len(data)), Len: uint32(len(data)), DATA: data}
}
