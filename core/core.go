package core

import "github.com/chuccp/cokeProxy/user"

type Context struct {
	storeMap *user.StoreMap
}

func NewContext() *Context {
	return &Context{storeMap:user.NewStoreMap()}
}
func (ctx *Context) GetUser(name string) user.IUser {
	ui,flag:=ctx.storeMap.Get(name)
	if flag{
		return ui
	}
	return nil
}
func (ctx *Context) AddUser(name string,user user.IUser)  {
	ctx.storeMap.Add(name,user)
}


