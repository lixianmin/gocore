/********************************************************************
created:    2018-12-02
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

package loom

import (
	"sync"
	"sync/atomic"
	"time"
)

type cacheItem struct {
	data       atomic.Value
	accessTime int64
	Loading    int32
	sync.Mutex
}

func newCacheItem() *cacheItem {
	var item = &cacheItem{}
	return item
}

func (item *cacheItem) setAccessTime() {
	atomic.StoreInt64(&item.accessTime, time.Now().UnixNano())
}

func (item *cacheItem) getAccessTime() int64 {
	return atomic.LoadInt64(&item.accessTime)
}

func (item *cacheItem) isExpired(expiration time.Duration) bool {
	return time.Now().UnixNano() >= item.getAccessTime()+int64(expiration)
}

func (item *cacheItem) setData(data interface{}) {
	item.data.Store(data)
}

func (item *cacheItem) getData() interface{} {
	return item.data.Load()
}
