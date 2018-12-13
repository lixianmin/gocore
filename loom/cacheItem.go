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
	accessTime unixTime
	Loading    int32
	sync.Mutex
}

func newCacheItem() *cacheItem {
	var item = &cacheItem{}
	return item
}

func (item *cacheItem) setAccessTime() {
	atomic.StoreInt64((*int64)(&item.accessTime), time.Now().Unix())
}

func (item *cacheItem) getAccessTime() unixTime {
	return unixTime(atomic.LoadInt64((*int64)(&item.accessTime)))
}

func (item *cacheItem) isExpired(expireTime unixTime) bool {
	return unixTime(time.Now().Unix()) >= item.getAccessTime()+expireTime
}

func (item *cacheItem) setData(data interface{}) {
	item.data.Store(data)
}

func (item *cacheItem) getData() interface{} {
	return item.data.Load()
}
