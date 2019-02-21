package controller

import (
	"fmt"
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
		idxNameKey     = c.PostForm("idxNameKey")
		individualVal  = c.PostForm("individualVal")
		startNumber, _ = strconv.Atoi(c.PostForm("startNumber"))
		endNumber, _   = strconv.Atoi(c.PostForm("endNumber"))
		startTime      = c.PostForm("startTime")
		endTime        = c.PostForm("endTime")
		ctype, _       = strconv.Atoi(c.PostForm("ctype"))

		individualVal_2  = c.PostForm("individualVal_2")
		startNumber_2, _ = strconv.Atoi(c.PostForm("startNumber_2"))
		endNumber_2, _   = strconv.Atoi(c.PostForm("endNumber_2"))
		startTime_2      = c.PostForm("startTime_2")
		endTime_2        = c.PostForm("endTime_2")
		ctype_2, _       = strconv.Atoi(c.PostForm("ctype_2"))
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
	timeStart_2, _ := time.ParseInLocation("2006-01-02 15:04:05", startTime_2, time.Local)
	timeEnd_2, _ := time.ParseInLocation("2006-01-02 15:04:05", endTime_2, time.Local)
	condition := &models.DataConditionInfo{
		CType:            ctype,
		IdxNameKey:       idxNameKey,
		IndividualValue:  individualVal,
		StartTime:        uint32(timeStart.Unix()),
		EndTime:          uint32(timeEnd.Unix()),
		StartNumber:      startNumber,
		EndNumber:        endNumber,
		CType2:           ctype_2,
		IndividualValue2: individualVal_2,
		StartTime2:       uint32(timeStart_2.Unix()),
		EndTime2:         uint32(timeEnd_2.Unix()),
		StartNumber2:     startNumber_2,
		EndNumber2:       endNumber_2,
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
		"primaryKey": table.PrimaryKey,
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

		"startTime_2":     startTime_2,
		"endTime_2":       endTime_2,
		"individualVal_2": individualVal_2,
		"startNumber_2":   startNumber_2,
		"endNumber_2":     endNumber_2,
		"ctype_2":         ctype_2,
	})
	return
}
func DataDel(c *gin.Context) {
	tbName, has := c.GetQuery("table_name")
	if !has {
		c.JSON(http.StatusOK, map[string]string{"statusCode": "300",
			"message":  "参数不对",
			"navTabId": "data_" + tbName})
		return
	}
	pkIdStr, has := c.GetQuery("pk_id")
	if !has {
		c.JSON(http.StatusOK, map[string]string{"statusCode": "300",
			"message":  "参数不对",
			"navTabId": "data_" + tbName})
		return
	}
	var pkId int64
	redis_orm.SetInt64FromStr(&pkId, pkIdStr)
	has, table, _ := redisORMSchemaBiz.BuildSchemaColumnsInfo(tbName)
	if !has {
		c.HTML(http.StatusBadRequest, "data_list.tmpl", gin.H{})
		return
	}
	err := redisORMDataBiz.Del(table, pkId)
	if err != nil {
		c.JSON(http.StatusOK, map[string]string{"statusCode": "300",
			"message":  "删除失败：" + err.Error(),
			"navTabId": "data_" + tbName})
	} else {
		c.JSON(http.StatusOK, map[string]string{"statusCode": "200",
			"message":  fmt.Sprint("删除成功：删除的ID=%d", pkId),
			"navTabId": "data_" + tbName})
	}
	//不用bean，直接传table和map~
	//schemaColumnsInfo := models.SchemaColumnsInfo{ColumnName: "colName"}
	//bys, _ := json.Marshal(schemaColumnsInfo)
	//val := reflect.New(reflect.TypeOf(schemaColumnsInfo)).Interface()
	//json.Unmarshal(bys, &val)
	//bys2,_:=json.Marshal(val)
	//fmt.Printf("val:%s\nval:%s\ntyp:%v\n",string(bys), string(bys2),reflect.TypeOf(schemaColumnsInfo))
}
func DataEdit(c *gin.Context) {
	tbName, has := c.GetQuery("table_name")
	if !has {
		c.JSON(http.StatusOK, map[string]string{"statusCode": "300",
			"message":  "参数不对",
			"navTabId": "data_" + tbName})
		return
	}
	//todo:
	c.HTML(http.StatusOK, "data_edit.tmpl", gin.H{
		"tableName":  tbName,
	})
	return
}
