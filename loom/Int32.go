/********************************************************************
created:    2018-10-21
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/
package loom

import "sync/atomic"

type Int32 int32

func (iam *Int32) Store(value int32) {
	atomic.StoreInt32((*int32)(iam), value)
}

func (iam *Int32) Load() int32 {
	return atomic.LoadInt32((*int32)(iam))
}

func (iam *Int32) AddIf(delta int32, predicate func(old int32) bool) bool {
	var expect, update int32
	for {
		expect = atomic.LoadInt32((*int32)(iam))
		if !predicate(expect) {
			return false
		}

		update = expect + delta

		if atomic.CompareAndSwapInt32((*int32)(iam), expect, update) {
			return true
		}
	}
}
