package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/weikaishio/redis_orm_workbench/common"
	"github.com/weikaishio/redis_orm_workbench/common/captcha"
	"github.com/weikaishio/redis_orm_workbench/config"
	"net/http"
)

type CaptchaModel struct {
	Id  string
	Src string
}

func Login(c *gin.Context) {
	c.HTML(http.StatusOK, "login.tmpl", gin.H{
		"title":   "Workbench",
		"captcha": getCaptcha(),
	})
}
func LogOut(c *gin.Context) {
	c.SetCookie("user", "", 0, "/", "", false, false)
	//Login(c)
	c.Redirect(http.StatusOK,"/login")
}
func LoginSubmit(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	captchaCode := c.PostForm("captchaCode")
	captchaID := c.PostForm("captchaId")

	if !captcha.VerifyString(captchaID, captchaCode) {
		/*
			data["callbackType"] = call
			data["forwardUrl"] = url
		*/
		c.JSON(http.StatusOK, map[string]string{"statusCode": "300",
			"message":  "验证码错误",
			"navTabId": "login"})
		return
	}
	domain := c.Request.Host
	if val, ok := config.Cfg.UserMap[username]; ok && val == password {
		c.SetCookie("user", username+":"+common.EncryptRC4Base64([]byte(username), password), 0, "/", domain, false, false)
		c.JSON(http.StatusOK, map[string]string{"statusCode": "200",
			"message":    "成功",
			"forwardUrl": "/index",
			"navTabId":   ""})
	} else {
		c.JSON(http.StatusOK, map[string]string{"statusCode": "300",
			"message":  "验证失败",
			"navTabId": "login"})
	}
}

func GetCaptchaImage(c *gin.Context) {
	id := c.Query("id")
	captcha.WriteImage(c.Writer, id, 120, 40)
}

func getCaptcha() *CaptchaModel {
	capResult := new(CaptchaModel)
	capResult.Id = captcha.NewLen(4)

	capResult.Src = fmt.Sprintf(
		"/login/getCaptchaImage?id=%s",
		capResult.Id,
	)
	return capResult
}

func GetCaptcha(c *gin.Context) {
	c.JSON(http.StatusOK, getCaptcha())
}
