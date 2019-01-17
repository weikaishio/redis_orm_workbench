package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/weikaishio/redis_orm_workbench/models"
	"net/http"
	"sort"
	"strconv"
)

func DataList(c *gin.Context) {
	var (
		pageSize, _    = strconv.Atoi(c.PostForm("pageSize"))
		pageIndex, _   = strconv.Atoi(c.PostForm("pageIndex"))
		startTime, _   = strconv.Atoi(c.PostForm("startTime"))
		endTime, _     = strconv.Atoi(c.PostForm("endTime"))
		columnName     = c.PostForm("columnName")
		individualVal  = c.PostForm("individualVal")
		startNumber, _ = strconv.Atoi(c.PostForm("startNumber"))
		endNumber, _   = strconv.Atoi(c.PostForm("endNumber"))
		ctype, _       = strconv.Atoi(c.PostForm("ctype"))
	)
	tbName, has := c.GetQuery("table_name")
	if !has {
		c.HTML(http.StatusBadRequest, "data_list.tmpl", gin.H{})
		return
	}
	has, tableName, columns := redisORMSchemaBiz.BuildSchemaColumnsInfo(tbName)
	if len(columns) > 0 {
		sort.Sort(columns)
	}
	condition := models.DataConditionInfo{
		CType:           ctype,
		ColumnName:      columnName,
		IndividualValue: individualVal,
		StartTime:       uint32(startTime),
		EndTime:         uint32(endTime),
		StartNumber:     startNumber,
		EndNumber:       endNumber,
	}

	c.HTML(http.StatusOK, "data_list.tmpl", gin.H{
		"table_name": tableName,
		"columns":    columns,
		"condition":  condition,
		"pageSize":   pageSize,
		"pageIndex":  pageIndex,
	})
}
