package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/mkideal/log"
	"github.com/mkideal/pkg/osutil/pid"
	"github.com/weikaishio/redis_orm/sync2db"
	"github.com/weikaishio/redis_orm_workbench/controller"
	"net/http"
)

func setupRouter() *gin.Engine {

	controller.InitBiz()

	r := gin.Default()

	r.LoadHTMLGlob("views/*")
	r.StaticFS("/static", http.Dir("static/"))
	r.GET("/monitor", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	r.GET("/index", controller.Index)
	r.GET("/schema", controller.Schema)
	return r
}

var (
	pidFile = flag.String("pid", "redis_orm_workbench.pid", "pid file")
)

func main() {
	log.SetLevelFromString("TRACE")

	if err := pid.New(*pidFile); err != nil {
		log.Fatal("pid.New: %v", err)
	}
	defer func() {
		pid.Remove(*pidFile)
	}()

	r := setupRouter()

	go func() { r.Run(":8881") }()

	log.Info("Start~")
	sync2db.ListenQuitAndDump()

	log.Info("Stop~")
}
