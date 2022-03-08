package web

import (
	"github.com/chuccp/cokeProxy/core"
	"github.com/chuccp/cokeProxy/entry"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Server struct {
	context *core.Context
	port int
}

func NewServer(context *core.Context,port int)*Server  {
	return &Server{context:context,port: port}
}

func (s *Server) Start()  {
	route := gin.Default()
	//context:=core.NewContext()
	//
	//server:=net.NewServer(9080)
	//go server.Start()
	route.GET("/:name/*action", func(ctx *gin.Context) {
		name:=ctx.Param("name")
		action:=ctx.Param("action")
		u:=s.context.GetUser(name)
		if u!=nil{
			steam,err:=u.Write(entry.NewUrlStream(action))
			if err==nil{
				web:=entry.NewWebStream(steam)
				ctx.DataFromReader(http.StatusOK,web.Length,web.ContentType,web,nil)
			}
		}
	})
	route.Run(":"+strconv.Itoa(s.port))
}

