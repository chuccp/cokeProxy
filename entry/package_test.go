package entry

import (
	"strconv"
	"strings"
	"testing"
)

func TestPackage_Bytes(t *testing.T) {

	var id uint64 = 0

	println(id<<0)

	var address = "127.0.0.1:2021"
	var addressPort = strings.Split(address,":")
	var host = addressPort[0]
	var port = addressPort[1]
	var hosts = strings.Split(host,".")

	for _, b := range hosts {
		num,_:=strconv.Atoi(b)
		id = id<<8|uint64(num)
	}
	println(id)
	p,_:=strconv.Atoi(port)
	id = (uint64)(p)|(id<<16)

	println("!!!",id)


}
