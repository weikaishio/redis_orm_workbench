package main

import (
	"context"
	"flag"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"time"

	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mkideal/log"
	"github.com/mkideal/pkg/osutil/pid"
	"github.com/weikaishio/redis_orm_workbench/common"
	"github.com/weikaishio/redis_orm_workbench/config"
	"github.com/weikaishio/redis_orm_workbench/controller"
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
		"IsTime":               common.IsTime,
		"LimitStrLen":          common.LimitStrLen,
	})
	r.Use(controller.UseMiddleware)

	r.LoadHTMLGlob("views/*")
	r.StaticFS("/static", http.Dir("static/"))
	r.GET("/monitor", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})
	r.GET("/login", controller.Login)
	r.POST("/login", controller.LoginSubmit)
	r.GET("/login/logOut", controller.LogOut)
	r.POST("/login/captcha", controller.GetCaptcha)
	r.GET("/login/getCaptchaImage", controller.GetCaptchaImage)

	r.GET("/", controller.Index)
	r.GET("/index", controller.Index)

	r.GET("/schema", controller.Schema)
	r.GET("/schema/create_table", controller.CreateTable)
	r.POST("/schema/drop_table", controller.DropTable)
	r.Any("/schema/table_create",controller.CreateTable)

	r.Any("/data_list", controller.DataList)
	r.POST("/data_list/del", controller.DataDel)
	r.Any("/data_list/edit", controller.DataEdit)
	r.POST("/data_list/truncate_table", controller.TruncateTable)
	r.POST("/data_list/rebuild_index", controller.RebuildIndex)
	return r
}

var (
	pidFile  = flag.String("pid", "redis_orm_workbench.pid", "pid file")
	dir      = flag.String("dir", "./config/", "config path")
	location = flag.String("loc", "local", "server location")
)

func main() {
	flag.Parse()
	log.SetLevelFromString("TRACE")
	var err error
	config.Cfg, err = config.NewConfig(*dir, *location)
	log.If(err != nil).Fatal("NewConfig(%s,%s) err:%v", *dir, *location, err)

	err = config.Cfg.Reload()
	log.If(err != nil).Fatal("config.Cfg.Reload() err:%v", err)

	log.Trace("cfg.UserMap:%v", config.Cfg.UserMap)

	err = pid.New(*pidFile)
	log.If(err != nil).Fatal("pid.New: %v", err)

	defer func() {
		pid.Remove(*pidFile)
		log.Uninit(nil)
	}()

	r := setupRouter()

	svr := http.Server{
		Addr:    fmt.Sprintf(":%d", config.Cfg.HttpPort),
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
