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
	accessTime UnixTime
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

func (item *cacheItem) getAccessTime() UnixTime {
	return UnixTime(atomic.LoadInt64((*int64)(&item.accessTime)))
}

func (item *cacheItem) isExpired(expireTime UnixTime) bool {
	return UnixTime(time.Now().Unix()) >= item.getAccessTime()+expireTime
}

func (item *cacheItem) setData(data interface{}) {
	item.data.Store(data)
}

func (item *cacheItem) getData() interface{} {
	return item.data.Load()
}
