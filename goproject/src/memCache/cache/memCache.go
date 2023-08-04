package cache

import (
	"fmt"
	"sync"
	"time"
)

type memCache struct {
	//最大内存
	MaxSize uint64
	//最大内存字符串表示
	MaxSizeStr string
	//当前使用
	CurrentSize uint64
	//内存中装有的值
	values map[string]*memCacheValues
	//内存锁，读写互斥锁
	locker sync.RWMutex
	//清除过期内存的时间间隔
	ClearTimeInterval time.Duration
}

type memCacheValues struct {
	//值
	value interface{}
	//过期时间
	deadTime time.Time
	//过期跨度
	duraTime time.Duration
	//占用大小
	ocSize uint64
}

func NewCache() Cache {
	newCache := &memCache{
		values:            make(map[string]*memCacheValues), //指针类型的必须make！！
		ClearTimeInterval: time.Second,
	}
	go newCache.timedClear()
	return newCache
}

//size : 1KB 100KB 1MB 2MB 1GB3
func (mc *memCache) SetMaxMemory(size string) bool {
	mc.MaxSize, mc.MaxSizeStr = parseSize(size)
	fmt.Println(mc.MaxSize, mc.MaxSizeStr)
	return true
}

//将value写入缓存
func (mc *memCache) Set(key string, val interface{}, expire time.Duration) bool {
	mc.locker.Lock()
	defer mc.locker.Unlock()
	tmpValue := &memCacheValues{
		value:    val,
		deadTime: time.Now().Add(expire),
		duraTime: expire,
		ocSize:   ParseOcSize(val),
	}
	tmp, ok := mc.values[key] //首先看看是不是更改
	if ok && tmp != nil {     //是更改，那么需要删除原来的数据，再增添新的数据
		mc.del(key, tmp)
	}
	mc.add(key, tmpValue)
	if mc.CurrentSize > mc.MaxSize {
		mc.del(key, tmpValue)
		panic(fmt.Sprintf("max memory size %s", mc.MaxSizeStr))
	}
	return true
}

func (mc *memCache) add(key string, val *memCacheValues) {
	mc.CurrentSize += val.ocSize
	mc.values[key] = val
}

func (mc *memCache) del(key string, val *memCacheValues) {
	tmp, ok := mc.values[key]
	if ok && tmp != nil {
		delete(mc.values, key)
		mc.CurrentSize -= val.ocSize
	}
}

//根据key值获取value
func (mc *memCache) Get(key string) (interface{}, bool) {
	mc.locker.RLock()
	defer mc.locker.RUnlock()
	tmp, ok := mc.values[key]
	if !ok && tmp == nil {
		return nil, false
	} else if tmp.duraTime != 0 && tmp.deadTime.Before(time.Now()) {
		return tmp.value, true
	} else {
		mc.del(key, tmp)
	}
	return nil, false
}

//删除key值
func (mc *memCache) Del(key string) bool {
	mc.locker.Lock()
	defer mc.locker.Unlock()
	tmp, ok := mc.values[key]
	if ok && tmp != nil {
		mc.del(key, tmp)
		return true
	}
	return false
}

//判断key是否存在
func (mc *memCache) Exists(key string) bool {
	mc.locker.RLock()
	defer mc.locker.RUnlock()
	_, ok := mc.values[key]
	if ok {
		return true
	}
	return false
}

//清空所有key
func (mc *memCache) Flush() bool {
	mc.locker.Lock()
	defer mc.locker.Unlock()
	mc.values = make(map[string]*memCacheValues)
	mc.CurrentSize = 0
	return true
}

//获取缓存中所有key的数量
func (mc *memCache) Keys() int64 {
	mc.locker.RLock()
	defer mc.locker.RUnlock()
	num := int64(len(mc.values))
	return num
}

//协程函数，定期清除过期的内存
func (mc *memCache) timedClear() {
	timeInterval := time.NewTicker(mc.ClearTimeInterval)
	defer timeInterval.Stop()
	for {
		select {
		case <-timeInterval.C:
			for key, value := range mc.values {
				if value.duraTime == 0 || value.deadTime.After(time.Now()) {
					mc.locker.Lock()
					mc.del(key, value)
					mc.locker.Unlock()
				}
			}
		}
	}
}
