package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/mkideal/log"
	"github.com/weikaishio/redis_orm"
	"github.com/weikaishio/redis_orm_workbench/models"
	"net/http"
	"sort"
	"strconv"
	"time"
)

func DataList(c *gin.Context) {
	var (
		numPerPage, _  = strconv.Atoi(c.PostForm("numPerPage"))
		pageNum, _     = strconv.Atoi(c.PostForm("pageNum"))
		startTime      = c.PostForm("startTime")
		endTime        = c.PostForm("endTime")
		idxNameKey     = c.PostForm("idxNameKey")
		individualVal  = c.PostForm("individualVal")
		startNumber, _ = strconv.Atoi(c.PostForm("startNumber"))
		endNumber, _   = strconv.Atoi(c.PostForm("endNumber"))
		ctype, _       = strconv.Atoi(c.PostForm("ctype"))
	)
	if pageNum == 0 {
		pageNum = 1
	}
	if numPerPage < 5 {
		numPerPage = 15
	}
	log.Info("numPerPage:%d, pageNum:%d,startTime:%v,endTime:%v", numPerPage, pageNum, startTime, endTime)
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
	timeStart, _ := time.ParseInLocation("2006-01-02 15:04:05", startTime, time.Local)
	timeEnd, _ := time.ParseInLocation("2006-01-02 15:04:05", endTime, time.Local)
	condition := &models.DataConditionInfo{
		CType:           ctype,
		IdxNameKey:      idxNameKey,
		IndividualValue: individualVal,
		StartTime:       uint32(timeStart.Unix()),
		EndTime:         uint32(timeEnd.Unix()),
		StartNumber:     startNumber,
		EndNumber:       endNumber,
	}
	valMap, count, err := redisORMDataBiz.Query(condition, (pageNum-1)*numPerPage, numPerPage, table)
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
		"tableName":  table.Name,
		"indexs":     table.IndexesMap,
		"columns":    columns,
		"numPerPage": numPerPage,
		"pageNum":    pageNum,
		"totalCount": count,
		"valAry":     valAry,

		"startTime":     startTime,
		"endTime":       endTime,
		"idxNameKey":    idxNameKey,
		"individualVal": individualVal,
		"startNumber":   startNumber,
		"endNumber":     endNumber,
		"ctype":         ctype,
	})
	return
}
