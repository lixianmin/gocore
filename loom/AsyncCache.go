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
	expireTime UnixTime
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
		// 将超时时间改为10年。原来写的是100年，但后面定期delete使用的是4*expiration，会超过最大time.Duration的值
		expiration = 10 * 365 * 24 * time.Hour
	}

	var cache = &AsyncCache{
		m:          &sync.Map{},
		pool:       NewGoroutinePool(goPoolSize),
		expireTime: fromDuration(expiration),
	}

	go cache.goAsyncDelete(expiration)
	return cache
}

func (cache *AsyncCache) goAsyncDelete(expiration time.Duration) {
	defer DumpIfPanic("")
	var m = cache.m
	var deleteDelay = expiration * 4

	for {
		time.Sleep(deleteDelay)

		m.Range(func(key, value interface{}) bool {
			var item = value.(*cacheItem)
			var accessTime = item.getAccessTime()
			// 虽然加到了cache.m中，但是accessTime=0的因为从未使用过，所以不做处理
			if accessTime > 0 && fromTime(time.Now()) >= accessTime+cache.expireTime {
				m.Delete(key)
				atomic.AddInt32(&cache.length, -1)
			}

			return true
		})
	}
}

func (cache *AsyncCache) Get(key interface{}, loader func() (data interface{}, err error)) interface{} {
	if key == nil {
		return nil
	}

	// 到这里为止，相同的key对应拿到的item一定是相同的
	var item = cache.fetchCacheItem(key)
	var expireTime = cache.expireTime

	// 如果从未加载，或过期了，则需要先加载
	if item.isExpired(expireTime) {
		// 真正获得锁的那个协程，不一定是store一个item到cache.m中的那个，它可能是任意一个协程
		item.Lock()
		defer item.Unlock()

		// 如果从未加载过，或过期了，则加载并设置
		if item.isExpired(expireTime) {
			var data, err = loader()
			if err == nil {
				item.setData(data)
			}
		}

		item.setAccessTime()
		var data = item.getData()
		return data
	}

	// 时间过半，则发起一次异步加载
	if item.isExpired(expireTime>>1) && atomic.CompareAndSwapInt32(&item.Loading, 0, 1) {
		cache.pool.Schedule(func() {
			defer atomic.StoreInt32(&item.Loading, 0)
			var data, err = loader()
			if err == nil {
				item.setData(data)
			}
			item.setAccessTime()
		})
	}

	var data = item.getData()
	return data
}

func (cache *AsyncCache) Set(key interface{}, data interface{}) {
	if key != nil {
		var item = cache.fetchCacheItem(key)
		item.setData(data)
		item.setAccessTime()
	}
}

func (cache *AsyncCache) fetchCacheItem(key interface{}) *cacheItem {
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
	return theItem
}

func (cache *AsyncCache) GetCount() int {
	return int(atomic.LoadInt32(&cache.length))
}
