/********************************************************************
created:    2018-12-09
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

package loom

import (
	"fmt"
	"testing"
	"time"
)

func TestAsyncCache_Get(t *testing.T) {
	var cache = NewAsyncCache(32, 1*time.Second)

	for i := 0; i < 1000; i++ {
		i := i
		go func() {
			for j := 0; j < 10; j++ {
				var val = cache.Get(i, func() interface{} {
					fmt.Printf("reload: i=%d\n\n", i)
					return i * (j + 1)
				})

				time.Sleep(time.Millisecond * 200)
				fmt.Printf("i=%d, val=%v\n", i, val)
			}
		}()
	}

	for cache.GetCount() > 0 {
		fmt.Printf("cache.size=%d\n", cache.GetCount())
		time.Sleep(time.Second)
	}

	//time.Sleep(time.Minute * 2)
}
