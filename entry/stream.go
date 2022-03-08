package entry

import (
	"bytes"
	"errors"
	"github.com/chuccp/utils/io"
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
	if v == nil {
		return nil, io2.EOF
	}
	return v.(*Package), nil
}

func NewStream() *Stream {
	return &Stream{queue: queue.NewVQueue()}
}
func NewUrlStream(url string) *Stream {
	stream := NewStream()
	stream.Write(HeaderUrlPackage(url))
	stream.Write(nil)
	return stream
}

type WebStream struct {
	stream      *Stream
	Length      int64
	ContentType string
	data        *bytes.Buffer
	readHead    bool
}

func NewWebStream(stream *Stream) *WebStream {
	return &WebStream{stream: stream, data: new(bytes.Buffer), readHead: false}
}

func (webStream *WebStream) Read(p []byte) (n int, err error) {
	if !webStream.readHead {
		webStream.readHead = true
		pack, err := webStream.stream.Read()
		if err != nil {
			return 0, err
		}
		if pack.PType == HeaderType {
			webStream.Length = int64(pack.Total)
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
		} else {
			return 0, errors.New("FORMAT ERROR")
		}
	}
	if webStream.data.Len() == 0 {
		pack, err := webStream.stream.Read()
		if err != nil {
			return 0, err
		}
		webStream.data.Write(pack.DATA)
		return webStream.data.Read(p)
	} else {
		return webStream.data.Read(p)
	}
}
