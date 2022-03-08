package net

import (
	"github.com/chuccp/cokeProxy/core"
	"github.com/chuccp/utils/io"
)

type Server struct {
	tcpServer *io.TCPServer
	context *core.Context
}

func NewServer(context *core.Context,port int)*Server  {
	return &Server{context:context,tcpServer:io.NewTCPServer(port)}
}
func (s *Server) Start()error  {
	err:=s.tcpServer.Bind()
	if err!=nil{
		return err
	}
	go s.accept()
	return nil
}
func (s *Server) accept() {
	for {
		io, err := s.tcpServer.Accept()
		if err!=nil{
			break
		}
		go newConn(io,s.context).Start()
	}
}
