/********************************************************************
created:    2018-10-30
author:     lixianmin

https://cloud.tencent.com/info/38078007c4e29f0acb54a33cb68c2543.html

Copyright (C) - All Rights Reserved
*********************************************************************/
package loom

type GoroutinePool struct {
	work chan func()
	sem  chan struct{}
}

func NewGoroutinePool(size int) *GoroutinePool {
	var pool = &GoroutinePool{
		work: make(chan func()),
		sem:  make(chan struct{}, size),
	}

	return pool
}

func (p *GoroutinePool) Schedule(task func()) {
	select {
	case p.work <- task:
	case p.sem <- struct{}{}:
		go p.worker(task)
	}
}

func (p *GoroutinePool) worker(task func()) {
	defer DumpIfPanic("")
	defer func() {
		<-p.sem
	}()

	for {
		task()
		task = <-p.work
	}
}
