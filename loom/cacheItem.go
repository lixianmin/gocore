/********************************************************************
created:    2018-12-02
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

package loom

import (
	"time"
)

type cacheItem struct {
	data      interface{}
	loader    func() interface{}
	loadTime  time.Time
	fetchTime time.Time
}

func newCacheItem(loader func() interface{}) *cacheItem {
	var item = &cacheItem{
		loader: loader,
	}

	return item
}

func (item *cacheItem) loadData() {
	item.data = item.loader()
	item.loadTime = time.Now()
}

func (item *cacheItem) fetchData() interface{} {
	item.fetchTime = time.Now()
	return item.data
}

func (item *cacheItem) getLoadTime() time.Time {
	return item.loadTime
}

func (item *cacheItem) getFetchTime() time.Time {
	return item.fetchTime
}
