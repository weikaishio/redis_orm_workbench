package business

import (
	"fmt"
	"github.com/mkideal/log"
	"github.com/weikaishio/redis_orm"
	"github.com/weikaishio/redis_orm/table_from_ast"
	"github.com/weikaishio/redis_orm_workbench/models"
	"reflect"
	"strings"
	"encoding/json"
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
		if schemaColumnTb.ColumnName == table.Created {
			tags = append(tags, redis_orm.TagCreatedAt)
		}
		if schemaColumnTb.ColumnName == table.Updated {
			tags = append(tags, redis_orm.TagUpdatedAt)
		}
		if schemaColumnTb.DefaultValue != "" {
			tags = append(tags, redis_orm.TagDefaultValue)
			tags = append(tags, fmt.Sprintf("'%s'", schemaColumnTb.DefaultValue))
		}
		if schemaColumnTb.ColumnComment != "" {
			tags = append(tags, redis_orm.TagComment)
			tags = append(tags, fmt.Sprintf("'%s'", schemaColumnTb.ColumnComment))
		}
		schemaColumn.Tags = strings.Join(tags, " ")
		columns = append(columns, schemaColumn)
	}
	return ok, table, columns
}

//ast from table struct
func (this *RedisORMSchemaBusiness) CreateTable(tableDefSchema string) error {
	tableDefSchema = "package schema\n" + tableDefSchema
	tables, err := table_from_ast.TableFromAst("", tableDefSchema)
	if err != nil {
		log.Error("CreateTable %s,err:%v", tableDefSchema, err)
		return err
	}
	var successTable []*redis_orm.Table
	for i, table := range tables {
		valTb, _ := json.Marshal(table)
		log.Info("valTb:%v\n", string(valTb))
		err = this.redisORMEngine.Schema.CreateTableByTable(table)
		if err != nil {
			if err != redis_orm.Err_DataHadAvailable {
				for j, tableRollback := range successTable {
					if i == j {
						break
					}
					this.redisORMEngine.Schema.TableDrop(tableRollback)
				}
				log.Error("CreateTable %s,err:%v", tableDefSchema, err)
				return err
			} else {
				log.Warn("CreateTable %s, table:%s HasExist", tableDefSchema, table.Name)
			}
		} else {
			successTable = append(successTable, table)
		}
	}
	return nil
}

//
func (this *RedisORMSchemaBusiness) AlterTable() {

}

//
func (this *RedisORMSchemaBusiness) DropTable(table *redis_orm.Table) error {
	return this.redisORMEngine.Schema.TableDrop(table)
}
