/********************************************************************
created:    2018-09-09
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/
package loom

import "sync/atomic"

type Bool int32

func (iam *Bool) Store(value bool) {
	if value {
		atomic.StoreInt32((*int32)(iam), 1)
	} else {
		atomic.StoreInt32((*int32)(iam), 0)
	}
}

func (iam *Bool) Load() bool {
	if atomic.LoadInt32((*int32)(iam)) != 0 {
		return true
	}

	return false
}
