/********************************************************************
created:    2018-11-29
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/
package safe

import (
	"unsafe"
	"sync/atomic"
)

type String struct {
	p unsafe.Pointer
}

func (str *String) Store(value string) {
	var val = unsafe.Pointer(&value)
	atomic.StorePointer(&str.p, val)
}

func (str *String) Load() string {
	var val = atomic.LoadPointer(&str.p)
	if val != nil {
		return *(*string)(val)
	}

	return ""
}
