package data

import "time"

func Unix2Str() string {
	return time.Now().Local().Format("2006-01-02 15:04:05")
}