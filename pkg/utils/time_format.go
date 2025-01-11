package utils

import "time"

// UnixSecondToTime 将时间戳转换为 time.Time 类型
func UnixSecondToTime(second int64) time.Time {
	return time.Unix(second, 0)
}

func GetCurrentTimestampByMill() int64 {
	return time.Now().UnixNano() / 1e6
}
