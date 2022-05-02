package entry

import (
	"bytes"
	"errors"
	"github.com/chuccp/cokeProxy/net"
	"github.com/chuccp/utils/io"
	"github.com/chuccp/utils/queue"
	io2 "io"
	"net/http"
	"strings"
)

type Stream struct {
	queue *queue.LiteQueue
	Id    uint64
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

func NewStream(Id uint64) *Stream {
	return &Stream{queue: queue.NewLiteQueue(5),Id:Id}
}
func NewUrlStream(url string,r *http.Request) *Stream {
	u,_:=net.AddrToNum(r.RemoteAddr)
	stream := NewStream(u)
	pack:=HeaderUrlPackage(url,stream.Id)
	var Range = r.Header.Get("Range")
	if  Range!=""{
		pack.head["Range"] = Range
	}
	stream.Write(pack)
	return stream
}

type WebStream struct {
	stream      *Stream
	Length      uint32
	rLen uint32
	ContentType string
	data        *bytes.Buffer
	readHead    bool
	Head map[string]string
}

func NewWebStream(stream *Stream) (*WebStream,error) {
	ws:= &WebStream{stream: stream, data: new(bytes.Buffer), readHead: false,rLen:0,Head:make(map[string]string,0)}
	_,err:=ws.header()
	if err!=nil{
		return nil, err
	}
	return ws,nil
}
func (webStream *WebStream) header() (n int, err error) {
	pack, err1 := webStream.stream.Read()
	if err1 != nil {
		return 0, err1
	}
	if pack.PType == HeaderType {
		webStream.Length = pack.Total
		var buff bytes.Buffer
		pack.Data(&buff)
		reader := io.NewReadStream(&buff)
		for  {
			data, err4 := reader.ReadLine()
			if len(data)==0{
				break
			}
			if err4 == nil {
				kv := strings.Split(string(data), ":")
				if kv[0] == "ContentType" {
					ContentType := kv[1]
					webStream.ContentType = ContentType
				}else{
					webStream.Head[kv[0]] = kv[1]
				}
			}
		}

		return int(webStream.Length),nil
	} else {
		return 0, errors.New("FORMAT ERROR")
	}
}

func (webStream *WebStream) ReadBuffer(buffer *bytes.Buffer) error {
	buffer.Reset()
	if webStream.rLen==webStream.Length{
		return io2.EOF
	}
	pack, err3 := webStream.stream.Read()
	if err3 != nil {
		return  err3
	}
	buffer.Write(pack.DATA)
	webStream.rLen = webStream.rLen+pack.Len
	return nil
}

func (webStream *WebStream) Read(p []byte) (int,  error) {
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
