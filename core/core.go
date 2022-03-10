package core

import "github.com/chuccp/cokeProxy/user"

type Context struct {
	userManage *user.Manage
}

func NewContext() *Context {
	return &Context{userManage:user.NewMange()}
}
func (ctx *Context) GetUserManage() *user.Manage {
	return ctx.userManage
}


