package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/mkideal/log"
	"github.com/weikaishio/redis_orm_workbench/business"
	"github.com/weikaishio/redis_orm_workbench/common"
	"github.com/weikaishio/redis_orm_workbench/config"
	"strings"
)

var (
	redisORMSchemaBiz *business.RedisORMSchemaBusiness
	redisORMDataBiz   *business.RedisORMDataBusiness
)

func InitBiz() {
	redisORMSchemaBiz = business.NewRedisORMSchemaBusiness(config.Cfg.RedisORM)
	redisORMDataBiz = business.NewRedisORMDataBusiness(config.Cfg.RedisORM)
}

func UseMiddleware(ctx *gin.Context) {
	url := ctx.Request.URL.String()
	if strings.HasPrefix(url, "/static/") ||
		strings.HasPrefix(url, "/login") {
		ctx.Next()
	} else {
		val, err := ctx.Cookie("user")
		for {
			if err != nil {
				log.Trace("ctx.Cookie(user),err:%v", err)
				break
			}
			valAry := strings.Split(val, ":")
			password, has := config.Cfg.UserMap[valAry[0]]
			if !has {
				log.Trace("config.Cfg.UserMap[%s] !has", valAry[0])
				break
			}
			if len(valAry) != 2 || valAry[1] != common.EncryptRC4Base64([]byte(valAry[0]), password) {
				log.Trace("cookie value is invalid:%s", valAry[1])
				break
			}
			ctx.Next()
			return
		}
		Login(ctx)
		ctx.Abort()
	}
}
