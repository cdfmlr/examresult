package server

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	ginlogrus "github.com/toorop/gin-logrus" // Logrus logger middleware for Gin
)

// HttpRouter 统一的 Http router。
// 用 HttpRouter.GET 去添加即可
//
// MAIN: 这个东西要在 main 里去 HttpRouter.Run()
var HttpRouter *gin.Engine

func initHttp() {
	// Logrus logger middleware
	log := logrus.New()
	logMiddleware := ginlogrus.Logger(log)

	log.Info("init HttpRouter")

	HttpRouter = gin.New()
	HttpRouter.Use(logMiddleware, gin.Recovery())

	// A simple ping-pong
	HttpRouter.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"ExamResultHttpRouter": "pong",
		})
	})
}
