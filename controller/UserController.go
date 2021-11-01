package controller

import (
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"tricyzhou.com/ginessential/util"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"tricyzhou.com/ginessential/common"
	"tricyzhou.com/ginessential/model"
)

func Register(ctx *gin.Context) {
	DB := common.GetDB()
	name := ctx.PostForm("name")
	tel := ctx.PostForm("telephone")
	passwd := ctx.PostForm("password")

	// 手机号11位
	if len(tel) != 11 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "手机号必须为11位"})
		return
	}

	// 密码不少于6位
	if len(passwd) < 6 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "密码不能少于6位"})
		return
	}

	// 没传name随机生成
	if len(name) == 0 {
		name = util.RandomString(10)
	}

	log.Println(name, passwd, tel)

	// 电话号码是否存在
	if isTelephoneExist(DB, tel) {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "用户已经存在"})
		return
	}
	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "加密错误"})
		return
	}
	// 新建用户
	newUser := model.User{
		Name:      name,
		Telephone: tel,
		Password:  string(hasedPassword),
	}
	DB.Create(&newUser)

	ctx.JSON(200, gin.H{
		"code":   200,
		"mssage": "注册成功",
	})
}

func Login(ctx *gin.Context) {
	DB := common.GetDB()
	// 获取参数
	tel := ctx.PostForm("telephone")
	passwd := ctx.PostForm("password")
	// 数据验证
	// 手机号11位

	// 判断手机号
	if len(tel) != 11 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "手机号必须为11位"})
		return
	}
	var user model.User
	DB.Where("telephone = ?", tel).First(&user)
	if user.ID == 0 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "用户不存在"})
		return
	}
	// 判断密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(passwd)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "密码错误"})
	}
	// 发放token
	token, err := common.ReleaseToken(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "系统异常"})
		log.Printf("token generate error : %v", err)
		return
	}

	// 返回结果
	ctx.JSON(200, gin.H{
		"code":    200,
		"data":    gin.H{"token": token},
		"message": "登陆成功",
	})

}

func Info(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	ctx.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"user": user}})
}

func isTelephoneExist(db *gorm.DB, tel string) bool {
	var user model.User
	db.Where("telephone = ?", tel).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}
