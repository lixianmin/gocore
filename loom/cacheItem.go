/********************************************************************
created:    2018-12-02
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

package loom

import (
	"time"
	"sync"
	"github.com/lixianmin/gocore/safe"
)

type cacheItem struct {
	data      interface{}
	loader    func() interface{}
	loadTime  safe.Int64
	fetchTime safe.Int64
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
	item.loadTime.Store(time.Now().UnixNano())
}

func (item *cacheItem) fetchData() interface{} {
	item.fetchTime.Store(time.Now().UnixNano())
	return item.data
}

func (item *cacheItem) getLoadTime() int64 {
	return item.loadTime.Load()
}

func (item *cacheItem) getFetchTime() int64 {
	return item.fetchTime.Load()
}
