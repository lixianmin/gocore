/********************************************************************
created:    2018-12-12
author:     lixianmin

unix时间戳，不能使用int32存储，因为会有2038问题，参考链接如下：
https://baike.baidu.com/item/unix%E6%97%B6%E9%97%B4%E6%88%B3/2078227?fr=aladdin

Copyright (C) - All Rights Reserved
*********************************************************************/
package loom

import "time"

type UnixTime int64

func (t *UnixTime) ToTime() time.Time {
	return time.Unix(int64(*t), 0)
}

func fromTime(t time.Time) UnixTime {
	return UnixTime(t.Unix())
}

func fromDuration(d time.Duration) UnixTime {
	return UnixTime(d / 1e9)
}
