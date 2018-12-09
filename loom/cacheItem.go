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
	data       interface{}
	accessTime int64
	loading    int32
	sync.Mutex
}

func newCacheItem() *cacheItem {
	var item = &cacheItem{}
	return item
}

func (item *cacheItem) setData(data interface{}) {
	item.data = data
	atomic.StoreInt64(&item.accessTime, time.Now().UnixNano())
}

func (item *cacheItem) getData() interface{} {
	atomic.StoreInt64(&item.accessTime, time.Now().UnixNano())
	return item.data
}

func (item *cacheItem) getAccessTime() int64 {
	return atomic.LoadInt64(&item.accessTime)
}

func (item *cacheItem) isExpired(expiration time.Duration) bool {
	return time.Now().UnixNano() >= item.getAccessTime()+int64(expiration)
}

func (item *cacheItem) isLoading() bool {
	return atomic.LoadInt32(&item.loading) == 1
}

func (item *cacheItem) setLoading(loading bool) {
	if loading {
		atomic.StoreInt32(&item.loading, 1)
	} else {
		atomic.StoreInt32(&item.loading, 0)
	}
}
