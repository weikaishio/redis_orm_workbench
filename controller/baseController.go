package controller

import (
	"github.com/go-redis/redis"
	"github.com/weikaishio/redis_orm"
	"github.com/weikaishio/redis_orm_workbench/business"
)

var (
	redisORM          *redis_orm.Engine
	redisORMSchemaBiz *business.RedisORMSchemaBusiness
)

func InitBiz() {
	options := redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       1,
	}

	redisClient := redis.NewClient(&options)
	redisORM = redis_orm.NewEngine(redisClient)
	redisORM.IsShowLog(true)
	redisORM.Schema.ReloadTables()
	redisORMSchemaBiz = business.NewRedisORMSchemaBusiness(redisORM)
}
