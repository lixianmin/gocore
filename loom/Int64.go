/********************************************************************
created:    2018-12-03
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/
package loom

import "sync/atomic"

type Int64 int64

func (iam *Int64) Store(value int64) {
	atomic.StoreInt64((*int64)(iam), value)
}

func (iam *Int64) Load() int64 {
	return atomic.LoadInt64((*int64)(iam))
}