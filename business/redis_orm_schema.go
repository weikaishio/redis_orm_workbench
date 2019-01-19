package business

import (
	"fmt"
	"github.com/weikaishio/redis_orm"
	"github.com/weikaishio/redis_orm_workbench/models"
	"reflect"
	"strings"
)

type RedisORMSchemaBusiness struct {
	redisORMEngine *redis_orm.Engine
}

func NewRedisORMSchemaBusiness(redisORMEngine *redis_orm.Engine) *RedisORMSchemaBusiness {
	return &RedisORMSchemaBusiness{
		redisORMEngine: redisORMEngine,
	}
}
func (this *RedisORMSchemaBusiness) LoadTables() map[string]*redis_orm.Table {
	return this.redisORMEngine.Tables
}
func (this *RedisORMSchemaBusiness) BuildSchemaColumnsInfo(tableName string) (bool, *redis_orm.Table, models.ColumnsSortModel) {
	tables := this.LoadTables()

	var columns models.ColumnsSortModel
	table, ok := tables[tableName]
	if !ok {
		return ok, table, columns
	}
	for _, idx := range table.IndexesMap {
		if len(idx.IndexColumn) > 1 {
			schemaColumn := &models.SchemaColumnsInfo{
				ColumnName: strings.Join(idx.IndexColumn, ""),
				Seq:        idx.Seq,
			}
			if idx.Type == redis_orm.IndexType_IdScore {
				schemaColumn.DataType = reflect.String.String()
			} else {
				schemaColumn.DataType = reflect.Int64.String()
			}
			var tags []string
			tags = append(tags, redis_orm.TagCombinedindex)
			tags = append(tags, strings.Join(idx.IndexColumn, "&"))
			if idx.Comment != "" {
				tags = append(tags, redis_orm.TagComment)
				tags = append(tags, fmt.Sprintf("'%s'", idx.Comment))
			}

			schemaColumn.Tags = strings.Join(tags, " ")
			schemaColumn.Comment = idx.Comment
			columns = append(columns, schemaColumn)
		}
	}
	for _, column := range table.ColumnsMap {
		schemaColumnTb := redis_orm.SchemaColumnsFromColumn(table.TableId, column)
		schemaColumn := &models.SchemaColumnsInfo{
			ColumnName: schemaColumnTb.ColumnName,
			DataType:   schemaColumnTb.DataType,
			Seq:        schemaColumnTb.Seq,
			Comment:    schemaColumnTb.ColumnComment,
		}
		var tags []string
		if schemaColumnTb.ColumnName == table.PrimaryKey {
			tags = append(tags, redis_orm.TagPrimaryKey)
			if table.IsSync2DB {
				tags = append(tags, redis_orm.TagSync2DB)
			}
		} else {
			for k, idx := range table.IndexesMap {
				if column.IsCombinedIndex {
					if strings.ToLower(schemaColumnTb.ColumnName) == strings.Join(idx.IndexColumn, "") {
						tags = append(tags, redis_orm.TagCombinedindex)
						tags = append(tags, strings.Join(idx.IndexColumn, "&"))
						break
					}
				} else if k == strings.ToLower(schemaColumnTb.ColumnName) {
					if idx.IsUnique {
						tags = append(tags, redis_orm.TagUniqueIndex)
					} else {
						tags = append(tags, redis_orm.TagIndex)
					}
					break
				}
			}
		}
		if schemaColumnTb.ColumnName == table.AutoIncrement {
			tags = append(tags, redis_orm.TagAutoIncrement)
		}
		if schemaColumnTb.ColumnComment != "" {
			tags = append(tags, redis_orm.TagComment)
			tags = append(tags, fmt.Sprintf("'%s'", schemaColumnTb.ColumnComment))
		}
		if schemaColumnTb.ColumnName == table.Created {
			tags = append(tags, redis_orm.TagCreatedAt)
		}
		if schemaColumnTb.ColumnName == table.Updated {
			tags = append(tags, redis_orm.TagUpdatedAt)
		}
		if schemaColumnTb.DefaultValue != "" {
			tags = append(tags, redis_orm.TagDefaultValue)
			tags = append(tags, schemaColumnTb.DefaultValue)
		}
		schemaColumn.Tags = strings.Join(tags, " ")
		columns = append(columns, schemaColumn)
	}
	return ok, table, columns
}
