/********************************************************************
created:    2018-10-14
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

package loom

import (
	"fmt"
	"testing"
	"time"
)

func TestLongTimeFunction (t *testing.T) {
	fmt.Println("startTime:" + time.Now().String())

	Repeat(time.Second, func() {
		time.Sleep(time.Second * 2)
		fmt.Println(time.Now().String())
	})

	time.Sleep(time.Minute)
}