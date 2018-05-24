package utils

import "time"

// 获取当前时间截
func TimestampNano() int64 {
	return time.Now().UnixNano()
}