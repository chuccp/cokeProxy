package web

import (
	"github.com/chuccp/cokeProxy/core"
	"github.com/chuccp/cokeProxy/user"
	"github.com/chuccp/utils/log"
	"github.com/gin-gonic/gin"
	io2 "io"
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
	route.GET("/:name/*action", func(ctx *gin.Context) {
		name:=ctx.Param("name")
		action:=ctx.Param("action")
		u, fa :=s.context.GetUserManage().Get(name)
		log.Info("==================",ctx.Request.RemoteAddr,"   ","name:",name,"  action:",action)
		if fa&&u!=nil{
			//urlStream:=entry.NewUrlStream(action,ctx.Request)
			//steam,err:=u.Write(urlStream)
			//if err==nil{
			//	web,err2:=entry.NewWebStream(steam)
			//	if err2==nil{
			//		ctx.Header("Content-Length",strconv.FormatUint(uint64(web.Length),10))
			//		ctx.Header("Content-Type",web.ContentType)
			//		for k, v := range web.Head {
			//			ctx.Header(k,v)
			//		}
			//		var buffer  = new(bytes.Buffer)
			//		ctx.Stream(func(w io2.Writer) bool {
			//			err4:=web.ReadBuffer(buffer)
			//			if err4!=nil{
			//				return false
			//			}
			//			_, err5 := w.Write(buffer.Bytes())
			//			return err5==nil
			//		})
			//	}else{
			//		ctx.HTML(http.StatusOK,"haha",err2.Error())
			//	}
			//
			//}
		}
	})

	route.GET("/test", func(ctx *gin.Context) {

		num := 5
		ctx.Stream(func(w io2.Writer) bool {

			num--
			if num==0{
				return false
			}
			w.Write([]byte("=============="))

			return true
		})
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

