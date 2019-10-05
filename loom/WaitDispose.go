/********************************************************************
created:    2018-10-09
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/
package loom

import (
	"sync"
	"sync/atomic"
)

type WaitDispose struct {
	DisposeChan chan struct{}
	m           sync.Mutex
	done        int32
}

func NewWaitDispose() *WaitDispose {
	var wd = &WaitDispose{
		DisposeChan: make(chan struct{}),
	}

	return wd
}

func (wd *WaitDispose) Dispose() {
	if 0 == atomic.LoadInt32(&wd.done) {
		wd.m.Lock()
		if 0 == wd.done {
			atomic.StoreInt32(&wd.done, 1)
			close(wd.DisposeChan)
		}
		wd.m.Unlock()
	}
}

func (wd *WaitDispose) IsDisposed() bool {
	return atomic.LoadInt32(&wd.done) == 1
}
