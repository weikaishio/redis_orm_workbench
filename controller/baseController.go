package controller

import (
	"github.com/go-redis/redis"
	"github.com/weikaishio/redis_orm"
	"github.com/weikaishio/redis_orm_workbench/business"
)

var (
	redisORM          *redis_orm.Engine
	redisORMSchemaBiz *business.RedisORMSchemaBusiness
	redisORMDataBiz *business.RedisORMDataBusiness
)

func InitBiz() {
	options := redis.Options{
		Addr:     "59.110.27.156:6379",
		Password: "testwashcar",
		DB:       1,
	}

	redisClient := redis.NewClient(&options)
	redisORM = redis_orm.NewEngine(redisClient)
	redisORM.IsShowLog(true)
	redisORMSchemaBiz = business.NewRedisORMSchemaBusiness(redisORM)
	redisORMDataBiz = business.NewRedisORMDataBusiness(redisORM)
}
