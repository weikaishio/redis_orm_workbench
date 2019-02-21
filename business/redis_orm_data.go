package business

import (
	"errors"
	"fmt"
	"github.com/mkideal/log"
	"github.com/weikaishio/redis_orm"
	"github.com/weikaishio/redis_orm_workbench/models"
)

type RedisORMDataBusiness struct {
	redisORMEngine *redis_orm.Engine
}

func NewRedisORMDataBusiness(redisORMEngine *redis_orm.Engine) *RedisORMDataBusiness {
	return &RedisORMDataBusiness{
		redisORMEngine: redisORMEngine,
	}
}

func (this *RedisORMDataBusiness) Query(condition *models.DataConditionInfo, offset, limit int, table *redis_orm.Table, cols ...string) ([]map[string]interface{}, int64, error) {
	searchCon := &redis_orm.SearchCondition{}
	if condition.IdxNameKey != "" {
		for _, idx := range table.IndexesMap {
			if idx.NameKey == condition.IdxNameKey {

				switch condition.CType {
				case models.CType_IndividualValue:
					searchCon.FieldMinValue = condition.IndividualValue
					searchCon.FieldMaxValue = condition.IndividualValue
				case models.CType_Number:
					searchCon.FieldMinValue = condition.StartNumber
					searchCon.FieldMaxValue = condition.EndNumber
				case models.CType_Time:
					searchCon.FieldMinValue = condition.StartTime
					searchCon.FieldMaxValue = condition.EndTime
				default:
					return nil, 0, errors.New("未知的查询条件")
				}

				if len(idx.IndexColumn) == 2 {
					switch condition.CType2 {
					case models.CType_IndividualValue:
						if idx.Type == redis_orm.IndexType_IdMember {
							var val, val2, individualVal int64
							redis_orm.SetInt64FromStr(&val, redis_orm.ToString(searchCon.FieldMinValue))
							redis_orm.SetInt64FromStr(&val2, redis_orm.ToString(searchCon.FieldMaxValue))
							redis_orm.SetInt64FromStr(&individualVal, redis_orm.ToString(condition.IndividualValue2))
							searchCon.FieldMinValue = val<<32 | individualVal
							searchCon.FieldMaxValue = val2<<32 | individualVal
						} else {
							searchCon.FieldMinValue = fmt.Sprintf("%s&%s", searchCon.FieldMinValue, condition.IndividualValue2)
							searchCon.FieldMaxValue = searchCon.FieldMinValue
						}
					case models.CType_Number:
						if idx.Type == redis_orm.IndexType_IdMember {
							var val, val2 int
							redis_orm.SetIntFromStr(&val, redis_orm.ToString(searchCon.FieldMinValue))
							redis_orm.SetIntFromStr(&val2, redis_orm.ToString(searchCon.FieldMaxValue))

							searchCon.FieldMinValue = val<<32 | condition.StartNumber2
							searchCon.FieldMaxValue = val2<<32 | condition.EndNumber2
						} else {
							searchCon.FieldMinValue = fmt.Sprintf("%s&%s", searchCon.FieldMinValue, condition.StartNumber2)
							searchCon.FieldMaxValue = fmt.Sprintf("%s&%s", searchCon.FieldMaxValue, condition.EndNumber2)
						}
					case models.CType_Time:
						if idx.Type == redis_orm.IndexType_IdMember {
							var val, val2 int
							redis_orm.SetIntFromStr(&val, redis_orm.ToString(searchCon.FieldMinValue))
							redis_orm.SetIntFromStr(&val2, redis_orm.ToString(searchCon.FieldMaxValue))

							searchCon.FieldMinValue = val<<32 | int(condition.StartTime2)
							searchCon.FieldMaxValue = val2<<32 | int(condition.EndTime2)
						} else {
							searchCon.FieldMinValue = fmt.Sprintf("%s&%s", searchCon.FieldMinValue, condition.StartTime2)
							searchCon.FieldMaxValue = fmt.Sprintf("%s&%s", searchCon.FieldMaxValue, condition.EndTime2)
						}
					default:
						return nil, 0, errors.New("未知的查询条件")
					}
				}

				searchCon.SearchColumn = idx.IndexColumn
				searchCon.IsAsc = false
				break
			}
		}
	}
	if len(searchCon.SearchColumn) == 0 {
		searchCon.SearchColumn = []string{table.PrimaryKey}
		searchCon.FieldMinValue = redis_orm.ScoreMin
		searchCon.FieldMaxValue = redis_orm.ScoreMax
	}
	log.Trace("seachCon:%v", *searchCon) //236223201284 236223201335
	val, count, err := this.redisORMEngine.Query(int64(offset), int64(limit), searchCon, table)
	if err != nil {
		log.Error("Query(%d,%d,searchCon:%v,tableName:%s) err:%v")
	}
	return val, count, err
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
	count, err := this.redisORMEngine.Find(offset, limit, searchCon, &resultAry)
	if err != nil {
		return nil, 0, err
	}

	if err != nil {
		log.Error("Count err:%v", err)
	}
	return resultAry, count, err
}

func (this *RedisORMDataBusiness) Del(table *redis_orm.Table, pkInt int64) error {
	err := this.redisORMEngine.DeleteByPK(table, pkInt)
	return err
}

func (this *RedisORMDataBusiness) Edit(table *redis_orm.Table, columnValMap map[string]string) error {
	return this.redisORMEngine.UpdateByMap(table,columnValMap)
}

func (this *RedisORMDataBusiness) Insert(table *redis_orm.Table, columnValMap map[string]string) error {
	return this.redisORMEngine.InsertByMap(table,columnValMap)
}