package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/mem"
	"github.com/weikaishio/redis_orm"
	"net/http"
	"runtime"
	"strings"
)

func Index(c *gin.Context) {
	title := "Redis ORM's workbench"
	systemInfo := make(map[string]string)

	systemInfo["os"] = strings.ToUpper(runtime.GOOS + " " + runtime.GOARCH)

	systemInfo["go_varsion"] = strings.ToUpper(runtime.Version())

	systemInfo["gin_varsion"] = strings.ToUpper(gin.Version)

	memoryInfo, _ := mem.VirtualMemory()
	systemInfo["main_server_total_memory"] = redis_orm.ToString(memoryInfo.Total)
	systemInfo["main_server_free_memory"] = redis_orm.ToString(int(memoryInfo.Free))
	systemInfo["main_server_available_memory"] = redis_orm.ToString(int(memoryInfo.Available))
	systemInfo["main_server_UsedPercent_memory"] = fmt.Sprintf("%10.2f%%", memoryInfo.UsedPercent)

	tableMap := redisORMSchemaBiz.LoadTables()
	var tables []string
	for tableName, _ := range tableMap {
		tables = append(tables, tableName)
	}
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"title":      title,
		"tables":     tables,
		"systemInfo": systemInfo,
		"memoryInfo": memoryInfo,
	})
}
