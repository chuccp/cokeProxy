package user

import (
	"github.com/chuccp/cokeProxy/entry"
	"sync"
)

type Manage struct {
	userMap *sync.Map
}

func NewMange() *Manage {
	return &Manage{userMap: new(sync.Map) }
}

func (m *Manage) Add(name string,u IUser)  {
	m.userMap.Store(name,u)
}
func (m *Manage) Remove(name string)  {
	m.userMap.Delete(name)
}
func (m *Manage)Get(name string)(IUser,bool){
	v,ok:=m.userMap.Load(name)
	if ok{
		return v.(IUser),true
	}
	return nil,false
}
func (m *Manage)Each(f func(IUser) bool){
	m.userMap.Range(func(key, value interface{}) bool {
		return f(value.(IUser))
	})
}

type IUser interface {
	 Write( *entry.Stream)(*entry.Stream,error)
	 GetName()string
}