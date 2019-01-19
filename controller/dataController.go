package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/mkideal/log"
	"github.com/weikaishio/redis_orm"
	"github.com/weikaishio/redis_orm_workbench/models"
	"net/http"
	"sort"
	"strconv"
)

func DataList(c *gin.Context) {
	var (
		pageSize, _    = strconv.Atoi(c.PostForm("page_size"))
		pageIndex, _   = strconv.Atoi(c.PostForm("page_index"))
		startTime, _   = strconv.Atoi(c.PostForm("startTime"))
		endTime, _     = strconv.Atoi(c.PostForm("endTime"))
		columnName     = c.PostForm("columnName")
		individualVal  = c.PostForm("individualVal")
		startNumber, _ = strconv.Atoi(c.PostForm("startNumber"))
		endNumber, _   = strconv.Atoi(c.PostForm("endNumber"))
		ctype, _       = strconv.Atoi(c.PostForm("ctype"))
	)
	log.Info("pageSize:%d, pageIndex:%d", pageSize, pageIndex)
	tbName, has := c.GetQuery("table_name")
	if !has {
		c.HTML(http.StatusBadRequest, "data_list.tmpl", gin.H{})
		return
	}
	has, table, _ := redisORMSchemaBiz.BuildSchemaColumnsInfo(tbName)
	if !has {
		c.HTML(http.StatusBadRequest, "data_list.tmpl", gin.H{})
		return
	}
	condition := &models.DataConditionInfo{
		CType:           ctype,
		ColumnName:      columnName,
		IndividualValue: individualVal,
		StartTime:       uint32(startTime),
		EndTime:         uint32(endTime),
		StartNumber:     startNumber,
		EndNumber:       endNumber,
	}
	valMap, count, err := redisORMDataBiz.Query(condition, pageIndex*pageSize, pageSize, table)
	if err != nil {
		c.HTML(http.StatusBadRequest, "data_list.tmpl", gin.H{})
		return
	}

	//var indexs []*redis_orm.Index
	//for _,v:=range table.IndexesMap{
	//	indexs=append(indexs,v)
	//}
	//log.Info("index:%v",indexs)
	var columns redis_orm.ColumnsModel
	for _, column := range table.ColumnsMap {
		columns = append(columns, column)
	}
	if len(columns) > 0 {
		sort.Sort(columns)
	}
	var valAry [][]interface{}
	for _, v := range valMap {
		var vals []interface{}
		for _, column := range columns {
			vals = append(vals, v[column.Name])
		}
		valAry = append(valAry, vals)
	}

	c.HTML(http.StatusOK, "data_list.tmpl", gin.H{
		"table_name":  table.Name,
		"indexs":      table.IndexesMap,
		"columns":     columns,
		"condition":   condition,
		"page_size":   pageSize,
		"page_index":  pageIndex,
		"total_count": count,
		"val_ary":     valAry,
	})
}
