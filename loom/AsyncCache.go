/********************************************************************
created:    2018-12-02
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

package loom

import (
	"sync"
	"time"
	"strconv"
)

type AsyncCache struct {
	m sync.Map
}

func NewAsyncCache(goPoolSize int, reloadDelay time.Duration, deleteDelay time.Duration) *AsyncCache {
	if goPoolSize <= 0 {
		var message = "Invalid goPoolSize= " + strconv.Itoa(goPoolSize)
		panic(message)
	}

	if reloadDelay <= 0 {
		var message = "Invalid reloadDelay= " + strconv.Itoa(int(reloadDelay))
		panic(message)
	}

	if deleteDelay <= 0 {
		var message = "Invalid deleteDelay= " + strconv.Itoa(int(deleteDelay))
		panic(message)
	}

	var cache = &AsyncCache{}
	go goAsyncCacheLoop(cache.m, goPoolSize, reloadDelay, deleteDelay)
	return cache
}

func goAsyncCacheLoop(m sync.Map, goPoolSize int, reloadDelay time.Duration, deleteDelay time.Duration) {
	var reloadTicker = time.NewTicker(reloadDelay / 10)
	var deleteTicker = time.NewTicker(deleteDelay)

	defer func() {
		DumpIfPanic()
		reloadTicker.Stop()
		deleteTicker.Stop()
	}()

	var pool = NewGoroutinePool(goPoolSize)

	for {
		select {
		case <-reloadTicker.C:
			m.Range(func(key, value interface{}) bool {
				var item = value.(*cacheItem)
				var loadTime = item.getLoadTime()
				// 虽然加到了cache.m中，但未loadTime=0的因为从未加载过，所以不做处理
				if loadTime > 0 && time.Now().UnixNano() >= loadTime+int64(reloadDelay) {
					pool.Schedule(func() {
						item.loadData()
					})
				}

				return true
			})
		case <-deleteTicker.C:
			m.Range(func(key, value interface{}) bool {
				var item = value.(*cacheItem)
				var fetchTime = item.getFetchTime()
				// 虽然加到了cache.m中，但未fetchTime=0的因为从未使用过，所以不做处理
				if fetchTime > 0 && time.Now().UnixNano() >= fetchTime+int64(deleteDelay) {
					m.Delete(key)
				}

				return true
			})
		}
	}
}

func (cache *AsyncCache) Get(key string, loader func() interface{}) interface{} {
	var item, ok = cache.m.Load(key)
	if !ok {
		item, _ = cache.m.LoadOrStore(key, newCacheItem(loader))
	}

	// 到这里为止，相同的key对应拿到的item一定是相同的
	var theItem = item.(*cacheItem)
	if theItem.getLoadTime() == 0 {
		// 真正获得锁的那个协程，不一定是store一个item到cache.m中的那一个，它可能是任意一个协程
		theItem.Lock()
		defer theItem.Unlock()

		// 如果从来未加载过，则加载并设置loadTime
		if theItem.getLoadTime() == 0 {
			theItem.loadData()
		}
	}

	var data = theItem.fetchData()
	return data
}
