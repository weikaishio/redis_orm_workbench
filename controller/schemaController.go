package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sort"
)

func Schema(c *gin.Context) {
	tbName, has := c.GetQuery("table_name")
	if !has {
		c.HTML(http.StatusBadRequest, "schema.tmpl", gin.H{})
		return
	}
	has, table, columns := redisORMSchemaBiz.BuildSchemaColumnsInfo(tbName)
	if len(columns) > 0 {
		sort.Sort(columns)
	}
	c.HTML(http.StatusOK, "schema.tmpl", gin.H{
		"table_name": table.Name,
		"columns":    columns,
	})
}

func CreateTable(c *gin.Context){

}
func RebuildIndex(c *gin.Context){

}
func AddColumn(c *gin.Context){
//engine not support
}