package middleware

import (
	"net/http"
	"strings"

	"tricyzhou.com/ginessential/model"

	"tricyzhou.com/ginessential/common"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 获取authorization header
		tokenString := ctx.GetHeader("Authorization")
		// valiate token formate
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "权限不足1"})
			ctx.Abort() // 抛弃请求
			return
		}

		tokenString = tokenString[7:] // 去掉Bearer

		token, claims, err := common.ParseToken(tokenString)
		if err != nil || !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "权限不足2"})
			ctx.Abort()
			return
		}
		// 验证通过后获取claim的userID
		userId := claims.UserId
		DB := common.GetDB()
		var user model.User
		DB.First(&user, userId)

		// 用户
		if user.ID == 0 {
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "权限不足3"})
			ctx.Abort()
			return
		}

		// 用户存在，将user信息写入上下文!!!
		ctx.Set("user", user)
		ctx.Next()
	}
}
