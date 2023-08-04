package cacheServer

import (
	"memCache/cache"
	"time"
)

type cacheServer struct {
	cache cache.Cache //声明一个接口，用来做中间代理层，代理层能使得我们在中间处理一次输入
}

func NewCacheServer() *cacheServer {
	tmp := &cacheServer{
		cache: cache.NewCache(),
	}
	return tmp
}

//size : 1KB 100KB 1MB 2MB 1GB3
func (cs *cacheServer) SetMaxMemory(size string) bool {
	return cs.cache.SetMaxMemory(size)
}

//将value写入缓存
func (cs *cacheServer) set(key string, val interface{}, expire ...time.Duration) bool {
	//在这里处理了可变参数
	expireTime := time.Second
	if len(expire) > 0 {
		expireTime = expire[0]
	}
	return cs.cache.Set(key, val, expireTime)
}

//根据key值获取value
func (cs *cacheServer) Get(key string) (interface{}, bool) {
	return cs.cache.Get(key)
}

//删除key值
func (cs *cacheServer) Del(key string) bool {
	return cs.cache.Del(key)
}

//判断key是否存在
func (cs *cacheServer) Exists(key string) bool {
	return cs.cache.Exists(key)
}

//清空所有key
func (cs *cacheServer) Flush() bool {
	return cs.cache.Flush()
}

//获取缓存中所有key的数量
func (cs *cacheServer) Keys() int64 {
	return cs.cache.Keys()
}
