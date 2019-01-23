package controller

import (
	"github.com/weikaishio/redis_orm_workbench/business"
	"github.com/weikaishio/redis_orm_workbench/config"
)

var (
	redisORMSchemaBiz *business.RedisORMSchemaBusiness
	redisORMDataBiz   *business.RedisORMDataBusiness
)

func InitBiz() {
	redisORMSchemaBiz = business.NewRedisORMSchemaBusiness(config.Cfg.RedisORM)
	redisORMDataBiz = business.NewRedisORMDataBusiness(config.Cfg.RedisORM)
}
