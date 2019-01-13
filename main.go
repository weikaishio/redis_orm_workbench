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

func CookieMiddleware(c *gin.Context) {
	//urlStr := c.Request.URL.String()
	//arr := strings.Split(urlStr, "/")
	//if len(arr) > 1 {
	//	first := arr[1]
	//	if first == "static" || first == "favicon.ico" || first == "login" {
	//		//log.Info("不验证cookie url = %s", urlStr)
	//	} else {
	//		val, _ := c.Cookie("ops")
	//		//log.Info("验证cookie url = %s,cookie=%s", urlStr, val)
	//		for _, cookie := range controllers.StoredMd5CookieStr {
	//			if cookie == val {
	//				c.Next()
	//				return
	//			}
	//		}
	//		controllers.UserVisit(c)
	//		c.Abort()
	//		return
	//	}
	//} else {
	//	controllers.UserVisit(c)
	//	c.Abort()
	//	return
	//}
	c.Next()
}
func setupRouter() *gin.Engine {

	controller.InitBiz()

	r := gin.Default()

	r.Use(CookieMiddleware)
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
