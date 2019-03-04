package controller

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/weikaishio/redis_orm"
	"net/http"
	"sort"
	"strings"
)

func Schema(c *gin.Context) {
	tbName, has := c.GetQuery("table_name")
	if !has {
		c.JSON(http.StatusOK, map[string]string{"statusCode": "300",
			"message":  fmt.Sprintf("table:%s, not exist", tbName),
			"navTabId": "data_" + c.Query("table_name")})
		return
	}
	has, table, columns := redisORMSchemaBiz.BuildSchemaColumnsInfo(tbName)
	if !has {
		c.JSON(http.StatusOK, map[string]string{"statusCode": "300",
			"message":  fmt.Sprintf("table:%s, not exist", tbName),
			"navTabId": "data_" + c.Query("table_name")})
		return
	}
	if len(columns) > 0 {
		sort.Sort(columns)
	}
	c.HTML(http.StatusOK, "schema.tmpl", gin.H{
		"table_name": redis_orm.Underline2Camel(table.Name),
		"columns":    columns,
	})
}

func CreateTable(c *gin.Context) {
	if strings.ToLower(c.Request.Method) == "get" {
		c.HTML(http.StatusOK, "table_create.tmpl", gin.H{})
	} else if strings.ToLower(c.Request.Method) == "post" {
		val := strings.Trim(c.PostForm("txt_content"), " ")
		if val == "" {
			c.JSON(http.StatusOK, map[string]string{"statusCode": "300",
				"message":  "content could not be null",
				"navTabId": "data_" + c.Query("table_name")})
			return
		}
		err := redisORMSchemaBiz.CreateTable(val)
		if err != nil {
			c.JSON(http.StatusOK, map[string]string{"statusCode": "300",
				"message":  "处理失败：" + err.Error(),
				"navTabId": "data_create_table"})
		} else {
			c.JSON(http.StatusOK, map[string]string{"statusCode": "200",
				"message":  "处理成功",
				"navTabId": "data_create_table"})
		}
	}
	return
}
func DropTable(c *gin.Context) {
	table, err := VerifyTable(c)
	if err != nil {
		c.JSON(http.StatusOK, map[string]string{"statusCode": "300",
			"message":  err.Error(),
			"navTabId": "data_" + c.Query("table_name")})
		return
	}
	//err = redisORMDataBiz.DropTable(table)
	err = errors.New("太危险了，功能先不放出来")
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
func AddColumn(c *gin.Context) {

}
func DropColumn(c *gin.Context) {

}
