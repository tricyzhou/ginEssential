package main

import (
	"github.com/gin-gonic/gin"
	"tricyzhou.com/ginessential/controller"
	"tricyzhou.com/ginessential/middleware"
)

func CollectRouter(r *gin.Engine) *gin.Engine {
	r.POST("/api/auth/register", controller.Register)
	r.POST("/api/auth/login", controller.Login)
	r.GET("/api/auth/info", middleware.AuthMiddleware(), controller.Info)
	return r
}
