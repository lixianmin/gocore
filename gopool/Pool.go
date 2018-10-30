/********************************************************************
created:    2018-10-30
author:     lixianmin

https://cloud.tencent.com/info/38078007c4e29f0acb54a33cb68c2543.html

Copyright (C) - All Rights Reserved
*********************************************************************/
package gopool

type Pool struct {
	work chan func()
	sem  chan struct{}
}

func New(size int) *Pool {
	var pool = &Pool{
		work: make(chan func()),
		sem:  make(chan struct{}, size),
	}

	return pool
}

func (p *Pool) Schedule(task func()) {
	select {
	case p.work <- task:
	case p.sem <- struct{}{}:
		go p.worker(task)
	}
}

func (p *Pool) worker(task func()) {
	defer func() {
		<-p.sem
	}()

	for {
		task()
		task = <-p.work
	}
}
