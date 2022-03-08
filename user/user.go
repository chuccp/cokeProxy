package user

import (
	"github.com/chuccp/cokeProxy/entry"
	"sync"
)

type StoreMap struct {
	userMap *sync.Map
}

func NewStoreMap() *StoreMap {
	return &StoreMap{userMap:new(sync.Map) }
}

func (m *StoreMap) Add(name string,u IUser)  {
	m.userMap.Store(name,u)
}
func (m *StoreMap)Get(name string)(IUser,bool){
	v,ok:=m.userMap.Load(name)
	if ok{
		return v.(IUser),true
	}
	return nil,false
}

type IUser interface {
	 Write( *entry.Stream)(*entry.Stream,error)
}