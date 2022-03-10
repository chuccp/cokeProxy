package web

import (
	"github.com/chuccp/cokeProxy/core"
	"github.com/chuccp/cokeProxy/entry"
	"github.com/chuccp/cokeProxy/user"
	"github.com/chuccp/utils/log"
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
		u, fa :=s.context.GetUserManage().Get(name)

		log.Info("name:",name,"  action:",action)

		if fa&&u!=nil{
			steam,err:=u.Write(entry.NewUrlStream(action))
			if err==nil{
				web,err2:=entry.NewWebStream(steam)
				if err2==nil{
					ctx.DataFromReader(http.StatusOK, int64(web.Length),web.ContentType,web,nil)
				}else{
					ctx.HTML(http.StatusOK,"haha",err2.Error())
				}

			}
		}
	})

	route.GET("/user/list", func(ctx *gin.Context) {
		values:=make([]string,0)
		s.context.GetUserManage().Each(func(user user.IUser) bool {
			values = append(values,user.GetName() )
			return true
		})
		ctx.JSON(http.StatusOK,values)
	})

	route.Run(":"+strconv.Itoa(s.port))
}

