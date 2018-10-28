/********************************************************************
created:    2018-10-21
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/
package safe

import "sync/atomic"

type Int32 int32

func (iam *Int32) Store(value int32) {
	atomic.StoreInt32((*int32)(iam), value)
}

func (iam *Int32) Load() int32 {
	return atomic.LoadInt32((*int32)(iam))
}
