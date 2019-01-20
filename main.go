package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"time"
	"html/template"

	"github.com/gin-gonic/gin"
	"github.com/mkideal/log"
	"github.com/mkideal/pkg/osutil/pid"
	"github.com/weikaishio/redis_orm_workbench/controller"
	"github.com/weikaishio/redis_orm_workbench/common"
)

/*
https://gin-gonic.com/api-example/
*/
func setupRouter() *gin.Engine {

	controller.InitBiz()

	r := gin.Default()
	//r.Delims("{[{", "}]}")
	r.SetFuncMap(template.FuncMap{
		"FormatInterface2Time": common.FormatInterface2Time,
	})

	r.LoadHTMLGlob("views/*")
	r.StaticFS("/static", http.Dir("static/"))
	r.GET("/monitor", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	r.GET("/index", controller.Index)
	r.GET("/schema", controller.Schema)
	r.Any("/data_list", controller.DataList)
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

	svr := http.Server{
		Addr:    ":8881",
		Handler: r,
	}
	go func() {
		if err := svr.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("listen: %s\n", err)
		}
	}()

	log.Info("Start~")
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := svr.Shutdown(ctx); err != nil {
		log.Error("svr.Shutdown err:%v", err)
	}

	log.Info("Stop~")
}
