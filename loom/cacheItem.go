/********************************************************************
created:    2018-12-02
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

package loom

import (
	"time"
	"sync"
	"sync/atomic"
)

type cacheItem struct {
	data      interface{}
	loader    func() interface{}
	loadTime  int64
	fetchTime int64
	sync.Mutex
}

func newCacheItem(loader func() interface{}) *cacheItem {
	var item = &cacheItem{
		loader: loader,
	}

	return item
}

func (item *cacheItem) loadData() {
	item.data = item.loader()
	atomic.StoreInt64(&item.loadTime, time.Now().UnixNano())
}

func (item *cacheItem) fetchData() interface{} {
	atomic.StoreInt64(&item.fetchTime, time.Now().UnixNano())
	return item.data
}

func (item *cacheItem) getLoadTime() int64 {
	return atomic.LoadInt64(&item.loadTime)
}

func (item *cacheItem) getFetchTime() int64 {
	return atomic.LoadInt64(&item.fetchTime)
}
