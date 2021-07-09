package service

import (
	"examresult/server"
	"github.com/gin-gonic/gin"
	"net/http"
)

// RegisterFrontEnd 注册的前端网页服务
func RegisterFrontEnd(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", gin.H{})
}

func initRegisterFrontEnd() {
	server.HttpRouter.LoadHTMLGlob("asset/register.html")
	server.HttpRouter.GET("/register", RegisterFrontEnd)
}