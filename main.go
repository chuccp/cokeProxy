package main

import (
	"github.com/chuccp/cokeProxy/core"
	"github.com/chuccp/cokeProxy/net"
	"github.com/chuccp/cokeProxy/web"
)

func main() {

	context:=core.NewContext()

	nnn:=net.NewServer(context,8091)

	go nnn.Start()

	server:=web.NewServer(context,8090)
	server.Start()

}
