/********************************************************************
created:    2018-12-02
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

package loom

import (
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

type AsyncCache struct {
	m          *sync.Map
	pool       *GoroutinePool
	expiration time.Duration
	length     int32
}

// expiration为0则代表永不超时
func NewAsyncCache(goPoolSize int, expiration time.Duration) *AsyncCache {
	if goPoolSize <= 0 {
		var message = "Invalid goPoolSize= " + strconv.Itoa(goPoolSize)
		panic(message)
	}

	if expiration < 0 || int64(expiration) >= time.Now().UnixNano() {
		var message = "Invalid expiration= " + strconv.Itoa(int(expiration))
		panic(message)
	}

	var isNeverExpire = expiration == 0
	if isNeverExpire {
		// 将超时时间改为100年
		expiration = 100 * 365 * 24 * time.Hour
	}

	var cache = &AsyncCache{
		m:          &sync.Map{},
		pool:       NewGoroutinePool(goPoolSize),
		expiration: expiration,
	}

	go cache.goAsyncDelete()
	return cache
}

func (cache *AsyncCache) goAsyncDelete() {
	defer DumpIfPanic()
	var m = cache.m
	var deleteDelay = cache.expiration * 4

	for {
		time.Sleep(deleteDelay)

		m.Range(func(key, value interface{}) bool {
			var item = value.(*cacheItem)
			var accessTime = item.getAccessTime()
			// 虽然加到了cache.m中，但是accessTime=0的因为从未使用过，所以不做处理
			if accessTime > 0 && time.Now().UnixNano() >= accessTime+int64(deleteDelay) {
				m.Delete(key)
				atomic.AddInt32(&cache.length, -1)
			}

			return true
		})
	}
}

func (cache *AsyncCache) Get(key interface{}, loader func() (interface{}, error)) interface{} {
	var item, ok = cache.m.Load(key)
	if !ok {
		var loaded bool
		item, loaded = cache.m.LoadOrStore(key, newCacheItem())
		if !loaded {
			atomic.AddInt32(&cache.length, 1)
		}
	}

	// 到这里为止，相同的key对应拿到的item一定是相同的
	var theItem = item.(*cacheItem)
	var expiration = cache.expiration

	// 如果从未加载，或过期了，则需要先加载
	if theItem.isExpired(expiration) {
		// 真正获得锁的那个协程，不一定是store一个item到cache.m中的那个，它可能是任意一个协程
		theItem.Lock()
		defer theItem.Unlock()

		// 如果从未加载过，或过期了，则加载并设置
		if theItem.isExpired(expiration) {
			var data, err = loader()
			if err == nil {
				theItem.setData(data)
			}
		}

		theItem.setAccessTime()
		var data = theItem.getData()
		return data
	}

	// 时间过半，则发起一次异步加载
	if theItem.isExpired(expiration>>1) && atomic.CompareAndSwapInt32(&theItem.Loading, 0, 1) {
		cache.pool.Schedule(func() {
			defer atomic.StoreInt32(&theItem.Loading, 0)
			var data, err = loader()
			if err == nil {
				theItem.setData(data)
			}
			theItem.setAccessTime()
		})
	}

	var data = theItem.getData()
	return data
}

func (cache *AsyncCache) GetCount() int {
	return int(atomic.LoadInt32(&cache.length))
}
