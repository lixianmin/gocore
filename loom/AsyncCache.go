/********************************************************************
created:    2018-12-02
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

package loom

import (
	"sync"
	"time"
)

type AsyncCache struct {
	m      sync.Map
	locker sync.Mutex
}

func NewAsyncCache(goPoolSize int, reloadDelay time.Duration, deleteDelay time.Duration) *AsyncCache {
	var cache = &AsyncCache{}
	go goAsyncCacheLoop(cache.m, goPoolSize, reloadDelay, deleteDelay)
	return cache
}

func goAsyncCacheLoop(m sync.Map, goPoolSize int, reloadDelay time.Duration, deleteDelay time.Duration) {
	var reloadTicker = time.NewTicker(reloadDelay / 8)
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
				if time.Now().After(item.getLoadTime().Add(reloadDelay)) {
					pool.Schedule(func() {
						item.loadData()
					})
				}

				return true
			})
		case <-deleteTicker.C:
			m.Range(func(key, value interface{}) bool {
				var item = value.(*cacheItem)
				if time.Now().After(item.getFetchTime().Add(deleteDelay)) {
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
		var locker = cache.locker
		locker.Lock()
		defer locker.Unlock()

		item, ok = cache.m.Load(key)
		if !ok {
			var newItem = newCacheItem(loader)
			newItem.loadData()
			cache.m.Store(key, newItem)
			item = newItem
		}
	}

	var lastItem = item.(*cacheItem)
	var data = lastItem.fetchData()
	return data
}
