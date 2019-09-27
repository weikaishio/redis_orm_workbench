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
	"strings"
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
	var (
		columnsPre, columns redis_orm.ColumnsModel
	)
	for _, column := range table.ColumnsMap {
		columnsPre = append(columnsPre, column)
	}
	if len(columnsPre) > 0 {
		sort.Sort(columnsPre)
	}
	for i, column := range columnsPre {
		if i > 9 && !column.IsCreated && !column.IsUpdated {
			continue
		}
		columns = append(columns, column)
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

func VerifyTable(c *gin.Context) (table *redis_orm.Table, err error) {
	tbName, has := c.GetQuery("table_name")
	if !has {
		err = fmt.Errorf("%s参数不对", "table_name")
		return
	}

	tables := redisORMSchemaBiz.LoadTables()

	table, has = tables[tbName]
	if !has {
		err = fmt.Errorf("表:%s不存在", tbName)
		return
	}
	return
}
func VerifyTableAndPkInt(c *gin.Context) (pkId int64, table *redis_orm.Table, err error) {
	tbName, has := c.GetQuery("table_name")
	if !has {
		err = fmt.Errorf("%s参数不对", "table_name")
		return
	}
	pkIdStr, has := c.GetQuery("pk_id")
	if !has {
		err = fmt.Errorf("%s参数不对", "pk_id")
		return
	}
	_ = redis_orm.SetInt64FromStr(&pkId, pkIdStr)

	tables := redisORMSchemaBiz.LoadTables()

	table, has = tables[tbName]
	if !has {
		err = fmt.Errorf("表:%s不存在", tbName)
		return
	}
	return
}
func DataDel(c *gin.Context) {
	pkId, table, err := VerifyTableAndPkInt(c)
	if err != nil {
		c.JSON(http.StatusOK, map[string]string{"statusCode": "300",
			"message":  err.Error(),
			"navTabId": "data_" + c.Query("table_name")})
		return
	}
	err = redisORMDataBiz.Del(table, pkId)
	if err != nil {
		c.JSON(http.StatusOK, map[string]string{"statusCode": "300",
			"message":  "删除失败：" + err.Error(),
			"navTabId": "data_" + table.Name})
	} else {
		c.JSON(http.StatusOK, map[string]string{"statusCode": "200",
			"message":  fmt.Sprintf("删除成功：删除的ID=%d", pkId),
			"navTabId": "data_" + table.Name})
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
	table, err := VerifyTable(c)
	if err != nil {
		c.JSON(http.StatusOK, map[string]string{"statusCode": "300",
			"message":  err.Error(),
			"navTabId": "data_" + c.Query("table_name")})
		return
	}
	var columns redis_orm.ColumnsModel
	for _, column := range table.ColumnsMap {
		columns = append(columns, column)
	}
	if len(columns) > 0 {
		sort.Sort(columns)
	}
	if strings.ToLower(c.Request.Method) == "get" {
		pkIdStr, hasId := c.GetQuery("pk_id")
		if !hasId {
			var vals []interface{}
			for _, column := range columns {
				if column.IsPrimaryKey {
					//if column.IsAutoIncrement {
					vals = append(vals, 0)
					//}
				} else if column.IsUpdated || column.IsCreated {
					vals = append(vals, time.Now().Unix())
				} else {
					vals = append(vals, column.DefaultValue)
				}
			}
			c.HTML(http.StatusOK, "data_edit.tmpl", gin.H{
				"tableName": table.Name,
				"valAry":    vals,
				"columns":   columns,
			})
			return
		}
		var pkId int64
		err = redis_orm.SetInt64FromStr(&pkId, pkIdStr)
		if err != nil {
			c.JSON(http.StatusOK, map[string]string{"statusCode": "300",
				"message":  fmt.Sprintf("%s参数值(%v)不对", "pk_id", pkIdStr),
				"navTabId": "data_" + c.Query("table_name")})
			return
		}

		valMap, has, err := redisORMDataBiz.Get(table, pkId)
		if err != nil {
			c.JSON(http.StatusOK, map[string]string{"statusCode": "300",
				"message":  err.Error(),
				"navTabId": "data_" + c.Query("table_name")})
			return
		}
		if !has {
			c.JSON(http.StatusOK, map[string]string{"statusCode": "300",
				"message":  fmt.Sprintf("数据(%d)不存在", pkId),
				"navTabId": "data_" + c.Query("table_name")})
			return
		}

		var vals []interface{}
		for _, column := range columns {
			vals = append(vals, valMap[column.Name])
		}

		c.HTML(http.StatusOK, "data_edit.tmpl", gin.H{
			"tableName": table.Name,
			"pk_id":     pkId,
			"valAry":    vals,
			"columns":   columns,
		})

	} else if strings.ToLower(c.Request.Method) == "post" {
		valMap := make(map[string]string)
		for colName, col := range table.ColumnsMap {
			v, has := c.GetPostForm(colName)
			if !has {
				valMap[colName] = col.DefaultValue
			} else {
				if strings.HasSuffix(colName, "At") && col.DataType == "int64" {
					if strings.Contains(v, "-") {
						t, err := time.ParseInLocation("2006-01-02 15:04:05", v, time.Local)
						if err != nil || t.IsZero() {
							valMap[colName] = ""
						} else {
							valMap[colName] = redis_orm.ToString(t.Unix())
						}
					} else {
						valMap[colName] = v
					}
				} else {
					valMap[colName] = v
				}
			}
		}
		log.Info("valMap:%v", valMap)
		pkIdStr, hasId := c.GetQuery("pk_id")
		if hasId {
			var pkId int64
			err = redis_orm.SetInt64FromStr(&pkId, pkIdStr)
			if err == nil && pkId > 0 {
				err = redisORMDataBiz.Edit(table, valMap)
			} else {
				err = redisORMDataBiz.Insert(table, valMap)
			}
		} else {
			err = redisORMDataBiz.Insert(table, valMap)
		}
		if err != nil {
			c.JSON(http.StatusOK, map[string]string{"statusCode": "300",
				"message":  err.Error(),
				"navTabId": "data_" + c.Query("table_name")})
			return
		} else {
			c.JSON(http.StatusOK, map[string]string{"statusCode": "200",
				"message":  "提交成功",
				"navTabId": "data_" + c.Query("table_name")})
		}
	}
	return
}

func RebuildIndex(c *gin.Context) {
	table, err := VerifyTable(c)
	if err != nil {
		c.JSON(http.StatusOK, map[string]string{"statusCode": "300",
			"message":  err.Error(),
			"navTabId": "data_" + c.Query("table_name")})
		return
	}
	err = redisORMDataBiz.RebuildIndex(table)
	if err != nil {
		c.JSON(http.StatusOK, map[string]string{"statusCode": "300",
			"message":  "处理失败：" + err.Error(),
			"navTabId": "data_" + table.Name})
	} else {
		c.JSON(http.StatusOK, map[string]string{"statusCode": "200",
			"message":  "处理成功",
			"navTabId": "data_" + table.Name})
	}
}

func TruncateTable(c *gin.Context) {
	table, err := VerifyTable(c)
	if err != nil {
		c.JSON(http.StatusOK, map[string]string{"statusCode": "300",
			"message":  err.Error(),
			"navTabId": "data_" + c.Query("table_name")})
		return
	}
	err = redisORMDataBiz.TruncateTable(table)
	//err = errors.New("太危险了，功能先不放出来")
	if err != nil {
		c.JSON(http.StatusOK, map[string]string{"statusCode": "300",
			"message":  "处理失败：" + err.Error(),
			"navTabId": "data_" + table.Name})
	} else {
		c.JSON(http.StatusOK, map[string]string{"statusCode": "200",
			"message":  "处理成功",
			"navTabId": "data_" + table.Name})
	}
}
