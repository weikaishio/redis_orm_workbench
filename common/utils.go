package common

import (
	"github.com/weikaishio/redis_orm"
	"strconv"
	"time"
)

func FormatInterface2Time(val interface{}) string {
	switch v := val.(type) {
	case int64:
		return FormatTime(v)
	default:
		valStr := redis_orm.ToString(val)
		if valStr != "" {
			var timeUnix int64
			err := redis_orm.SetInt64FromStr(&timeUnix, valStr)
			if err == nil && timeUnix > 0 {
				return FormatTime(timeUnix)
			}
		}
		return valStr
	}
}
func FormatTime(timeUnix int64) string {
	timeUnixStr := strconv.FormatInt(timeUnix, 10)
	formatedTime := timeUnixStr
	if timeUnix > 0 {
		switch len(timeUnixStr) {
		case 10:
			formatedTime = time.Unix(timeUnix, 0).Format("2006-01-02 15:04:05")
		case 13:
			//毫秒
			formatedTime = time.Unix(timeUnix/1e3, 0).Format("2006-01-02 15:04:05")
		case 16:
			//微秒
			formatedTime = time.Unix(timeUnix/1e6, 0).Format("2006-01-02 15:04:05")
		case 19:
			//纳秒
			formatedTime = time.Unix(timeUnix/1e9, 0).Format("2006-01-02 15:04:05")
		default:
		}
	}
	return formatedTime
}
