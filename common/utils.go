package common

import (
	"github.com/weikaishio/redis_orm"
	"strconv"
	"time"
	"encoding/base64"
	"crypto/rc4"
	"github.com/mkideal/log"
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

func DescryptRC4Base64(p, keystr string) []byte {
	key := []byte(keystr)
	str, err := base64.StdEncoding.DecodeString(p)
	if err != nil {
		log.Info("DescryptBase64 debase64 err:%v", err)
		return nil
	}
	data := []byte(str)
	ct, err := rc4.NewCipher(key)
	if err != nil {
		log.Info("DescryptBase64  err:%v", err)
		return nil
	}
	dst := make([]byte, len(data))
	ct.XORKeyStream(dst, data)
	return dst
}
func EncryptRC4Base64(p []byte, key string) string {
	k := []byte(key)
	cl, _ := rc4.NewCipher(k)
	dst := make([]byte, len(p))
	cl.XORKeyStream(dst, p)
	str := base64.StdEncoding.EncodeToString(dst)
	return str
}