package business

import (
	"github.com/mkideal/log"
	"github.com/weikaishio/redis_orm"
)

type RedisORMDataBusiness struct {
	redisORMEngine *redis_orm.Engine
}

func NewRedisORMDataBusiness(redisORMEngine *redis_orm.Engine) *RedisORMDataBusiness {
	return &RedisORMDataBusiness{
		redisORMEngine: redisORMEngine,
	}
}

/*
find by table, return map array
get by table & id, return map
edit by table & map, return affected count
add by table & map, return map
*/
func (this *RedisORMDataBusiness) List(searchCon *redis_orm.SearchCondition, pageNum, pageSize int64) ([]interface{}, int64, error) {
	var (
		resultAry     []interface{}
		offset, limit int64 = 0, 100
	)

	if pageNum > 0 && pageSize > 0 {
		offset = (pageNum - 1) * pageSize
		limit = pageSize
	}
	_, err := this.redisORMEngine.Find(offset, limit, searchCon, &resultAry)
	if err != nil {
		return nil, 0, err
	}
	// 查询数量
	count, err := this.redisORMEngine.Count(searchCon, &resultAry)

	if err != nil {
		log.Error("Count err:%v", err)
	}
	return resultAry, count, err
}
