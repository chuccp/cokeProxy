package entry

import (
	"bytes"
	"errors"
	"github.com/chuccp/utils/io"
	"github.com/chuccp/utils/log"
	"github.com/chuccp/utils/math"
	"github.com/chuccp/utils/queue"
	io2 "io"
	"strings"
)

type Stream struct {
	queue *queue.VQueue
	Id    uint32
}

func (stream *Stream) Write(p *Package) error {
	stream.queue.Offer(p)
	return nil
}
func (stream *Stream) Read() (*Package, error) {
	v, _ := stream.queue.Poll()
	p,ok:=v.(*Package)
	if ok{
		if p == nil {
			return nil, io2.EOF
		}
	}else{
		return nil, io2.EOF
	}
	return p, nil
}

func NewStream(Id uint32) *Stream {
	return &Stream{queue: queue.NewVQueue(),Id:Id}
}
func NewUrlStream(url string) *Stream {
	stream := NewStream(math.RandUInt32())
	stream.Write(HeaderUrlPackage(url,stream.Id))
	stream.Write(nil)
	return stream
}

type WebStream struct {
	stream      *Stream
	Length      uint32
	rLen uint32
	ContentType string
	data        *bytes.Buffer
	readHead    bool
}

func NewWebStream(stream *Stream) (*WebStream,error) {
	ws:= &WebStream{stream: stream, data: new(bytes.Buffer), readHead: false,rLen:0}
	_,err:=ws.header()
	if err!=nil{
		return nil, err
	}
	return ws,nil
}
func (webStream *WebStream) header() (n int, err error) {
	pack, err1 := webStream.stream.Read()
	log.Info("header11111",pack, err1)
	if err1 != nil {
		return 0, err1
	}
	if pack.PType == HeaderType {
		webStream.Length = pack.Total
		var buff bytes.Buffer
		pack.Data(&buff)
		reader := io.NewReadStream(&buff)
		data, err := reader.ReadLine()
		if err == nil {
			kv := strings.Split(string(data), ":")
			if kv[0] == "ContentType" {
				ContentType := kv[1]
				webStream.ContentType = ContentType
			}
		}
		return int(webStream.Length),nil
	} else {
		return 0, errors.New("FORMAT ERROR")
	}
}
func (webStream *WebStream) Read(p []byte) (int,  error) {
	log.Info("header============")
	if webStream.data.Len() == 0 {
		if webStream.rLen==webStream.Length{
			return 0,io2.EOF
		}
		pack, err3 := webStream.stream.Read()
		if err3 != nil {
			return 0, err3
		}
		webStream.data.Write(pack.DATA)
		webStream.rLen = webStream.rLen+pack.Len
		return webStream.data.Read(p)
	} else {
		return webStream.data.Read(p)
	}
}
